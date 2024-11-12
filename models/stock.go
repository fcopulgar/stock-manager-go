package models

import "time"

type Stock struct {
	ID       int
	Symbol   string
	Quantity int
	BuyDate  time.Time
	BuyPrice float64
}
