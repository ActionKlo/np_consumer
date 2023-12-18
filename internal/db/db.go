package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"np_consumer/internal/db/gen"
	"np_consumer/internal/types"
)

type DB struct {
	Pool *pgxpool.Pool
}

func CreatePool() (*pgxpool.Pool, error) {
	url := "postgresql://consumerAdmin:supersecret@100.66.158.79:5430/consumerdb"
	c, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	c.MaxConns = 64

	dbPool, err := pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func NewDB() (*DB, error) {
	conn, err := CreatePool()
	if err != nil {
		return nil, err
	}

	return &DB{
		Pool: conn,
	}, nil
}

func (d *DB) InsertMessage(ms types.KafkaMessage) error { // TODO thing about 2 func: InsMes and InsStatus
	conn := stdlib.OpenDBFromPool(d.Pool)

	q := gen.New(conn)

	resMs, err := q.InsertMessage(context.Background(), gen.InsertMessageParams{
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
		fmt.Println("failed to insert message:", zap.Error(err))
		return err
	}
	fmt.Println(resMs.ID)

	resSt, err := q.InsertStatus(context.Background(), gen.InsertStatusParams{
		ID:        uuid.NewString(),
		Messageid: ms.ID,
		Status:    ms.Status,
		Time:      ms.Time,
	})
	if err != nil {
		fmt.Println("failed to insert status", zap.Error(err))
		return err
	}
	fmt.Println(resSt.ID)

	return nil
}
