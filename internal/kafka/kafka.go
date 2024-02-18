package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"np_consumer/internal/db"
	"np_consumer/internal/models"
	"np_consumer/internal/webhooks"
)

type (
	Config struct {
		Kafka Kafka
	}

	Kafka struct {
		ExternalHost string
		GroupID      string
		Topic        string
		Partition    int
	}
)

type Service struct {
	logger *zap.Logger
	config *Config
	reader *kafka.Reader
}

func New(logger *zap.Logger, config *Config) *Service {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Kafka.ExternalHost},
		//GroupID:   config.Kafka.GroupID, // uncomment after start work with db
		Topic:     config.Kafka.Topic,
		Partition: config.Kafka.Partition,
		MaxBytes:  10e6,
	})

	return &Service{
		logger: logger,
		config: config,
		reader: reader,
	}
}

func (k *Service) Reader(dbService *db.PostgresService) error {
	k.logger.Info("kafka consumer started")

	for {
		m, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			k.logger.Error("error read message", zap.Error(err))
			break
		}

		var payload models.Payload
		if err = json.Unmarshal(m.Value, &payload); err != nil {
			k.logger.Error("failed to unmarshal message: "+string(m.Value), zap.Error(err))
			return err
		}

		go func(ms models.Payload) {
			err = dbService.SavePayload(ms)
			if err != nil {
				k.logger.Error("failed to insert message: "+ms.MessageID.String(), zap.Error(err))
			} else {
				if webhookUrl := dbService.GetSettingsByReceiverID(ms.ReceiverID); webhookUrl != "" {
					err = wh.SendNotification(webhookUrl, ms, k.logger)
					if err != nil {
						k.logger.Error("failed to send webhooks: "+ms.ReceiverID.String(), zap.Error(err))
					}
				}
			}
		}(payload)
	}

	if err := k.reader.Close(); err != nil {
		k.logger.Error("failed to close reader:", zap.Error(err))
		return err
	}

	k.logger.Info("kafka consumer ended")

	return nil
}
