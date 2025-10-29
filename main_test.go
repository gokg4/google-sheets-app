package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchSheetData(t *testing.T) {
	// A mock server that will act as our fake Google Sheet.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is the CSV content our test server will provide.
		w.Header().Set("Content-Type", "text/csv")
		// Correctly formatted CSV with newlines
		w.Write([]byte(`Header1,Header2
Row1Col1,Row1Col2
Row2Col1,Row2Col2`))
	}))
	defer server.Close()

	// The URL from our mock server is passed to the function under test.
	headers, rows, err := fetchSheetData(server.URL)

	// 1. Test for unexpected errors.
	if err != nil {
		t.Fatalf("fetchSheetData returned an unexpected error: %v", err)
	}

	// 2. Define the expected results.
	expectedHeaders := []string{"Header1", "Header2"}
	expectedRows := [][]string{
		{"Row1Col1", "Row1Col2"},
		{"Row2Col1", "Row2Col2"},
	}

	// 3. Compare the actual results with our expectations.
	if !reflect.DeepEqual(headers, expectedHeaders) {
		t.Errorf("expected headers %v, got %v", expectedHeaders, headers)
	}
	if !reflect.DeepEqual(rows, expectedRows) {
		t.Errorf("expected rows %v, got %v", expectedRows, rows)
	}
}

// Test for what happens when the server returns a non-200 status code.
func TestFetchSheetData_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, _, err := fetchSheetData(server.URL)

	if err == nil {
		t.Fatal("expected an error, but got none")
	}
}
