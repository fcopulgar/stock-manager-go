package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type AlphaVantageService struct {
	APIKey string
	Client *resty.Client
}

func NewAlphaVantageService(apiKey string) *AlphaVantageService {
	client := resty.New()
	return &AlphaVantageService{
		APIKey: apiKey,
		Client: client,
	}
}

func (avs *AlphaVantageService) GetPrice(symbol string, date time.Time) (float64, error) {
	dateStr := date.Format("2006-01-02")
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&apikey=%s", symbol, avs.APIKey)

	resp, err := avs.Client.R().Get(url)
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return 0, err
	}

	timeSeries, ok := result["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("no price history found for %s", symbol)
	}

	dayData, ok := timeSeries[dateStr].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("no data found for the date %s", dateStr)
	}

	closeStr, ok := dayData["5. adjusted close"].(string)
	if !ok {
		return 0, fmt.Errorf("no adjusted closing price found for %s on date %s", symbol, dateStr)
	}

	var price float64
	fmt.Sscanf(closeStr, "%f", &price)

	return price, nil
}
