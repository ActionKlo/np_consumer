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
	"strconv"
)

type ServiceDB struct {
	Pool   *pgxpool.Pool
	Logger *zap.Logger
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
		Pool:   dbPool,
		Logger: logger,
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
	//qtx, tx, err := BeginTransaction(d.Pool)
	//if err != nil {
	//	d.Logger.Error("failed to create query transaction:", zap.Error(err))
	//	return err
	//}
	//defer tx.Rollback()
	//
	//addressID, err := qtx.CreateCustomerAddress(context.Background(), gen.CreateCustomerAddressParams{
	//	Country: customer.Address.Country,
	//	Street:  customer.Address.Street,
	//	City:    customer.Address.City,
	//	ZipCode: customer.Address.ZipCode,
	//})
	//
	//if err != nil {
	//	d.Logger.Error("failed to create customer address:", zap.Error(err))
	//	return err
	//}
	//
	//customerID, err := qtx.CreateCustomer(context.Background(), gen.CreateCustomerParams{
	//	CustomerAddressID: addressID,
	//	Name:              customer.Customer.Name,
	//	Email:             customer.Customer.Email,
	//	PhoneNumber:       customer.Customer.PhoneNumber,
	//})
	//if err != nil {
	//	d.Logger.Error("failed to create customer:", zap.Error(err))
	//	return err
	//}
	//d.Logger.Info("customer created:", zap.Any("customer id:", customerID))
	//
	//return tx.Commit()
	return nil
}

func (d *ServiceDB) GetAllCustomers() ([]gen.GetAllCustomersRow, error) {
	conn := stdlib.OpenDBFromPool(d.Pool)
	defer conn.Close()

	q := gen.New(conn)

	res, err := q.GetAllCustomers(context.Background())
	if err != nil {
		d.Logger.Error("failed to get customers:", zap.Error(err))
		return nil, err
	}

	fmt.Println(res[0])

	return res, nil
}

func (d *ServiceDB) SaveMeessage(ms models.Shipment) error {
	qtx, tx, err := BeginTransaction(d.Pool)
	if err != nil {
		d.Logger.Error("failed to create query transaction:", zap.Error(err))
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
		Weight:     float32(ms.Weight),
		Count:      int32(ms.Count),
	})

	err = qtx.CreateEvent(context.Background(), gen.CreateEventParams{
		StatusID:         ms.Event.EventID,
		ShipmentID:       ms.ShipmentID,
		EventTimestamp:   ms.Event.EventTime,
		EventDescription: ms.Event.EventDescription,
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}
