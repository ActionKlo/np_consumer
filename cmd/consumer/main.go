package main

import (
	"go.uber.org/zap"
	"np_consumer/config"
	"np_consumer/internal/db"
	"np_consumer/internal/kafka"
	"np_consumer/logger"
)

func main() {
	log := logger.Init()
	cfg := config.New()

	serviceDB, err := db.Init(log, cfg)
	if err != nil {
		log.Fatal("failed to init database:", zap.Error(err))
	}

	k := kafka.New(log, cfg, serviceDB)

	if err := k.Reader(); err != nil {
		log.Fatal("kafka reader fall down")
	}
}
