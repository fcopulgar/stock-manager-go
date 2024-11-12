package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// MockHTTPClient is a mock implementation of HTTPClient.
type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

func TestGetSP500Symbols_Success(t *testing.T) {
	// Mock CSV data representing a subset of the S&P 500
	mockCSV := `Symbol,Name,Sector
AAPL,Apple Inc,Information Technology
MSFT,Microsoft Corporation,Information Technology
AMZN,Amazon.com Inc,Consumer Discretionary
GOOGL,Alphabet Inc Class A,Communication Services
FB,Meta Platforms Inc Class A,Communication Services`

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(mockCSV)),
	}

	mockClient := &MockHTTPClient{
		Response: mockResponse,
		Err:      nil,
	}

	symbols, err := GetSP500Symbols(mockClient)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSymbols := []string{"AAPL", "MSFT", "AMZN", "GOOGL", "FB"}

	if len(symbols) != len(expectedSymbols) {
		t.Fatalf("Expected %d symbols, got %d", len(expectedSymbols), len(symbols))
	}

	for i, symbol := range symbols {
		if symbol != expectedSymbols[i] {
			t.Errorf("Expected symbol %s at index %d, got %s", expectedSymbols[i], i, symbol)
		}
	}
}

func TestGetSP500Symbols_HTTPError(t *testing.T) {
	mockClient := &MockHTTPClient{
		Response: nil,
		Err:      fmt.Errorf("network error"),
	}

	_, err := GetSP500Symbols(mockClient)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}
