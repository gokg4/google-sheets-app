package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

var config = make(map[string]string)

// loadConfig loads key-value pairs from a CSV file into the config map
func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return fmt.Errorf("config file must have at least a header and one row")
	}

	// Assuming the first row is the header
	for _, record := range records[1:] {
		if len(record) == 2 {
			config[record[0]] = record[1]
		}
	}

	return nil
}
