package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_negroni"
	"github.com/urfave/negroni"
)

var apiList Entries

// getList returns an Entries struct filled from the public-apis project
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

func main() {
	jsonFile := os.Getenv("JSONFILE")
	if jsonFile == "" {
		jsonFile = "/entries.json"
	}
	getList(jsonFile)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	rate := os.Getenv("RATE")
	if rate == "" {
		rate = "10"
	}
	i, err := strconv.Atoi(rate)
	if err != nil {
		panic(err)
	}
	limiter := tollbooth.NewLimiter(int64(i), nil)

	filename := os.Getenv("LOGFILE")
	if filename == "" {
		filename = "/tmp/public-api.log"
	}
	f, _ := os.Create(filename)
	logger := NewLogger(Options{
		Out: io.MultiWriter(f, os.Stdout),
	})

	mux := http.NewServeMux()
	mux.Handle("/entries", negroni.New(
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
	n.Use(negroni.HandlerFunc(logger.logFunc))
	n.UseHandler(mux)

	log.Println("logging requests in " + filename)
	log.Printf("listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, n))
}
