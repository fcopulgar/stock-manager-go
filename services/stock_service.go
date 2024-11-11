package services

import (
	"time"
)

type StockService interface {
	GetPrice(symbol string, date time.Time) (float64, error)
}
