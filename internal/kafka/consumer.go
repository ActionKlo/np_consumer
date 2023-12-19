package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"np_consumer/config"
	"np_consumer/internal/db"
	"np_consumer/internal/types"
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
			k.logger.Error("error read message:", zap.Error(err))
			break
		}

		var ms types.KafkaMessage
		if err = json.Unmarshal(m.Value, &ms); err != nil {
			k.logger.Error("failed to unmarshal message:", zap.Error(err))
			return err
		}

		go func(ms types.KafkaMessage) {
			err := k.dbStr.InsertMessage(&ms)
			if err != nil {
				k.logger.Error("failed to insert message", zap.Error(err))
			}
		}(ms)
		//fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		k.logger.Error("failed to close reader:", zap.Error(err))
		return err
	}

	k.logger.Info("kafka consumer ended")
	return nil
}
