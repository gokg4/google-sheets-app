package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// PageData defines the data structure for the HTML template.
type PageData struct {
	Headers []string
	Rows    [][]string
	Error   string
}

func main() {
	// Get the spreadsheet ID from the environment variable
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		log.Fatal("FATAL: SPREADSHEET_ID environment variable not set")
	}

	// Hardcode the sheet ID (gid)
	sheetID := "0"

	// Construct the sheet URL
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&gid=%s", spreadsheetID, sheetID)

	// Fetch and parse the data from the Google Sheet
	headers, rows, err := fetchSheetData(url)

	// Prepare the data for the template
	data := PageData{
		Headers: headers,
		Rows:    rows,
	}
	if err != nil {
		data.Error = fmt.Sprintf("Could not load data: %v", err)
	}

	// Create the public directory if it doesn't exist
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		log.Fatalf("FATAL: could not create public directory: %v", err)
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/index.html.tmpl")
	if err != nil {
		log.Fatalf("FATAL: could not parse template: %v", err)
	}

	// Create the output file
	file, err := os.Create("public/index.html")
	if err != nil {
		log.Fatalf("FATAL: could not create output file: %v", err)
	}
	defer file.Close()

	// Execute the template and write to the file
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatalf("FATAL: could not execute template: %v", err)
	}

	log.Println("Successfully generated static site to public/index.html")
}

// fetchSheetData fetches and parses CSV data from the given URL.
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

	return records[0], records[1:], nil
}
