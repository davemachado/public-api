package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
)

func encodeURL(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	encoded := strings.Replace(r.URL.String(), "%20&%20", "%20%26%20", -1)
	r.URL, _ = url.Parse(encoded)
	next(rw, r)
}

// getList initializes an Entries struct filled from the public-apis project
func getList(jsonFile string) {
	file, err := os.OpenFile(jsonFile, os.O_RDONLY, 0644)
	if err != nil {
		panic("failed to open file: " + err.Error())
	}

	err = json.NewDecoder(file).Decode(&apiList)
	if err != nil {
		panic("failed to decode JSON from file: " + err.Error())
	}
	file.Close()
}

// getCategories initializes a string slice containing
// all unique categories from a given slice of Entries
func parseCategories(entries []Entry) []string {
	var cats []string
	set := make(map[string]struct{})
	for _, entry := range entries {
		if _, exists := set[entry.Category]; !exists {
			cats = append(cats, entry.Category)
			set[entry.Category] = struct{}{}
		}
	}
	return cats
}

// checkEntryMatches checks if the given entry matches the given request's parameters.
// it returns true if the entry matches, and returns false otherwise.
func checkEntryMatches(entry Entry, request *SearchRequest) bool {
	if strings.Contains(strings.ToLower(entry.API), strings.ToLower(request.Title)) &&
		strings.Contains(strings.ToLower(entry.Description), strings.ToLower(request.Description)) &&
		((strings.ToLower(request.Auth) == "null" && entry.Auth == "") || strings.Contains(strings.ToLower(entry.Auth), strings.ToLower(request.Auth))) &&
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

// processSearchRequestToMatchingEntries decodes the request body into a SearchRequest struct that can
// be processed in a call to checkEntryMatches to return all matching entries.
func processSearchRequestToMatchingEntries(req *http.Request) ([]Entry, error) {
	searchReq := new(SearchRequest)
	// Decode incoming search request off the query parameters map.
	if err := schema.NewDecoder().Decode(searchReq, req.URL.Query()); err != nil {
		return nil, fmt.Errorf("server failed to parse request: %v", err)
	}
	var results []Entry
	for _, e := range apiList.Entries {
		if checkEntryMatches(e, searchReq) {
			results = append(results, e)
		}
	}
	return results, nil
}
