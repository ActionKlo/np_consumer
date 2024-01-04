package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"np_consumer/internal/db"
	"np_consumer/internal/models"
)

type Config struct {
	KafkaExternalHost string
	KafkaGroupID      string
	KafkaTopic        string
	KafkaPartition    int
}

type Service struct {
	logger *zap.Logger
	config *Config
}

func New(logger *zap.Logger, config *Config) *Service {
	return &Service{
		logger: logger,
		config: config,
	}
}

func (k *Service) Reader(dbService *db.Service) error {
	k.logger.Info("kafka consumer started")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.config.KafkaExternalHost},
		//GroupID:   k.config.KafkaGroupID, // uncomment after start work with db
		Topic:     k.config.KafkaTopic,
		Partition: k.config.KafkaPartition,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil { // TODO add retry if kafka_producer fall down
			k.logger.Error("error read message:", zap.Error(err))
			break
		}

		var ms models.Shipment
		if err = json.Unmarshal(m.Value, &ms); err != nil {
			k.logger.Error("failed to unmarshal message:", zap.Error(err))
			return err
		}

		go func(ms models.Shipment) {
			err = dbService.SaveOrder(ms)
			if err != nil {
				k.logger.Error("failed to insert message", zap.Error(err))
			}
			fmt.Println(ms)
		}(ms)
	}

	if err := r.Close(); err != nil {
		k.logger.Error("failed to close reader:", zap.Error(err))
		return err
	}

	k.logger.Info("kafka consumer ended")
	return nil
}
