package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"np_consumer/internal/db/gen"
	"np_consumer/internal/models"
)

const (
	MaxPoolConnections = 64
)

type (
	Config struct {
		DB       DB
		Postgres Postgres
	}

	DB struct {
		Host string
		Port string
	}

	Postgres struct {
		User     string
		Password string
		DB       string
	}
)

type Service struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
	config *Config
}

func Init(logger *zap.Logger, cfg *Config) *Service {
	urlDB := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.Postgres.DB)

	poolConfig, err := pgxpool.ParseConfig(urlDB)
	if err != nil {
		logger.Fatal("failed to parse pool config", zap.Error(err))
	}
	poolConfig.MaxConns = MaxPoolConnections

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatal("failed to create pool", zap.Error(err))
	}

	return &Service{
		pool:   dbPool,
		logger: logger,
	}
}

func (d *Service) SavePayload(pl models.Payload) error {
	data, err := json.Marshal(pl.Order)
	if err != nil {
		d.logger.Error("failed to marshal order data", zap.Error(err))
		return err
	}

	q := gen.New(stdlib.OpenDBFromPool(d.pool))
	messageID, err := q.SavePayload(context.Background(), gen.SavePayloadParams{
		MessageID:      pl.MessageID,
		TrackingNumber: pl.TrackingNumber,
		EventID:        pl.EventID,
		EventType:      pl.EventType,
		EventTime:      pl.EventTime,
		Data:           data,
		ReceiverID:     pl.ReceiverID,
	})
	if err != nil {
		d.logger.Error("failed to save payload", zap.Error(err))
		return err
	}
	d.logger.Info("payload saved", zap.String("messageID", messageID.String()))

	return nil
}

func (d *Service) GetSettingsByReceiverID(receiverID uuid.UUID) string {
	q := gen.New(stdlib.OpenDBFromPool(d.pool))
	url, err := q.GetSettingsByReceiverID(context.TODO(), receiverID)
	if err != nil {
		d.logger.Error("failed to get setting by receiverID", zap.Error(err))
		return ""
	}

	return url
}
