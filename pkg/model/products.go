package model

import "time"

type Product struct {
	ID        int64	`json:"id"`
	Name      string	`json:"name"`
	Price     float64	`json:"price"`
	Currency  string `json:"currency"`
	PriceID   string	`json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
