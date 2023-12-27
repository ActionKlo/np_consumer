package models

import (
	"np_consumer/internal/db/gen"
	"time"
)

type KafkaMessage struct {
	ID          string    `json:"ID"`
	Time        time.Time `json:"Time"`
	Sender      string    `json:"Sender"`
	Status      string    `json:"Status"`
	TrackNumber string    `json:"TrackNumber"`
	Country     string    `json:"Country"`
	City        string    `json:"City"`
	Street      string    `json:"Street"`
	PostCode    string    `json:"PostCode"`
}

type CustomerInfo struct {
	Address  *gen.Address  `json:"*Gen.Address"`
	Customer *gen.Customer `json:"*Gen.Customer"`
}

type CustInf struct {
	*gen.Address  `json:"*Gen.Address"`
	*gen.Customer `json:"*Gen.Customer"`
}
