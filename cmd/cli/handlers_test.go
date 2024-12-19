package cli

import (
	"bufio"
	"bytes"
	"github.com/fcopulgar/stock-manager-go/models"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

// TestViewPortfolios_NoPortfolios verifica el comportamiento de viewPortfolios() cuando no hay carteras.
func TestViewPortfolios_NoPortfolios(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Configurar el mock para que GetAllPortfolios devuelva una lista vacía
	mockService.On("GetAllPortfolios").Return([]models.Portfolio{}, nil)

	// No necesitamos que el usuario ingrese nada en este caso,
	// porque viewPortfolios() solo lista carteras y luego pide un ID.
	// Simularemos que el usuario solo presiona Enter para volver.
	input := "\n" // Usuario no ingresa ningún ID, solo Enter
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	// Crear la instancia de CLI con el mock
	cli := CLI{
		portfolioService: mockService,
		reader:           bufio.NewReader(inputReader),
		writer:           &outputBuffer,
	}

	// Llamar directamente a viewPortfolios()
	cli.viewPortfolios()

	output := outputBuffer.String()

	// Verificar que la salida contenga "No portfolios available."
	if !strings.Contains(output, "No portfolios available.") {
		t.Errorf("Expected output to contain 'No portfolios available.', got '%s'", output)
	}

	mockService.AssertExpectations(t)
}

// TestCreatePortfolioManual_NoInput prueba que si el usuario no ingresa acciones, se muestre el mensaje correspondiente.
func TestCreatePortfolioManual_NoInput(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Configurar el mock para GetSP500Symbols
	mockService.On("GetSP500Symbols").Return([]string{"AAPL", "MSFT"}, nil)

	// Simulamos que el usuario ingresa el nombre de la cartera ("My Portfolio") y luego presiona Enter sin seleccionar acciones.
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

	// Verificar que la salida contenga "No stocks added to the portfolio."
	if !strings.Contains(output, "No stocks added to the portfolio.") {
		t.Errorf("Expected output to contain 'No stocks added to the portfolio.', got '%s'", output)
	}

	// Dado que el usuario no agregó acciones, no se debe llamar a CreatePortfolioManual.
	mockService.AssertNotCalled(t, "CreatePortfolioManual", mock.Anything)
	mockService.AssertExpectations(t)
}

// TestCreatePortfolioManual_SingleStock prueba que si el usuario agrega una acción y una fecha válida, se cree la cartera.
func TestCreatePortfolioManual_SingleStock(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Configurar los mocks
	mockService.On("GetSP500Symbols").Return([]string{"AAPL"}, nil)
	// Cuando se llama a GetPriceClose con AAPL y una fecha, devolveremos un precio fijo.
	mockService.On("GetPriceClose", "AAPL", mock.AnythingOfType("time.Time")).Return(300.0, nil)
	mockService.On("CreatePortfolioManual", mock.AnythingOfType("*models.Portfolio")).Return(nil)

	// Simular entrada:
	// Nombre de la cartera: "Test Portfolio"
	// Seleccionar stock 1 (AAPL)
	// Cantidad: 10
	// Fecha de compra: 2020-01-15
	// Presionar Enter para terminar selección
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

	// Verificar que indique "Portfolio created successfully."
	if !strings.Contains(output, "Portfolio created successfully.") {
		t.Errorf("Expected output to contain 'Portfolio created successfully.', got '%s'", output)
	}

	// Verificar que CreatePortfolioManual se haya llamado con la cartera esperada
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
