// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package gen

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	MessageID      uuid.UUID
	ReceiverID     uuid.UUID
	TrackingNumber string
	EventID        uuid.UUID
	EventType      string
	EventTime      time.Time
	Data           json.RawMessage
}

type Receiver struct {
	ReceiverID uuid.UUID
	Url        string
}
