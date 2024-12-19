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

// TestCLI_Run prueba que el menú se imprima y que la opción 4 salga sin errores.
func TestCLI_Run(t *testing.T) {
	mockService := new(MockPortfolioService)

	// Simular la entrada del usuario: selecciona opción 4 (Exit)
	input := "4\n"
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	cli := NewCLI(mockService, inputReader, &outputBuffer)
	cli.Run()

	output := outputBuffer.String()

	// Verificar que el menú esté presente en la salida
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

	// Verificar que al ingresar 4 se imprima "Exiting..."
	if !strings.Contains(output, "Exiting...") {
		t.Errorf("Expected output to contain 'Exiting...', got '%s'", output)
	}
}

// TestCLI_ViewPortfolios verifica que al seleccionar la opción 1 se llame a GetAllPortfolios()
// y se maneje el caso de no haber carteras.
func TestCLI_ViewPortfolios_NoPortfolios(t *testing.T) {
	mockService := new(MockPortfolioService)
	// Configurar el mock para que GetAllPortfolios devuelva una lista vacía
	mockService.On("GetAllPortfolios").Return([]models.Portfolio{}, nil)

	// Simular la entrada del usuario: opción 1 (ver carteras), luego 4 (salir)
	input := "1\n4\n"
	inputReader := strings.NewReader(input)
	var outputBuffer bytes.Buffer

	cli := NewCLI(mockService, inputReader, &outputBuffer)
	cli.Run()

	output := outputBuffer.String()

	// Verificar que llame a GetAllPortfolios al menos una vez
	mockService.AssertExpectations(t)

	// Verificar la salida
	if !strings.Contains(output, "No portfolios available.") {
		t.Errorf("Expected output to contain 'No portfolios available.', got '%s'", output)
	}
}
