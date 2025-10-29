package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

// fetchSheetData fetches and parses the CSV data from the given URL
func fetchSheetData(url string) ([]string, [][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch sheet data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse CSV: %w", err)
	}

	if len(records) < 1 {
		return nil, nil, fmt.Errorf("no data found in sheet")
	}

	headers := records[0]
	rows := records[1:]

	return headers, rows, nil
}
