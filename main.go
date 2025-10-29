package main

import (
	"log"
	"net/http"
)

// Data structure for the HTML template
type PageData struct {
	Headers []string
	Rows    [][]string
	Error   string
}

func main() {
	// Load configuration from config.csv
	if err := loadConfig("config.csv"); err != nil {
		log.Fatalf("FATAL: Could not load config.csv: %v", err)
	}

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
