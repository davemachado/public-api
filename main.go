package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.HandleFunc("/api", getEntriesHandler)
	r.HandleFunc("/health-check", healthCheckHandler)

	n := negroni.New()
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n.Use(recovery)
	n.Use(negroni.NewLogger())
	n.UseHandler(r)

	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", n))
}
