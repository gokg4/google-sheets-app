package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
)

// PageData now only needs to hold the Spreadsheet ID for the template.
type PageData struct {
	SpreadsheetID string
}

func main() {
	// Get the spreadsheet ID from the environment variable
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		log.Fatal("FATAL: SPREADSHEET_ID environment variable not set")
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

	// Prepare the data for the template
	data := PageData{
		SpreadsheetID: spreadsheetID,
	}

	// Execute the template and write to the file
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatalf("FATAL: could not execute template: %v", err)
	}

	// Copy static assets
	if err := copyStaticAssets("static", "public"); err != nil {
		log.Fatalf("FATAL: could not copy static assets: %v", err)
	}

	log.Println("Successfully generated static site to public/index.html")
}

// copyStaticAssets copies all files and directories from src to dst.
func copyStaticAssets(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a parallel destination path
		dstPath := filepath.Join(dst, path[len(src):])

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// If it's a file, copy it
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
