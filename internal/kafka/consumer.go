package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
	"np_consumer/config"
	"np_consumer/internal/db"
	"np_consumer/internal/models"
)

type ServiceKafka struct {
	logger    *zap.Logger
	config    *config.Config
	serviceDB *db.ServiceDB
}

func New(logger *zap.Logger, config *config.Config, serviceDB *db.ServiceDB) *ServiceKafka {
	serviceDB, err := db.Init(logger, config)
	if err != nil {
		log.Fatal("failed to create new db service", zap.Error(err))
	}

	return &ServiceKafka{
		logger:    logger,
		config:    config,
		serviceDB: serviceDB,
	}
}

func (k *ServiceKafka) Reader() error {
	k.logger.Info("kafka consumer started")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.config.KafkaExternalHost},
		//GroupID:   k.config.KafkaGroupID, // uncomment after start work with db
		Topic:     k.config.KafkaTopic,
		Partition: k.config.KafkaPartition,
		MaxBytes:  10e6,
	})

	r.SetOffset(19818)

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
			err = k.serviceDB.SaveMeessage(ms)
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
