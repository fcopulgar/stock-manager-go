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
}

func NewFinancialModelingPrepService(apiKey string) *FinancialModelingPrepService {
	client := resty.New()
	return &FinancialModelingPrepService{
		APIKey: apiKey,
		Client: client,
	}
}

func (fmp *FinancialModelingPrepService) GetPrice(symbol string, date time.Time) (float64, error) {
	startDateStr := date.Format("2006-01-02")

	endDate := date.AddDate(0, 0, 1)
	endDateStr := endDate.Format("2006-01-02")

	url := fmt.Sprintf("https://financialmodelingprep.com/api/v3/historical-price-full/%s?from=%s&to=%s&apikey=%s",
		symbol, startDateStr, endDateStr, fmp.APIKey)

	fmt.Println(url)

	resp, err := fmp.Client.R().Get(url)
	if err != nil {
		return 0, err
	}

	// Parse the response
	var result struct {
		Symbol     string `json:"symbol"`
		Historical []struct {
			Date  string  `json:"date"`
			Close float64 `json:"close"`
		} `json:"historical"`
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return 0, err
	}

	if len(result.Historical) == 0 {
		return 0, fmt.Errorf("no price data available for %s on %s", symbol, startDateStr)
	}

	price := result.Historical[0].Close
	fmt.Println(price)
	return price, nil
}
