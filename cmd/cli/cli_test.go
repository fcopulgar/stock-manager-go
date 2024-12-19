package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/fcopulgar/stock-manager-go/models"
	"github.com/stretchr/testify/mock"
)

type MockPortfolioService struct {
	mock.Mock
}

func (m *MockPortfolioService) GetAllPortfolios() ([]models.Portfolio, error) {
	args := m.Called()
	return args.Get(0).([]models.Portfolio), args.Error(1)
}

func (m *MockPortfolioService) GetPortfolioByID(id int) (*models.Portfolio, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Portfolio), args.Error(1)
}

func (m *MockPortfolioService) CreatePortfolioManual(portfolio *models.Portfolio) error {
	args := m.Called(portfolio)
	return args.Error(0)
}

func (m *MockPortfolioService) DeletePortfolio(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPortfolioService) CalculateAPR(portfolio *models.Portfolio, startDate, endDate time.Time) (float64, error) {
	args := m.Called(portfolio, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPortfolioService) GetPriceClose(symbol string, date time.Time) (float64, error) {
	args := m.Called(symbol, date)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPortfolioService) GetSP500Symbols() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// TestCLI_Run tests that the menu prints and that option 4 exits without errors.
func TestCLI_Run(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Simulate user login: select option 4 (Exit)
	input := "4\n"
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	cli := NewCLI(mockService, inputReader, &outputBuffer)
	cli.Run()

	output := outputBuffer.String()

	// Verify that the menu is present at the output
	expectedMenuItems := []string{
		"Select an option:",
		"1. View portfolios",
		"2. Create portfolio manually",
		"3. Create random portfolio",
		"4. Exit",
	}

	for _, item := range expectedMenuItems {
		if !strings.Contains(output, item) {
			t.Errorf("Expected output to contain '%s', got '%s'", item, output)
		}
	}

	// Verify that when 4 is entered, “Exiting..." is printed.
	if !strings.Contains(output, "Exiting...") {
		t.Errorf("Expected output to contain 'Exiting...', got '%s'", output)
	}
}

// TestCLI_ViewPortfolios verifies that when option 1 is selected, GetAllPortfolios() is called.
// and the case of no portfolios is handled.
func TestCLI_ViewPortfolios_NoPortfolios(t *testing.T) {
	mockService := new(MockPortfolioService)
	// Configure the mock so that GetAllPortfolios returns an empty list
	mockService.On("GetAllPortfolios").Return([]models.Portfolio{}, nil)

	// Simulate user login: option 1 (view portfolios), then 4 (exit)
	input := "1\n4\n"
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	cli := NewCLI(mockService, inputReader, &outputBuffer)
	cli.Run()

	output := outputBuffer.String()

	// Verify that GetAllPortfolios is called at least once
	mockService.AssertExpectations(t)

	// Verify the output
	if !strings.Contains(output, "No portfolios available.") {
		t.Errorf("Expected output to contain 'No portfolios available.', got '%s'", output)
	}
}
