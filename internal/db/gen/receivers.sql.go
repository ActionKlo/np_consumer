// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: receivers.sql

package gen

import (
	"context"

	"github.com/google/uuid"
)

const createReceiver = `-- name: CreateReceiver :one

INSERT INTO receivers (
    receiver_id, url
) VALUES (
    $1, $2
) RETURNING receiver_id
`

type CreateReceiverParams struct {
	ReceiverID uuid.UUID
	Url        string
}

// ----------------------------------
func (q *Queries) CreateReceiver(ctx context.Context, arg CreateReceiverParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createReceiver, arg.ReceiverID, arg.Url)
	var receiver_id uuid.UUID
	err := row.Scan(&receiver_id)
	return receiver_id, err
}

const deleteReceiver = `-- name: DeleteReceiver :exec
DELETE FROM receivers WHERE receiver_id = $1
`

func (q *Queries) DeleteReceiver(ctx context.Context, receiverID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteReceiver, receiverID)
	return err
}

const getSettingsByReceiverID = `-- name: GetSettingsByReceiverID :one
SELECT url FROM receivers where receiver_id = $1
`

func (q *Queries) GetSettingsByReceiverID(ctx context.Context, receiverID uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, getSettingsByReceiverID, receiverID)
	var url string
	err := row.Scan(&url)
	return url, err
}

const retrieveReceiver = `-- name: RetrieveReceiver :one
SELECT receiver_id, url FROM receivers WHERE receiver_id = $1
`

func (q *Queries) RetrieveReceiver(ctx context.Context, receiverID uuid.UUID) (Receiver, error) {
	row := q.db.QueryRowContext(ctx, retrieveReceiver, receiverID)
	var i Receiver
	err := row.Scan(&i.ReceiverID, &i.Url)
	return i, err
}

const saveSettings = `-- name: SaveSettings :one
INSERT INTO receivers (
    receiver_id, url
) VALUES (
    $1, $2
) RETURNING receiver_id
`

type SaveSettingsParams struct {
	ReceiverID uuid.UUID
	Url        string
}

func (q *Queries) SaveSettings(ctx context.Context, arg SaveSettingsParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, saveSettings, arg.ReceiverID, arg.Url)
	var receiver_id uuid.UUID
	err := row.Scan(&receiver_id)
	return receiver_id, err
}
