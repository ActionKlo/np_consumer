package main

import (
	"fmt"
	"go.uber.org/zap"
	"np_consumer/api/handlers"
	"np_consumer/config"
	"np_consumer/internal/db"
	"np_consumer/logger"
	"time"
)

func main() {
	log := logger.Init()
	cfg := config.New()

	d, err := db.Init(log, cfg)
	if err != nil {
		log.Fatal("failed to init database:", zap.Error(err))
	}
	t := time.Now()

	handlers.CreateSomeConsumers(d)

	fmt.Println(time.Since(t))

	customers, err := handlers.GetALlCustomersWithAddress(d)
	if err != nil {
		log.Fatal("failed to get customers", zap.Error(err))
	}
	_ = customers

	//k := kafka.New(log, cfg)
	//
	//if err := k.Reader(); err != nil {
	//	log.Fatal("kafka reader fall down")
	//}
}
