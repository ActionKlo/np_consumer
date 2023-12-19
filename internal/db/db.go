package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"np_consumer/config"
	"np_consumer/internal/db/gen"
	"np_consumer/internal/types"
)

type DB struct {
	Pool *pgxpool.Pool
}

func CreateDB(cfg *config.Config) (*DB, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.PostgresDB,
	)

	dbCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	dbCfg.MaxConns = 64

	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbCfg)
	if err != nil {
		return nil, err
	}

	return &DB{Pool: dbPool}, nil
}

func (d *DB) InsertMessage(ms *types.KafkaMessage) error { // TODO think about 2 func: InsMes and InsStatus
	// TODO Add logger? For what?
	// but if add from kafka it will be circle import error
	conn := stdlib.OpenDBFromPool(d.Pool) // why is it here? But why not?
	defer conn.Close()

	q := gen.New(conn) // TODO How to close ConnectionFromPool(d.Pool) connection inside New() ????

	_, err := q.InsertMessage(context.Background(), gen.InsertMessageParams{
		ID:          ms.ID,
		Time:        ms.Time,
		Sender:      ms.Sender,
		Tracknumber: ms.TrackNumber,
		Country:     ms.Country,
		City:        ms.City,
		Street:      ms.Street,
		Postcode:    ms.PostCode,
	})
	if err != nil {
		return err
	}

	_, err = q.InsertStatus(context.Background(), gen.InsertStatusParams{
		ID:        uuid.NewString(),
		Messageid: ms.ID,
		Status:    ms.Status,
		Time:      ms.Time,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetMessageByID(id string) (*gen.Message, error) {
	conn := stdlib.OpenDBFromPool(d.Pool)
	defer conn.Close()

	q := gen.New(conn)

	resMes, err := q.GetMessageByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &resMes, nil // why need pointer?
}
