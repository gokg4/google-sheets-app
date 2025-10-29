package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Data structure for the HTML template
type PageData struct {
	Headers []string
	Rows    [][]string
	Error   string // Add an Error field
}

func main() {
	// The handler for the root path
	http.HandleFunc("/", rootHandler)

	// Add a file server for the "public" directory
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Start the server
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// rootHandler fetches the sheet data and serves the HTML page
func rootHandler(w http.ResponseWriter, r *http.Request) {
	spreadsheetID := "1-k5N3-G-p2E92s6s-O0s-a-K2FwHYs3hp2pDEsJjCR4"
	sheetID := "0"
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&gid=%s", spreadsheetID, sheetID)

	// Fetch the data
	headers, rows, err := fetchSheetData(url)

	// Prepare the data for the template
	var data PageData
	if err != nil {
		// If there's an error, populate the Error field
		data.Error = fmt.Sprintf("Could not load data: %v", err)
	} else {
		data.Headers = headers
		data.Rows = rows
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("public/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse template: %v", err), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

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
