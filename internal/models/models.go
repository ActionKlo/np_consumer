package models

import (
	"github.com/google/uuid"
	"time"
)

type (
	Address struct {
		AddressID uuid.UUID
		Country   string
		City      string
		Street    string
		Zip       string
	}

	Receiver struct {
		ReceiverID  uuid.UUID
		Name        string
		LastName    string
		Email       string
		PhoneNumber string
		Address     Address
	}
	Sender struct {
		SenderID    uuid.UUID
		Name        string
		LastName    string
		Email       string
		PhoneNumber string
		Address     Address
	}

	Order struct {
		OrderID uuid.UUID
		Size    string
		Weight  int
		Count   int

		Receiver Receiver
		Sender   Sender
	}

	Payload struct {
		MessageID      uuid.UUID
		EventID        uuid.UUID
		EventType      string
		EventTime      time.Time
		TrackingNumber string
		Order          Order
		ReceiverID     uuid.UUID
	}
)
