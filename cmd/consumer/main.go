package main

import (
	"np_consumer/config"
	"np_consumer/logger"
)

func main() {
	log := logger.Init()
	cfg := config.New()

	kafkaCfg := cfg.NewKafkaConfig(log)
	k := kafkaCfg.Kafka

	dbCfg := cfg.NewDBConfig(log)
	dbService := dbCfg.DB

	//k := cfg.NewKafkaConfig().Kafka // just reminder

	if err := k.Reader(dbService); err != nil {
		log.Fatal("kafka reader fall down")
	}
}
