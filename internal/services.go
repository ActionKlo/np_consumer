package internal

import (
	"go.uber.org/zap"
	"np_consumer/internal/db"
	"np_consumer/internal/kafka"
)

type MasterService struct {
	DB     *db.Service
	Kafka  *kafka.Service
	Logger *zap.Logger
}

func (m MasterService) ListenIncomingMessages() {
	if err := m.Kafka.Reader(m.DB); err != nil {
		m.Logger.Fatal("kafka reader fall down")
	}
}

func (m MasterService) StartGRPCServer() {

}
