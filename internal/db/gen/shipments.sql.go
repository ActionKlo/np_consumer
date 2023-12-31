// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: shipments.sql

package gen

import (
	"context"

	"github.com/google/uuid"
)

const createShipment = `-- name: CreateShipment :exec
INSERT INTO shipments (
    shipment_id, sender_id, customer_id, size, weight, count
) VALUES (
    $1, $2, $3, $4, $5, $6
)ON CONFLICT DO NOTHING
`

type CreateShipmentParams struct {
	ShipmentID uuid.UUID
	SenderID   uuid.UUID
	CustomerID uuid.UUID
	Size       string
	Weight     float32
	Count      int32
}

func (q *Queries) CreateShipment(ctx context.Context, arg CreateShipmentParams) error {
	_, err := q.db.ExecContext(ctx, createShipment,
		arg.ShipmentID,
		arg.SenderID,
		arg.CustomerID,
		arg.Size,
		arg.Weight,
		arg.Count,
	)
	return err
}
