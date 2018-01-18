package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckEntryMatches(t *testing.T) {
	entry := Entry{
		API:         "examplesAsAService",
		Description: "provide classic examples of classic things",
		Auth:        "apiKey",
		HTTPS:       true,
		Link:        "http://www.example.com",
		Category:    "Development",
	}
	search := &SearchRequest{}
	if !checkEntryMatches(entry, search) {
		t.Errorf("failed to match entry and search")
	}
	search.HTTPS = "true"
	if !checkEntryMatches(entry, search) {
		t.Errorf("failed to match entry and search")
	}
	search.Auth = "OAuth"
	if checkEntryMatches(entry, search) {
		t.Errorf("failed to match entry and search")
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetEntriesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getEntriesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func TestGetEntriesWithBadMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getEntriesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
