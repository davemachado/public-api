package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func mockEntry() []Entry {
	return []Entry{
		{
			API:         "title",
			Description: "description",
			Auth:        "apiKey",
			HTTPS:       false,
			Cors:        "Cors",
			Link:        "link",
			Category:    "category",
		},
	}
}

func assertResponseValid(t *testing.T, body *bytes.Buffer, expected []Entry) {
	var resp Entries
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resp.Entries, expected) {
		t.Fatalf("handler returned wrong entry: got %v want %v",
			resp.Entries, expected)
	}

}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := healthCheckHandler()
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

func TestGetCategoriesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		categories   []string
		expectedBody string
	}{
		{[]string{}, "[]\n"},
		{[]string{"cat1"}, "[\"cat1\"]\n"},
		{[]string{"cat1", "cat2", "cat3"}, "[\"cat1\",\"cat2\",\"cat3\"]\n"},
	}
	for _, tc := range testCases {
		categories = tc.categories
		rr := httptest.NewRecorder()
		handler := getCategoriesHandler()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if rr.Body.String() != tc.expectedBody {
			t.Errorf("handler returned wrong body: got %q want %q", rr.Body, tc.expectedBody)
		}
	}
}

func TestGetRandomHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/random", nil)
	if err != nil {
		t.Fatal(err)
	}
	apiList.Entries = mockEntry()
	rr := httptest.NewRecorder()
	handler := getRandomHandler()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	assertResponseValid(t, rr.Body, apiList.Entries)
}

func TestGetEntriesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	apiList.Entries = mockEntry()
	rr := httptest.NewRecorder()
	handler := getEntriesHandler()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	assertResponseValid(t, rr.Body, apiList.Entries)
}

func TestGetEntriesWithBadMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := getEntriesHandler()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
