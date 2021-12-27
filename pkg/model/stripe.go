package model

import "time"

type StripeCustomer struct {
	UserID     int64
	CustomerID string
	CreatedAt  time.Time
}
