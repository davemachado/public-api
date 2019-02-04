package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
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
		API         string `json:"API"`
		Description string `json:"Description"`
		Auth        string `json:"Auth"`
		HTTPS       bool   `json:"HTTPS"`
		Cors        string `json:"Cors"`
		Link        string `json:"Link"`
		Category    string `json:"Category"`
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
		results, err := processSearchRequestToMatchingEntries(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
		results, err := processSearchRequestToMatchingEntries(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err = json.NewEncoder(w).Encode(Entries{
			Count:   1,
			Entries: []Entry{results[rand.Intn(len(results))]},
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
