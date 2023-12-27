// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: cutomers.sql

package gen

import (
	"context"

	"github.com/google/uuid"
)

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customers (
    customer_address_id, name, last_name, email, phone_number
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING customer_id
`

type CreateCustomerParams struct {
	CustomerAddressID uuid.UUID
	Name              string
	LastName          string
	Email             string
	PhoneNumber       string
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createCustomer,
		arg.CustomerAddressID,
		arg.Name,
		arg.LastName,
		arg.Email,
		arg.PhoneNumber,
	)
	var customer_id uuid.UUID
	err := row.Scan(&customer_id)
	return customer_id, err
}

const createCustomerAddress = `-- name: CreateCustomerAddress :one
INSERT INTO addresses (
    country, street, city, zip_code
) VALUES (
    $1, $2, $3, $4
) RETURNING address_id
`

type CreateCustomerAddressParams struct {
	Country string
	Street  string
	City    string
	ZipCode string
}

func (q *Queries) CreateCustomerAddress(ctx context.Context, arg CreateCustomerAddressParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createCustomerAddress,
		arg.Country,
		arg.Street,
		arg.City,
		arg.ZipCode,
	)
	var address_id uuid.UUID
	err := row.Scan(&address_id)
	return address_id, err
}

const getAllCustomers = `-- name: GetAllCustomers :many
SELECT customer_id, customer_address_id, name, last_name, email, phone_number, address_id, country, street, city, zip_code FROM customers JOIN addresses a ON a.address_id = customers.customer_address_id
ORDER BY customer_id
`

type GetAllCustomersRow struct {
	CustomerID        uuid.UUID
	CustomerAddressID uuid.UUID
	Name              string
	LastName          string
	Email             string
	PhoneNumber       string
	AddressID         uuid.UUID
	Country           string
	Street            string
	City              string
	ZipCode           string
}

func (q *Queries) GetAllCustomers(ctx context.Context) ([]GetAllCustomersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllCustomers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllCustomersRow
	for rows.Next() {
		var i GetAllCustomersRow
		if err := rows.Scan(
			&i.CustomerID,
			&i.CustomerAddressID,
			&i.Name,
			&i.LastName,
			&i.Email,
			&i.PhoneNumber,
			&i.AddressID,
			&i.Country,
			&i.Street,
			&i.City,
			&i.ZipCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
