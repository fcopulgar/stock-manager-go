package services

import (
	"testing"
	"time"

	"github.com/fcopulgar/stock-manager-go/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockPortfolioRepository es un mock de PortfolioRepository
type MockPortfolioRepository struct {
	mock.Mock
}

func (m *MockPortfolioRepository) GetAll() ([]models.Portfolio, error) {
	args := m.Called()
	return args.Get(0).([]models.Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) GetByID(id int) (*models.Portfolio, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) Save(portfolio *models.Portfolio) error {
	args := m.Called(portfolio)
	return args.Error(0)
}

func (m *MockPortfolioRepository) Update(portfolio *models.Portfolio) error {
	args := m.Called(portfolio)
	return args.Error(0)
}

func (m *MockPortfolioRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockStockService is a mock of StockServiceInterface
type MockStockService struct {
	mock.Mock
}

func (m *MockStockService) GetPriceOpen(symbol string, date time.Time) (float64, error) {
	args := m.Called(symbol, date)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockStockService) GetPriceClose(symbol string, date time.Time) (float64, error) {
	args := m.Called(symbol, date)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockStockService) GetSP500Symbols() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// TestGetAllPortfolios test TestGetAllPortfolios()
func TestGetAllPortfolios(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	mockStock := new(MockStockService)
	service := NewPortfolioService(mockRepo, mockStock)

	mockRepo.On("GetAll").Return([]models.Portfolio{
		{Name: "Test Portfolio"},
	}, nil)

	portfolios, err := service.GetAllPortfolios()
	require.NoError(t, err)
	require.Len(t, portfolios, 1)
	require.Equal(t, "Test Portfolio", portfolios[0].Name)

	mockRepo.AssertExpectations(t)
}

// TestGetPortfolioByID test TestGetPortfolioByID()
func TestGetPortfolioByID(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	mockStock := new(MockStockService)
	service := NewPortfolioService(mockRepo, mockStock)

	mockRepo.On("GetByID", 1).Return(&models.Portfolio{Name: "ID Portfolio"}, nil)

	p, err := service.GetPortfolioByID(1)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, "ID Portfolio", p.Name)

	mockRepo.AssertExpectations(t)
}

// TestCreatePortfolioManual test TestCreatePortfolioManual()
func TestCreatePortfolioManual(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	mockStock := new(MockStockService)
	service := NewPortfolioService(mockRepo, mockStock)

	portfolio := &models.Portfolio{Name: "New Portfolio"}
	mockRepo.On("Save", portfolio).Return(nil)

	err := service.CreatePortfolioManual(portfolio)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestDeletePortfolio test DeletePortfolio()
func TestDeletePortfolio(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	mockStock := new(MockStockService)
	service := NewPortfolioService(mockRepo, mockStock)

	mockRepo.On("Delete", 1).Return(nil)

	err := service.DeletePortfolio(1)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestCalculateAPR test CalculateAPR()
func TestCalculateAPR(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	mockStock := new(MockStockService)
	service := NewPortfolioService(mockRepo, mockStock)

	// Test portfolio
	startDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.UTC)

	portfolio := &models.Portfolio{
		Name: "APR Portfolio",
		Stocks: []models.Stock{
			{
				Symbol:   "AAPL",
				Quantity: 10,
				BuyDate:  startDate,
				BuyPrice: 100.0,
			},
		},
	}

	mockStock.On("GetPriceClose", "AAPL", endDate).Return(110.0, nil)

	apr, err := service.CalculateAPR(portfolio, startDate, endDate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedAPR := 0.1
	epsilon := 0.001 // tolerance of 0.1%

	if apr < expectedAPR-epsilon || apr > expectedAPR+epsilon {
		t.Errorf("APR expected %.4f, got %.4f (within ±%.4f)", expectedAPR, apr, epsilon)
	}

	mockStock.AssertExpectations(t)
}

// TestGetPriceClose test GetPriceClose()
func TestGetPriceClose(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	mockStock := new(MockStockService)
	service := NewPortfolioService(mockRepo, mockStock)

	date := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	mockStock.On("GetPriceClose", "AAPL", date).Return(300.0, nil)

	price, err := service.GetPriceClose("AAPL", date)
	require.NoError(t, err)
	require.Equal(t, 300.0, price)

	mockStock.AssertExpectations(t)
}

// TestGetSP500Symbols test GetSP500Symbols()
func TestGetSP500Symbols(t *testing.T) {
	mockRepo := new(MockPortfolioRepository)
	mockStock := new(MockStockService)
	service := NewPortfolioService(mockRepo, mockStock)

	mockStock.On("GetSP500Symbols").Return([]string{"AAPL", "MSFT"}, nil)

	symbols, err := service.GetSP500Symbols()
	require.NoError(t, err)
	require.Len(t, symbols, 2)
	require.Equal(t, "AAPL", symbols[0])
	require.Equal(t, "MSFT", symbols[1])

	mockStock.AssertExpectations(t)
}
