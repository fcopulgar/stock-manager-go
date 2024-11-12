package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

const sp500URL = "https://raw.githubusercontent.com/datasets/s-and-p-500-companies/main/data/constituents.csv"

// HTTPClient defines the interface for an HTTP client.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultHTTPClient is an implementation of HTTPClient using the standard library.
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// GetSP500Symbols fetches the S&P 500 symbols using the provided HTTP client.
func GetSP500Symbols(client HTTPClient) ([]string, error) {
	fmt.Println("Downloading: " + sp500URL)
	resp, err := client.Get(sp500URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var symbols []string
	for i, record := range records {
		if i == 0 {
			continue
		}
		symbols = append(symbols, record[0])
	}

	return symbols, nil
}
