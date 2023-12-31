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

	Customer struct {
		CustomerID  uuid.UUID
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

	Event struct {
		EventID          uuid.UUID
		EventTime        time.Time
		EventDescription string
	}

	Shipment struct {
		ShipmentID uuid.UUID
		Size       string
		Weight     float64
		Count      int

		Customer Customer
		Sender   Sender

		Event Event
	}
)
