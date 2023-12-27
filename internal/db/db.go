package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"np_consumer/config"
	"np_consumer/internal/db/gen"
	"np_consumer/internal/models"
)

type ServiceDB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func Init(logger *zap.Logger, cfg *config.Config) (*ServiceDB, error) {
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

	return &ServiceDB{
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

func (d *ServiceDB) CreateCustomer(customer *models.CustomerInfo) error {
	qtx, tx, err := BeginTransaction(d.pool)
	if err != nil {
		d.logger.Error("failed to create query transaction:", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	addressID, err := qtx.CreateCustomerAddress(context.Background(), gen.CreateCustomerAddressParams{
		Country: customer.Address.Country,
		Street:  customer.Address.Street,
		City:    customer.Address.City,
		ZipCode: customer.Address.ZipCode,
	})

	if err != nil {
		d.logger.Error("failed to create customer address:", zap.Error(err))
		return err
	}

	customerID, err := qtx.CreateCustomer(context.Background(), gen.CreateCustomerParams{
		CustomerAddressID: addressID,
		Name:              customer.Customer.Name,
		Email:             customer.Customer.Email,
		PhoneNumber:       customer.Customer.PhoneNumber,
	})
	if err != nil {
		d.logger.Error("failed to create customer:", zap.Error(err))
		return err
	}
	d.logger.Info("customer created:", zap.Any("customer id:", customerID))

	return tx.Commit()
}

func (d *ServiceDB) GetAllCustomers() ([]gen.GetAllCustomersRow, error) {
	conn := stdlib.OpenDBFromPool(d.pool)
	defer conn.Close()

	q := gen.New(conn)

	res, err := q.GetAllCustomers(context.Background())
	if err != nil {
		d.logger.Error("failed to get customers:", zap.Error(err))
		return nil, err
	}

	fmt.Println(res[0])

	return res, nil
}
