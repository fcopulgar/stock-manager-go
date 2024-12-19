package services

import (
	"time"

	"github.com/fcopulgar/stock-manager-go/models"
)

type PortfolioServiceInterface interface {
	GetAllPortfolios() ([]models.Portfolio, error)
	GetPortfolioByID(id int) (*models.Portfolio, error)
	CreatePortfolioManual(portfolio *models.Portfolio) error
	DeletePortfolio(id int) error
	CalculateAPR(portfolio *models.Portfolio, startDate, endDate time.Time) (float64, error)
	GetPriceClose(symbol string, date time.Time) (float64, error)
	GetSP500Symbols() ([]string, error)
}
