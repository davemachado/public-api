package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
)

// SearchRequest describes an incoming search request.
type SearchRequest struct {
	Title       string `schema:"title"`
	Description string `schema:"description"`
	Auth        string `schema:"auth"`
	HTTPS       string `schema:"https"`
	Cors        string `schema:"cors"`
	Category    string `schema:"category"`
}

// Entries contains an array of API entries, and a count representing the length of that array.
type Entries struct {
	Count   int     `json:"count"`
	Entries []Entry `json:"entries"`
}

// Entry describes a single API reference.
type Entry struct {
	API         string
	Description string
	Auth        string
	HTTPS       bool
	Cors        string
	Link        string
	Category    string
}

// checkEntryMatches checks if the given entry matches the given request's parameters.
// it returns true if the entry matches, and returns false otherwise.
func checkEntryMatches(entry Entry, request *SearchRequest) bool {
	if strings.Contains(strings.ToLower(entry.API), strings.ToLower(request.Title)) &&
		strings.Contains(strings.ToLower(entry.Description), strings.ToLower(request.Description)) &&
		strings.Contains(strings.ToLower(entry.Auth), strings.ToLower(request.Auth)) &&
		strings.Contains(strings.ToLower(entry.Cors), strings.ToLower(request.Cors)) &&
		strings.Contains(strings.ToLower(entry.Category), strings.ToLower(request.Category)) {
		if request.HTTPS == "" {
			return true
		}
		if value, err := strconv.ParseBool(request.HTTPS); err == nil {
			if entry.HTTPS == value {
				return true
			}
		}
	}
	return false
}

// getEntriesHandler returns an Entries object with the matching entries filtered
// by the search request
func getEntriesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Only allow GET requests
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var err error
		searchReq := new(SearchRequest)
		// Only check query parameters if the request's Body is not nil
		if req.Body != nil {
			// Decode incoming search request off the query parameters map.
			err = schema.NewDecoder().Decode(searchReq, req.URL.Query())
			if err != nil {
				http.Error(w, "server failed to parse request: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer req.Body.Close()
		}

		var results []Entry
		for _, e := range apiList.Entries {
			if checkEntryMatches(e, searchReq) {
				results = append(results, e)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err = json.NewEncoder(w).Encode(Entries{
			Count:   len(results),
			Entries: results,
		})
		if err != nil {
			http.Error(w, "server failed to encode response object: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// healthCheckHandler returns a simple indication on whether or not the core http service is running
func healthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, `{"alive": true}`)
	})
}
