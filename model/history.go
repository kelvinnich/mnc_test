package model

import "time"

type History struct {
	CustomerID string `json:"customer_id"`
	MerchantID string `json:"merchant_id"`
	Amount     int    `json:"amount"`
	Time time.Time `json:"date"`
}