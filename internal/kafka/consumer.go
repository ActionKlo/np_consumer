package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"np_consumer/config"
	"np_consumer/internal/db"
	"np_consumer/internal/types"
	"sync"
	"time"
)

type ServiceKafka struct {
	logger *zap.Logger
	config config.Config
	dbStr  *db.DB
}

func NewKafka(logger *zap.Logger, config config.Config, dbStr *db.DB) *ServiceKafka {
	return &ServiceKafka{
		logger: logger,
		config: config,
		dbStr:  dbStr,
	}
}

type KafkaMessage struct {
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

//func CreatePool() (*pgxpool.Pool, error) {
//	url := "postgresql://consumerAdmin:supersecret@100.66.158.79:5430/consumerdb"
//	dbPool, err := pgxpool.New(context.Background(), url)
//	if err != nil {
//		return nil, err
//	}
//
//	return dbPool, nil
//}

func (k *ServiceKafka) Reader() error {
	k.logger.Info("kafka consumer started")

	//pool, err := CreatePool()
	//if err != nil {
	//	k.logger.Error("failed create pgxpool:", zap.Error(err))
	//	return err
	//}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.config.KafkaExternalHost},
		//GroupID:   k.config.KafkaGroupID, // uncomment after start work with db
		Topic:     k.config.KafkaTopic,
		Partition: k.config.KafkaPartition,
		MaxBytes:  10e6, // 10MB
	})

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 500; i++ {
		wg.Add(1)
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			k.logger.Error("error read message:", zap.Error(err))
			break
		}

		var ms types.KafkaMessage
		if err = json.Unmarshal(m.Value, &ms); err != nil {
			k.logger.Error("failed to unmarshal message:", zap.Error(err))
			return err
		}

		go func(ms types.KafkaMessage) {
			defer wg.Done()
			err := k.dbStr.InsertMessage(ms)
			if err != nil {
				k.logger.Error("failed to insert message", zap.Error(err))
			}
		}(ms)

		/*
			go func(ms KafkaMessage) {
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
		*/

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	wg.Wait()
	fmt.Println(time.Since(start))

	if err := r.Close(); err != nil {
		k.logger.Error("failed to close reader:", zap.Error(err))
		return err
	}

	k.logger.Info("kafka consumer ended")
	return nil
}
