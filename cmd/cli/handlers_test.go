package cli

import (
	"bufio"
	"bytes"
	"github.com/fcopulgar/stock-manager-go/models"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

// TestViewPortfolios_NoPortfolios checks the behavior of viewPortfolios() when there are no portfolios.
func TestViewPortfolios_NoPortfolios(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Configure the mock so that GetAllPortfolios returns an empty list
	mockService.On("GetAllPortfolios").Return([]models.Portfolio{}, nil)

	// We don't need the user to enter anything in this case,
	// because viewPortfolios() just lists portfolios and then asks for an ID.
	// We'll simulate that the user just presses Enter to return.
	input := "\n" // User does not enter any ID, just Enter
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	// Create the CLI instance with the mock
	cli := CLI{
		portfolioService: mockService,
		reader:           bufio.NewReader(inputReader),
		writer:           &outputBuffer,
	}

	// Directly call viewPortfolios()
	cli.viewPortfolios()

	output := outputBuffer.String()

	// Verify that the output contains “No portfolios available."
	if !strings.Contains(output, "No portfolios available.") {
		t.Errorf("Expected output to contain 'No portfolios available.', got '%s'", output)
	}

	mockService.AssertExpectations(t)
}

// TestCreatePortfolioManual_NoInput tests that if the user does not enter actions, the corresponding message is displayed.
func TestCreatePortfolioManual_NoInput(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Configure the mock for GetSP500Symbols
	mockService.On("GetSP500Symbols").Return([]string{"AAPL", "MSFT"}, nil)

	// We simulate that the user enters the name of the portfolio (“My Portfolio”) and then presses Enter without selecting stocks.
	input := "My Portfolio\n\n"
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	cli := CLI{
		portfolioService: mockService,
		reader:           bufio.NewReader(inputReader),
		writer:           &outputBuffer,
	}

	cli.createPortfolioManual()

	output := outputBuffer.String()

	// Verify that the output contains “No stocks added to the portfolio.”
	if !strings.Contains(output, "No stocks added to the portfolio.") {
		t.Errorf("Expected output to contain 'No stocks added to the portfolio.', got '%s'", output)
	}

	// Since the user did not add shares, CreatePortfolioManual should not be called.
	mockService.AssertNotCalled(t, "CreatePortfolioManual", mock.Anything)
	mockService.AssertExpectations(t)
}

// TestCreatePortfolioManual_SingleStock tests that if the user adds a stock and a valid date, the portfolio is created.
func TestCreatePortfolioManual_SingleStock(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Configure the mocks
	mockService.On("GetSP500Symbols").Return([]string{"AAPL"}, nil)
	// When GetPriceClose is called with AAPL and a date, we will return a fixed price.
	mockService.On("GetPriceClose", "AAPL", mock.AnythingOfType("time.Time")).Return(300.0, nil)
	mockService.On("CreatePortfolioManual", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Simulate entry:
	// Portfolio name: “Test Portfolio”.
	// Select stock 1 (AAPL)
	// Quantity: 10
	// Purchase date: 2020-01-15
	// Press Enter to end selection
	input := "Test Portfolio\n1\n10\n2020-01-15\n\n"
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	cli := CLI{
		portfolioService: mockService,
		reader:           bufio.NewReader(inputReader),
		writer:           &outputBuffer,
	}

	cli.createPortfolioManual()

	output := outputBuffer.String()

	// Verify that it indicates “Portfolio created successfully”.
	if !strings.Contains(output, "Portfolio created successfully.") {
		t.Errorf("Expected output to contain 'Portfolio created successfully.', got '%s'", output)
	}

	// Verify that CreatePortfolioManual has been called with the expected portfolio
	mockService.AssertCalled(t, "CreatePortfolioManual", mock.MatchedBy(func(p *models.Portfolio) bool {
		if p.Name != "Test Portfolio" {
			return false
		}
		if len(p.Stocks) != 1 {
			return false
		}
		s := p.Stocks[0]
		if s.Symbol != "AAPL" || s.Quantity != 10 || s.BuyPrice != 300.0 {
			return false
		}
		return true
	}))

	mockService.AssertExpectations(t)
}
