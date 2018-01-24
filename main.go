package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_negroni"
	"github.com/urfave/negroni"
)

const jsonUrl = "https://raw.githubusercontent.com/toddmotto/public-apis/master/json/entries.min.json"

var apiList Entries

// getList returns an Entries struct filled from the public-apis project
func getList() Entries {
	res, err := http.Get(jsonUrl)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	entries := new(Entries)
	err = json.Unmarshal(body, &entries)
	if err != nil {
		panic(err)
	}
	return *entries
}

func main() {
	apiList = getList()
	mux := http.NewServeMux()

	limiter := tollbooth.NewLimiter(1, nil)

	mux.Handle("/api", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(getEntriesHandler()),
	))
	mux.Handle("/health-check", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(healthCheckHandler()),
	))

	f, _ := os.Create("requests.log")
	logWriter = io.MultiWriter(f, os.Stdout)

	n := negroni.New()
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(logger))
	n.UseHandler(mux)

	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", n))
}
