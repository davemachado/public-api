package main

import "testing"

func TestGetCategories(t *testing.T) {
	actual := parseCategories([]Entry{
		Entry{Category: "A"},
		Entry{Category: "B"},
		Entry{Category: "B"},
		Entry{Category: "C"},
		Entry{Category: "D"},
	})
	expected := []string{"A", "B", "C", "D"}
	if len(actual) != len(expected) {
		t.Fatalf("bad parsing: expected %v, got %v", expected, actual)
	}
	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("bad element: expected %q, got %q", actual[i], expected[i])
		}
	}
}

func TestCheckEntryMatches(t *testing.T) {
	entry := Entry{
		API:         "examplesAsAService",
		Description: "provide classic examples of classic things",
		Auth:        "apiKey",
		HTTPS:       true,
		Cors:        "Unknown",
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
	search.Cors = "unknown"
	if checkEntryMatches(entry, search) {
		t.Errorf("failed to match entry and search")
	}
}
