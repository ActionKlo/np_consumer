package internal

import (
	"go.uber.org/zap"
	"net"
	"np_consumer/internal/api"
	"np_consumer/internal/db"
	"np_consumer/internal/kafka"
)

type MasterService struct {
	DB     *db.PostgresService
	Kafka  *kafka.Service
	Logger *zap.Logger
}

func (m MasterService) ListenIncomingMessages() {
	if err := m.Kafka.Reader(m.DB); err != nil {
		m.Logger.Fatal("kafka reader fall down")
	}
}

func (m MasterService) StartGRPCServer() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		m.Logger.Fatal("failed to listen address", zap.Error(err))
	}
	srv := api.NewGRPCServer(m.DB, m.Logger)

	m.Logger.Info("gRPC server should been started")

	if err = srv.Serve(listen); err != nil {
		m.Logger.Fatal("failed to start server", zap.Error(err))
	}
}
