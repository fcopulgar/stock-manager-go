package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

// Simulated structure for API response
type fmpHistoricalResponse struct {
	Symbol     string `json:"symbol"`
	Historical []struct {
		Date  string  `json:"date"`
		Open  float64 `json:"open"`
		Close float64 `json:"close"`
	} `json:"historical"`
}

func TestFinancialModelingPrepService_GetPriceOpen(t *testing.T) {
	// Create a test server that simulates the Financial Modeling Prep response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We simulate a response with a history with one day
		resp := fmpHistoricalResponse{
			Symbol: "AAPL",
			Historical: []struct {
				Date  string  "json:\"date\""
				Open  float64 "json:\"open\""
				Close float64 "json:\"close\""
			}{
				{Date: "2020-01-15", Open: 300.0, Close: 305.0},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	fmp := &FinancialModelingPrepService{
		APIKey: "dummykey",
		Client: resty.New(),
	}

	// Configure the base URL of the Resty client to the test server
	fmp.Client.SetBaseURL(ts.URL)

	// Test GetPriceOpen
	date := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	price, err := fmp.GetPriceOpen("AAPL", date)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if price != 300.0 {
		t.Errorf("Expected open price 300.0, got %f", price)
	}
}

func TestFinancialModelingPrepService_GetPriceClose(t *testing.T) {
	// Create a test server that simulates the response for GetPriceClose
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmpHistoricalResponse{
			Symbol: "AAPL",
			Historical: []struct {
				Date  string  "json:\"date\""
				Open  float64 "json:\"open\""
				Close float64 "json:\"close\""
			}{
				{Date: "2020-01-15", Open: 300.0, Close: 305.0},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	fmp := &FinancialModelingPrepService{
		APIKey: "dummykey",
		Client: resty.New(),
	}

	fmp.Client.SetBaseURL(ts.URL)

	date := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	price, err := fmp.GetPriceClose("AAPL", date)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if price != 305.0 {
		t.Errorf("Expected close price 305.0, got %f", price)
	}
}

func TestFinancialModelingPrepService_GetPriceClose_APIError(t *testing.T) {
	// Simulate an error in the API (e.g. return 500)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	fmp := &FinancialModelingPrepService{
		APIKey: "dummykey",
		Client: resty.New(),
	}

	fmp.Client.SetBaseURL(ts.URL)

	date := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	_, err := fmp.GetPriceClose("AAPL", date)
	if err == nil {
		t.Fatalf("Expected an error due to API error, got none")
	}
}
