package main

import (
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
var categories []string

func main() {
	jsonFile := os.Getenv("JSONFILE")
	if jsonFile == "" {
		jsonFile = "/entries.json"
	}
	getList(jsonFile)
	categories = parseCategories(apiList.Entries)
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
	limiter := tollbooth.NewLimiter(float64(i), nil)

	filename := os.Getenv("LOGFILE")
	if filename == "" {
		filename = "/tmp/public-api.log"
	}
	// If the file does not exist, create it. Otherwise, append to the file.
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		panic(err)
	}
	logger := NewLogger(Options{
		Out: io.MultiWriter(f, os.Stdout),
	})

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.Handle("/entries", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(getEntriesHandler()),
	))
	mux.Handle("/categories", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(getCategoriesHandler()),
	))
	mux.Handle("/random", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(getRandomHandler()),
	))
	mux.Handle("/health", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(healthCheckHandler()),
	))

	n := negroni.New(negroni.HandlerFunc(logger.logFunc), negroni.HandlerFunc(encodeURL))
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n.Use(recovery)
	n.UseHandler(mux)

	log.Println("logging requests in " + filename)
	log.Printf("listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, n))
}
