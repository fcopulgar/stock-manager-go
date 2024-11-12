package services

import (
	"github.com/fcopulgar/stock-manager-go/models"
	"github.com/fcopulgar/stock-manager-go/repositories"
	"math"
	"time"
)

type PortfolioService struct {
	Repo         repositories.PortfolioRepository
	StockService StockService
}

func NewPortfolioService(repo repositories.PortfolioRepository, stockService StockService) *PortfolioService {
	return &PortfolioService{
		Repo:         repo,
		StockService: stockService,
	}
}

func (ps *PortfolioService) CalculateAPR(portfolio *models.Portfolio, startDate, endDate time.Time) (float64, error) {
	initialValue := 0.0
	finalValue := 0.0

	for _, stock := range portfolio.Stocks {
		initialPrice := stock.BuyPrice

		finalPrice, err := ps.StockService.GetPriceClose(stock.Symbol, endDate)
		if err != nil {
			return 0, err
		}

		initialValue += initialPrice * float64(stock.Quantity)
		finalValue += finalPrice * float64(stock.Quantity)
	}

	years := endDate.Sub(startDate).Hours() / (24 * 365)

	if years == 0 {
		return 0, nil
	}

	apr := math.Pow(finalValue/initialValue, 1/years) - 1

	return apr, nil
}
