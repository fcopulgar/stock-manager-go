package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type FinancialModelingPrepService struct {
	APIKey string
	Client *resty.Client
	cache  map[string]map[string]StockPrices
}

type StockPrices struct {
	Open  float64
	Close float64
}

func NewFinancialModelingPrepService(apiKey string) *FinancialModelingPrepService {
	client := resty.New()
	return &FinancialModelingPrepService{
		APIKey: apiKey,
		Client: client,
		cache:  make(map[string]map[string]StockPrices),
	}
}

func (fmp *FinancialModelingPrepService) GetPriceOpen(symbol string, date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")

	// local cache
	if prices, ok := fmp.cache[symbol]; ok {
		if stockPrices, ok := prices[dateStr]; ok {
			return stockPrices.Open, nil
		}
	}

	// Fetch the price data from the API
	stockPrices, err := fmp.fetchStockPrices(symbol, date)
	if err != nil {
		return 0, err
	}

	return stockPrices.Open, nil
}

func (fmp *FinancialModelingPrepService) GetPriceClose(symbol string, date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")

	// local cache
	if prices, ok := fmp.cache[symbol]; ok {
		if stockPrices, ok := prices[dateStr]; ok {
			return stockPrices.Close, nil
		}
	}

	// Fetch the price data from the API
	stockPrices, err := fmp.fetchStockPrices(symbol, date)
	if err != nil {
		return 0, err
	}

	return stockPrices.Close, nil
}

func (fmp *FinancialModelingPrepService) fetchStockPrices(symbol string, date time.Time) (StockPrices, error) {
	dateStr := date.Format("2006-01-02")
	endDateStr := date.AddDate(0, 0, 1).Format("2006-01-02")

	url := fmt.Sprintf("https://financialmodelingprep.com/api/v3/historical-price-full/%s?from=%s&to=%s&apikey=%s",
		symbol, dateStr, endDateStr, fmp.APIKey)

	resp, err := fmp.Client.R().Get(url)

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

	// Store in cache
	if _, ok := fmp.cache[symbol]; !ok {
		fmp.cache[symbol] = make(map[string]StockPrices)
	}
	fmp.cache[symbol][dateStr] = stockPrices

	return stockPrices, nil
}
