package main

import (
	"np_consumer/config"
	"np_consumer/internal"
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

	masterService.ListenIncomingMessages()
}
