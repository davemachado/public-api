package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"

	"github.com/gorilla/schema"
)

type (
	// SearchRequest describes an incoming search request.
	SearchRequest struct {
		Title       string `schema:"title"`
		Description string `schema:"description"`
		Auth        string `schema:"auth"`
		HTTPS       string `schema:"https"`
		Cors        string `schema:"cors"`
		Category    string `schema:"category"`
	}
	// Entries contains an array of API entries, and a count representing the length of that array.
	Entries struct {
		Count   int     `json:"count"`
		Entries []Entry `json:"entries"`
	}
	// Entry describes a single API reference.
	Entry struct {
		API         string
		Description string
		Auth        string
		HTTPS       bool
		Cors        string
		Link        string
		Category    string
	}
)

// getEntriesHandler returns an Entries object with the matching entries filtered
// by the search request
func getEntriesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

// getCategoriesHandler returns a string slice object with all unique categories
func getCategoriesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := json.NewEncoder(w).Encode(categories)
		if err != nil {
			http.Error(w, "server failed to encode response object: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// getRandomHandler returns an Entries object containing a random element from the Entries slice
func getRandomHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := json.NewEncoder(w).Encode(Entries{
			Count:   1,
			Entries: []Entry{apiList.Entries[rand.Intn(len(apiList.Entries))]},
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
