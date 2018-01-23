package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_negroni"
	"github.com/urfave/negroni"
)

var apiList Entries

func main() {
	file, err := os.OpenFile("./entries.min.json", os.O_RDONLY, 0644)
	if err != nil {
		panic("failed to open entries.min.json: " + err.Error())
	}

	err = json.NewDecoder(file).Decode(&apiList)
	if err != nil {
		panic("failed to decode JSON from file: " + err.Error())
	}
	file.Close()

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

	n := negroni.New()
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n.Use(recovery)
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", n))
}
