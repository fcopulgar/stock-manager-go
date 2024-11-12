package services

import (
	"time"
)

type StockService interface {
	GetPriceOpen(symbol string, date time.Time) (float64, error)
	GetPriceClose(symbol string, date time.Time) (float64, error)
}
