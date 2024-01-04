package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"np_consumer/internal/db/gen"
	"np_consumer/internal/models"
	"strconv"
)

type Config struct {
	DBHost string
	DBPort string

	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

type Service struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
	config *Config
}

func Init(logger *zap.Logger, cfg *Config) (*Service, error) {
	urlDB := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.PostgresDB)

	poolConfig, err := pgxpool.ParseConfig(urlDB)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = 64

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return &Service{
		pool:   dbPool,
		logger: logger,
	}, nil
}

func BeginTransaction(pool *pgxpool.Pool) (*gen.Queries, *sql.Tx, error) {
	singleConn := stdlib.OpenDBFromPool(pool)
	defer singleConn.Close()

	tx, err := singleConn.Begin()
	if err != nil {
		return nil, nil, err
	}
	qtx := gen.New(singleConn).WithTx(tx)

	return qtx, tx, nil
}

func (d *Service) SaveOrder(ms models.Shipment) error {
	qtx, tx, err := BeginTransaction(d.pool)
	if err != nil {
		d.logger.Error("failed to create query transaction:", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	ca := ms.Customer.Address
	err = qtx.CreateAddress(context.Background(), gen.CreateAddressParams{
		AddressID: ca.AddressID,
		Country:   ca.Country,
		City:      ca.City,
		Street:    ca.Street,
		ZipCode:   ca.Zip,
	})

	c := ms.Customer
	err = qtx.CreateCustomer(context.Background(), gen.CreateCustomerParams{
		CustomerID:        c.CustomerID,
		CustomerAddressID: c.Address.AddressID,
		Name:              c.Name,
		LastName:          c.LastName,
		Email:             c.Email,
		PhoneNumber:       c.PhoneNumber,
	})

	sa := ms.Sender.Address
	err = qtx.CreateAddress(context.Background(), gen.CreateAddressParams{
		AddressID: sa.AddressID,
		Country:   sa.Country,
		Street:    sa.Street,
		City:      sa.City,
		ZipCode:   sa.Zip,
	})

	s := ms.Sender
	sp, _ := strconv.Atoi(s.PhoneNumber)
	err = qtx.CreateSender(context.Background(), gen.CreateSenderParams{
		SenderID:        s.SenderID,
		SenderAddressID: s.Address.AddressID,
		Name:            s.Name,
		Email:           s.Email,
		PhoneNumber:     int32(sp),
	})

	err = qtx.CreateShipment(context.Background(), gen.CreateShipmentParams{
		ShipmentID: ms.ShipmentID,
		SenderID:   ms.Sender.SenderID,
		CustomerID: ms.Customer.CustomerID,
		Size:       ms.Size,
		Weight:     ms.Weight,
		Count:      int32(ms.Count),
	})

	err = qtx.CreateEvent(context.Background(), gen.CreateEventParams{
		EventID:          ms.Event.EventID,
		ShipmentID:       ms.ShipmentID,
		EventTimestamp:   ms.Event.EventTime,
		EventDescription: ms.Event.EventDescription,
	})

	return tx.Commit()
}
