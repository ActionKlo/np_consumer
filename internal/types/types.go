package types

import "time"

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
