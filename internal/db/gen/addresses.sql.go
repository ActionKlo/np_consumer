// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: addresses.sql

package gen

import (
	"context"

	"github.com/google/uuid"
)

const createAddress = `-- name: CreateAddress :exec
INSERT INTO addresses (
    address_id, country, street, city, zip_code
) VALUES (
    $1, $2, $3, $4, $5
)ON CONFLICT DO NOTHING
`

type CreateAddressParams struct {
	AddressID uuid.UUID
	Country   string
	Street    string
	City      string
	ZipCode   string
}

func (q *Queries) CreateAddress(ctx context.Context, arg CreateAddressParams) error {
	_, err := q.db.ExecContext(ctx, createAddress,
		arg.AddressID,
		arg.Country,
		arg.Street,
		arg.City,
		arg.ZipCode,
	)
	return err
}
