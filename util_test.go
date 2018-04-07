package main

import "testing"

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
