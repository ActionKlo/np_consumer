package main

import (
	"go.uber.org/zap"
	"net"
	"np_consumer/config"
	"np_consumer/internal"
	"np_consumer/internal/api"
	"np_consumer/logger"
)

func main() {
	log := logger.Init()
	cfg := config.New()

	services := cfg.NewServices(log)

	masterService := internal.MasterService{ // is it equal kafkaService?
		DB:     services.DB,
		Kafka:  services.Kafka,
		Logger: log,
	}

	go masterService.ListenIncomingMessages()

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen address", zap.Error(err))
	}
	srv := api.NewGRPCServer(*services.DB)

	if err = srv.Serve(listen); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}
}
