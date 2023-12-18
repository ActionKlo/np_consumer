package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	//KafkaHost         string `mapstructure:"KAFKA_HOST"`
	KafkaExternalHost string `mapstructure:"KAFKA_EXTERNAL_HOST"`
	KafkaGroupID      string `mapstructure:"KAFKA_GROUPID"`
	KafkaTopic        string `mapstructure:"KAFKA_TOPIC"`
	KafkaPartition    int    `mapstructure:"KAFKA_PARTITION"`
}

//type Services struct {
//	Db    db.DB
//	Kafka kafka.ServiceKafka
//}

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