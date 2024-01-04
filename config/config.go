package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"np_consumer/internal/db"
	"np_consumer/internal/kafka"
)

type Config struct {
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`

	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`

	KafkaExternalHost string `mapstructure:"KAFKA_EXTERNAL_HOST"`
	KafkaGroupID      string `mapstructure:"KAFKA_GROUPID"`
	KafkaTopic        string `mapstructure:"KAFKA_TOPIC"`
	KafkaPartition    int    `mapstructure:"KAFKA_PARTITION"`
}

type Services struct {
	KafkaService *kafka.Service
	DBService    *db.Service
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
	if err := v.Unmarshal(&appConfig); err != nil {
		log.Fatal(err)
	}
	return &appConfig
}

func (c *Config) NewKafkaConfig(logger *zap.Logger) *Services {
	kafkaConfig := kafka.New(logger, &kafka.Config{
		KafkaExternalHost: c.KafkaExternalHost,
		KafkaGroupID:      c.KafkaGroupID,
		KafkaTopic:        c.KafkaTopic,
		KafkaPartition:    0,
	})

	return &Services{
		KafkaService: kafkaConfig,
	}
}

func (c *Config) NewDBConfig(logger *zap.Logger) *Services {
	dbConfig, err := db.Init(logger, &db.Config{
		DBHost:           c.DBHost,
		DBPort:           c.DBPort,
		PostgresUser:     c.PostgresUser,
		PostgresPassword: c.PostgresPassword,
		PostgresDB:       c.PostgresDB,
	})
	if err != nil {
		logger.Error("for what?", zap.Error(err))
	}

	return &Services{
		DBService: dbConfig,
	}
}
