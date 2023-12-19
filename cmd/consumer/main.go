package main

import (
	"go.uber.org/zap"
	"np_consumer/config"
	"np_consumer/internal/db"
	"np_consumer/internal/kafka"
	"np_consumer/logger"
)

func main() {
	cfg := config.New()
	log := logger.Init()

	dbCfg, err := db.CreateDB(cfg)
	if err != nil {
		log.Fatal("failed create pgxpool:", zap.Error(err))
	}

	k := kafka.NewKafka(log, *cfg, dbCfg)
	if err := k.Reader(); err != nil {
		log.Fatal("kafka reader fall down")
	}
}
