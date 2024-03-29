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
	gapi "np_consumer/proto"
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

type PostgresService struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
	config *Config
}

func Init(logger *zap.Logger, cfg *Config) *PostgresService {
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

	return &PostgresService{
		pool:   dbPool,
		logger: logger,
	}
}

func (d *PostgresService) SavePayload(pl models.Payload) error {
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

func (d *PostgresService) GetSettingsByReceiverID(receiverID uuid.UUID) string {
	q := gen.New(stdlib.OpenDBFromPool(d.pool))
	url, err := q.GetSettingsByReceiverID(context.TODO(), receiverID)
	if err != nil {
		d.logger.Error("failed to get setting by receiverID", zap.Error(err))
		return ""
	}

	return url
}

type ReceiverRepository interface {
	CreateReceiver(ctx context.Context, receiver *gapi.Receiver) (uuid.UUID, error)
	RetrieveReceiver(ctx context.Context, id uuid.UUID) (*gapi.Receiver, error)
	UpdateReceiver(ctx context.Context, receiver *gapi.Receiver) error
	DeleteReceiver(ctx context.Context, id uuid.UUID) error
}

func (d *PostgresService) CreateReceiver(ctx context.Context, receiver *gapi.Receiver) (uuid.UUID, error) {
	q := gen.New(stdlib.OpenDBFromPool(d.pool))

	receiverID, err := q.CreateReceiver(ctx, gen.CreateReceiverParams{
		ReceiverID: uuid.New(),
		Url:        receiver.Url,
	})
	if err != nil {
		d.logger.Error("failed to create receiver", zap.Error(err))
		return uuid.Nil, err // TODO is it correct?
	}

	return receiverID, nil
}

func (d *PostgresService) RetrieveReceiver(ctx context.Context, id uuid.UUID) (*gapi.Receiver, error) {
	q := gen.New(stdlib.OpenDBFromPool(d.pool))

	//var receiver *gapi.Receiver
	receiver, err := q.RetrieveReceiver(ctx, id)
	if err != nil {
		d.logger.Error("failed to retrieve receiver", zap.Error(err))
		return nil, err
	}
	return &gapi.Receiver{
		Id:  receiver.ReceiverID.String(),
		Url: receiver.Url,
	}, nil
}

func (d *PostgresService) UpdateReceiver(ctx context.Context, receiver *gapi.Receiver) error {
	rid, err := uuid.Parse(receiver.Id)
	if err != nil {
		d.logger.Error("failed to parse uuid", zap.Error(err))
		return err
	}

	q := gen.New(stdlib.OpenDBFromPool(d.pool))

	rows, err := q.UpdateReceiver(ctx, gen.UpdateReceiverParams{
		ReceiverID: rid,
		Url:        receiver.Url,
	})
	if err != nil {
		d.logger.Error("filed to update receiver", zap.Error(err))
		return err
	}

	if rows == 0 {
		// TODO should I return error?
		d.logger.Debug("receiver not found")
	}

	return nil
}

func (d *PostgresService) DeleteReceiver(ctx context.Context, id uuid.UUID) error {
	q := gen.New(stdlib.OpenDBFromPool(d.pool))

	rows, err := q.DeleteReceiver(ctx, id)
	if err != nil {
		d.logger.Error("filed to delete receiver", zap.Error(err))
		return err
	}

	if rows == 0 {
		// TODO should I return error?
		d.logger.Debug("receiver not found")
	}

	return nil
}
