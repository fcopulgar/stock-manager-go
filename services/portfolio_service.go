package services

import (
	"math"
	"time"

	"github.com/fcopulgar/stock-manager-go/models"
	"github.com/fcopulgar/stock-manager-go/repositories"
)

type PortfolioService struct {
	Repo         repositories.PortfolioRepository
	StockService StockServiceInterface
}

func NewPortfolioService(repo repositories.PortfolioRepository, stockService StockServiceInterface) *PortfolioService {
	return &PortfolioService{
		Repo:         repo,
		StockService: stockService,
	}
}

func (ps *PortfolioService) GetAllPortfolios() ([]models.Portfolio, error) {
	return ps.Repo.GetAll()
}

func (ps *PortfolioService) GetPortfolioByID(id int) (*models.Portfolio, error) {
	return ps.Repo.GetByID(id)
}

func (ps *PortfolioService) CreatePortfolioManual(portfolio *models.Portfolio) error {
	return ps.Repo.Save(portfolio)
}

func (ps *PortfolioService) DeletePortfolio(id int) error {
	return ps.Repo.Delete(id)
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

func (ps *PortfolioService) GetPriceClose(symbol string, date time.Time) (float64, error) {
	return ps.StockService.GetPriceClose(symbol, date)
}

func (ps *PortfolioService) GetSP500Symbols() ([]string, error) {
	return ps.StockService.GetSP500Symbols()
}
