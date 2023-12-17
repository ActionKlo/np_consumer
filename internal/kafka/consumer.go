package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"np_consumer/config"
	"np_consumer/internal/db/gen"
	"time"
)

type ServiceKafka struct {
	logger *zap.Logger
	config config.Config
}

func NewKafka(logger *zap.Logger, config config.Config) *ServiceKafka {
	return &ServiceKafka{
		logger: logger,
		config: config,
	}
}

type Message struct {
	ID          string    `json:"ID"`
	Time        time.Time `json:"Time"`
	Sender      string    `json:"Sender"`
	Status      string    `json:"Status"`
	TrackNumber string    `json:"TrackNumber"`
	Country     string    `json:"Country"`
	City        string    `json:"City"`
	Street      string    `json:"Street"`
	PostCode    string    `json:"PostCode"`
}

func CreatePool() (*pgxpool.Pool, error) {
	url := "postgresql://consumerAdmin:supersecret@100.66.158.79:5430/consumerdb"
	dbPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func (k *ServiceKafka) Reader() error {
	k.logger.Info("kafka consumer started")

	pool, err := CreatePool()
	if err != nil {
		k.logger.Error("failed create pgxpool:", zap.Error(err))
		return err
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.config.KafkaExternalHost},
		//GroupID:   k.config.KafkaGroupID, // uncomment after start work with db
		Topic:     k.config.KafkaTopic,
		Partition: k.config.KafkaPartition,
		MaxBytes:  10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			k.logger.Error("error read message:", zap.Error(err))
			break
		}

		var ms Message
		if err = json.Unmarshal(m.Value, &ms); err != nil {
			k.logger.Error("failed to unmarshal message:", zap.Error(err))
			return err
		}

		go func(ms Message) {
			conn := stdlib.OpenDBFromPool(pool)

			q := gen.New(conn)

			resMs, err := q.InsertMessage(context.Background(), gen.InsertMessageParams{
				ID:          ms.ID,
				Time:        ms.Time,
				Sender:      ms.Sender,
				Tracknumber: ms.TrackNumber,
				Country:     ms.Country,
				City:        ms.City,
				Street:      ms.Street,
				Postcode:    ms.PostCode,
			})
			if err != nil {
				k.logger.Error("failed to insert message:", zap.Error(err))
			}
			k.logger.Info(resMs.ID)

			resSt, err := q.InsertStatus(context.Background(), gen.InsertStatusParams{
				ID:        uuid.NewString(),
				Messageid: ms.ID,
				Status:    ms.Status,
				Time:      ms.Time,
			})
			if err != nil {
				k.logger.Error("failed to insert status", zap.Error(err))
			}
			k.logger.Info(resSt.ID)
		}(ms)

		k.logger.Info(ms.ID)

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		k.logger.Error("failed to close reader:", zap.Error(err))
		return err
	}

	k.logger.Info("kafka consumer ended")
	return nil
}
