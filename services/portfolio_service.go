package services

import (
	"fmt"
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
	var initialValue, finalValue float64
	var validStocks int

	for _, stock := range portfolio.Stocks {
		initialPrice, err := ps.StockService.GetPriceClose(stock.Symbol, startDate)
		if err != nil {
			// Log the error and skip this stock
			fmt.Printf("Could not retrieve initial price for %s (%v). Skipping...\n", stock.Symbol, err)
			continue
		}

		lastEndDateStr := endDate.AddDate(0, 0, -1)
		finalPrice, err := ps.StockService.GetPriceClose(stock.Symbol, lastEndDateStr)
		if err != nil {
			// Log the error and skip this stock
			fmt.Printf("Could not retrieve final price for %s (%v). Skipping...\n", stock.Symbol, err)
			continue
		}

		initialValue += initialPrice * float64(stock.Quantity)
		finalValue += finalPrice * float64(stock.Quantity)
		validStocks++
	}

	// Check if we have any valid stocks
	if validStocks == 0 || initialValue == 0 {
		return 0, fmt.Errorf("no valid stock prices available to calculate APR")
	}

	// Calculate the number of years between the dates
	years := endDate.Sub(startDate).Hours() / (24 * 365)

	if years == 0 {
		return 0, nil
	}

	// Calculate the annualized return
	apr := math.Pow(finalValue/initialValue, 1/years) - 1

	return apr, nil
}
