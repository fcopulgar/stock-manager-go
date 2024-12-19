package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fcopulgar/stock-manager-go/api"
	"github.com/go-resty/resty/v2"
)

type FinancialModelingPrepService struct {
	APIKey  string
	Client  *resty.Client
	BaseURL string
}

type StockPrices struct {
	Open  float64
	Close float64
}

func NewFinancialModelingPrepService(apiKey string) *FinancialModelingPrepService {
	client := resty.New()
	client.SetBaseURL("https://financialmodelingprep.com")
	return &FinancialModelingPrepService{
		APIKey: apiKey,
		Client: client,
	}
}

func (fmp *FinancialModelingPrepService) GetPriceOpen(symbol string, date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")

	stockPrices, err := fmp.fetchStockPrices(symbol, date)
	if err != nil {
		fmt.Printf("Failed to retrieve open price for %s on %s.\n", symbol, dateStr)
		return fmp.promptUserForPrice(symbol, date, "open")
	}

	return stockPrices.Open, nil
}

func (fmp *FinancialModelingPrepService) GetPriceClose(symbol string, date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")

	stockPrices, err := fmp.fetchStockPrices(symbol, date)
	if err != nil {
		fmt.Printf("Failed to retrieve close price for %s on %s. %s\n", symbol, dateStr, err)
		return fmp.promptUserForPrice(symbol, date, "close")
	}

	return stockPrices.Close, nil
}

func (fmp *FinancialModelingPrepService) GetSP500Symbols() ([]string, error) {
	return api.GetSP500Symbols(&api.DefaultHTTPClient{})
}

func (fmp *FinancialModelingPrepService) fetchStockPrices(symbol string, date time.Time) (StockPrices, error) {
	dateStr := date.Format("2006-01-02")
	endDateStr := date.AddDate(0, 0, 1).Format("2006-01-02")

	path := fmt.Sprintf("/api/v3/historical-price-full/%s?from=%s&to=%s&apikey=%s",
		symbol, dateStr, endDateStr, fmp.APIKey)

	resp, err := fmp.Client.R().Get(path)
	if err != nil {
		return StockPrices{}, err
	}

	if resp.IsError() {
		return StockPrices{}, fmt.Errorf("API request failed with status code %d", resp.StatusCode())
	}

	var result struct {
		Symbol     string `json:"symbol"`
		Historical []struct {
			Date  string  `json:"date"`
			Open  float64 `json:"open"`
			Close float64 `json:"close"`
		} `json:"historical"`
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return StockPrices{}, err
	}

	if len(result.Historical) == 0 {
		return StockPrices{}, fmt.Errorf("no price data available for %s on %s", symbol, dateStr)
	}

	historicalData := result.Historical[0]
	stockPrices := StockPrices{
		Open:  historicalData.Open,
		Close: historicalData.Close,
	}

	return stockPrices, nil
}

func (fmp *FinancialModelingPrepService) promptUserForPrice(symbol string, date time.Time, priceType string) (float64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Please enter the %s price for %s on %s: ", priceType, symbol, date.Format("2006-01-02"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	price, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price input")
	}

	return price, nil
}
