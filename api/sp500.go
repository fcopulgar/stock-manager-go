package api

import (
	"encoding/csv"
	"net/http"
)

func GetSP500Symbols() ([]string, error) {
	url := "https://raw.githubusercontent.com/datasets/s-and-p-500-companies/main/data/constituents.csv"

	resp, err := http.Get(url)
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
