package kafka

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

func Reader() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	brokerUrl := os.Getenv("HOST") + ":9094"
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerUrl},
		Topic:     "test-topic",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	fmt.Println("Consumer started")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
