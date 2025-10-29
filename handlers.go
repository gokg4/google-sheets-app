package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// rootHandler fetches the sheet data and serves the HTML page
func rootHandler(w http.ResponseWriter, r *http.Request) {
	sheetID := "0"
	spreadsheetID, ok := config["spreadsheetID"]
	if !ok {
		log.Fatal("FATAL: spreadsheetID not found in config")
	}

	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&gid=%s", spreadsheetID, sheetID)

	// Fetch the data
	headers, rows, err := fetchSheetData(url)

	// Prepare the data for the template
	var data PageData
	if err != nil {
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
