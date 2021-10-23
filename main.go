package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	bind_addr := os.Getenv("BIND_ADDR")
	if bind_addr == "" {
		bind_addr = ":8080"
	}
	symbol := os.Getenv("SYMBOL")
	if symbol == "" {
		panic("The SYMBOL environment variable must be set")
	}
	ndays := os.Getenv("NDAYS")
	if ndays == "" {
		panic("The NDAYS environment variable must be set")
	}
	apikey := os.Getenv("APIKEY")
	if apikey == "" {
		panic("The APIKEY environment variable must be set")
	}

	log.Printf("Starting web server on %s", bind_addr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "symbol=%s, ndays=%s", symbol, ndays)
	})
	log.Fatal(http.ListenAndServe(bind_addr, nil))
}
