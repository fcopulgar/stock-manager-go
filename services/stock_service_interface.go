package services

import (
	"time"
)

type StockServiceInterface interface {
	GetPriceOpen(symbol string, date time.Time) (float64, error)
	GetPriceClose(symbol string, date time.Time) (float64, error)
}
