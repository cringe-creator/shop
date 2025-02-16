package models

import "time"

type Purchase struct {
	ItemName  string    `json:"item_name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
