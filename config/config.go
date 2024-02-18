package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"np_consumer/internal/db"
	"np_consumer/internal/kafka"
)

type Config struct {
	DB struct {
		Host string `mapstructure:"DB_HOST"`
		Port string `mapstructure:"DB_PORT"`
	}
	Postgres struct {
		User     string `mapstructure:"POSTGRES_USER"`
		Password string `mapstructure:"POSTGRES_PASSWORD"`
		DB       string `mapstructure:"POSTGRES_DB"`
	}

	Kafka struct {
		ExternalHost string `mapstructure:"KAFKA_EXTERNAL_HOST"`
		GroupID      string `mapstructure:"KAFKA_GROUPID"`
		Topic        string `mapstructure:"KAFKA_TOPIC"`
		Partition    int    `mapstructure:"KAFKA_PARTITION"`
	}
}

type Services struct {
	Kafka *kafka.Service
	DB    *db.PostgresService
}

func New() *Config {
	var appConfig Config
	v := viper.New()
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// TODO why doesn't work: v.Unmarshal(&appConfig)????
	if err := v.Unmarshal(&appConfig.DB); err != nil {
		log.Fatal(err)
	}
	if err := v.Unmarshal(&appConfig.Postgres); err != nil {
		log.Fatal(err)
	}
	if err := v.Unmarshal(&appConfig.Kafka); err != nil {
		log.Fatal(err)
	}

	return &appConfig
}

func (c *Config) NewServices(logger *zap.Logger) *Services {
	kafkaService := kafka.New(logger, &kafka.Config{
		Kafka: kafka.Kafka{
			ExternalHost: c.Kafka.ExternalHost,
			GroupID:      c.Kafka.GroupID,
			Topic:        c.Kafka.Topic,
			Partition:    c.Kafka.Partition,
		},
	})

	DBService := db.Init(logger, &db.Config{
		DB: db.DB{
			Host: c.DB.Host,
			Port: c.DB.Port,
		},
		Postgres: db.Postgres{
			User:     c.Postgres.User,
			Password: c.Postgres.Password,
			DB:       c.Postgres.DB,
		},
	})

	return &Services{
		Kafka: kafkaService,
		DB:    DBService,
	}
}
