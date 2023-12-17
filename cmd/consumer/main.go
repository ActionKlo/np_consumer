package main

import (
	"np_consumer/config"
	"np_consumer/internal/kafka"
	"np_consumer/logger"
)

func main() {
	cfg := config.New()
	log := logger.Init()

	k := kafka.NewKafka(log, *cfg)
	if err := k.Reader(); err != nil {
		log.Error("kafka reader fall down")
	}
}
