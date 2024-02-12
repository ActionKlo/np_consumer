// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: payloads.sql

package gen

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const getPayloads = `-- name: GetPayloads :many
SELECT message_id, receiver_id, tracking_number, event_id, event_type, event_time, data FROM payloads
`

func (q *Queries) GetPayloads(ctx context.Context) ([]Payload, error) {
	rows, err := q.db.QueryContext(ctx, getPayloads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Payload
	for rows.Next() {
		var i Payload
		if err := rows.Scan(
			&i.MessageID,
			&i.ReceiverID,
			&i.TrackingNumber,
			&i.EventID,
			&i.EventType,
			&i.EventTime,
			&i.Data,
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

const getPayloadsWithSettings = `-- name: GetPayloadsWithSettings :many
SELECT message_id, tracking_number, event_type, event_time, data, s.url
FROM payloads JOIN receivers s ON s.receiver_id = payloads.receiver_id
`

type GetPayloadsWithSettingsRow struct {
	MessageID      uuid.UUID
	TrackingNumber string
	EventType      string
	EventTime      time.Time
	Data           json.RawMessage
	Url            string
}

func (q *Queries) GetPayloadsWithSettings(ctx context.Context) ([]GetPayloadsWithSettingsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPayloadsWithSettings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPayloadsWithSettingsRow
	for rows.Next() {
		var i GetPayloadsWithSettingsRow
		if err := rows.Scan(
			&i.MessageID,
			&i.TrackingNumber,
			&i.EventType,
			&i.EventTime,
			&i.Data,
			&i.Url,
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

const savePayload = `-- name: SavePayload :one
INSERT INTO payloads (
    message_id, tracking_number, event_id, event_type, event_time, data, receiver_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING message_id
`

type SavePayloadParams struct {
	MessageID      uuid.UUID
	TrackingNumber string
	EventID        uuid.UUID
	EventType      string
	EventTime      time.Time
	Data           json.RawMessage
	ReceiverID     uuid.UUID
}

func (q *Queries) SavePayload(ctx context.Context, arg SavePayloadParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, savePayload,
		arg.MessageID,
		arg.TrackingNumber,
		arg.EventID,
		arg.EventType,
		arg.EventTime,
		arg.Data,
		arg.ReceiverID,
	)
	var message_id uuid.UUID
	err := row.Scan(&message_id)
	return message_id, err
}
