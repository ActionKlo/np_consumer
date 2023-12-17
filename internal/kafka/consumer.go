package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"np_consumer/config"
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

func (k *ServiceKafka) Reader() error {
	k.logger.Info("kafka consumer started")

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
			k.logger.Error("error read message: ", zap.Error(err))
			break
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		k.logger.Error("failed to close reader:", zap.Error(err))
		return err
	}

	k.logger.Info("kafka consumer ended")
	return nil
}
