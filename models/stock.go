package models

import "time"

type Stock struct {
	Symbol   string
	Quantity int
	BuyDate  time.Time
}
