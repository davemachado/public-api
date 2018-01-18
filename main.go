package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

	http.HandleFunc("/api", getEntriesHandler)
	http.HandleFunc("/health-check", healthCheckHandler)

	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
