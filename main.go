package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_negroni"
	"github.com/urfave/negroni"
)

const (
	jsonURL = "https://raw.githubusercontent.com/toddmotto/public-apis/master/json/entries.min.json"
)

var apiList Entries

// getList returns an Entries struct filled from the public-apis project
func getList() Entries {
	res, err := http.Get(jsonURL)
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
	mux := http.NewServeMux()

	rate := os.Getenv("RATE")
	if rate == "" {
		log.Fatal("$RATE not set")
	}
	i, err := strconv.Atoi(rate)
	if err != nil {
		panic(err)
	}
	limiter := tollbooth.NewLimiter(int64(i), nil)

	mux.Handle("/", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(getEntriesHandler()),
	))
	mux.Handle("/health-check", negroni.New(
		tollbooth_negroni.LimitHandler(limiter),
		negroni.Wrap(healthCheckHandler()),
	))

	filename := os.Getenv("LOGFILE")
	if filename == "" {
		log.Fatal("$LOGFILE not set")
	}
	f, _ := os.Create(filename)
	logger := NewLogger(Options{
		Out: io.MultiWriter(f, os.Stdout),
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT not set")
	}

	n := negroni.New()
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n.Use(recovery)
	n.Use(negroni.HandlerFunc(logger.logFunc))
	n.UseHandler(mux)

	apiList = getList()

	log.Println("logging requests in " + filename)
	log.Printf("listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, n))
}
