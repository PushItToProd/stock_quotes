package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/pushittoprod/stock-quotes/alphavantage"
)

var (
	bindAddr  = os.Getenv("BIND_ADDR")
	symbol    = os.Getenv("SYMBOL")
	ndays_str = os.Getenv("NDAYS")
)
var ndays int

func init() {
	if bindAddr == "" {
		bindAddr = ":8080"
	}

	if symbol == "" {
		panic("The SYMBOL environment variable must be set")
	}

	if ndays_str == "" {
		panic("The NDAYS environment variable must be set")
	}
	ndays, err := strconv.Atoi(ndays_str)
	log.Printf("ndays_str=%s, ndays=%d", ndays_str, ndays)
	if err != nil {
		fmt.Errorf("Invalid value for NDAYS: %v", ndays)
		panic("The NDAYS environment variable must be a valid integer")
	}
}

type ApiResponse struct {
	symbol  string
	data    []float64 // in real code, I would not use float64 for currency
	average float64
}

func createApiResponse(symbol string, ndays int) *ApiResponse {
	data := alphavantage.GetClosingData(symbol, ndays)
	average := mean(data)
	return &ApiResponse{
		symbol:  symbol,
		data:    data,
		average: average,
	}
}

func mean(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func main() {
	log.Printf("Starting web server on %s", bindAddr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling request")
		log.Printf("Getting %d days of data for %s", ndays, symbol)
		data := createApiResponse(symbol, ndays)
		fmt.Fprintf(w, "%s, data=%v, average=%.2f", data.symbol, data.data, data.average)
	})

	log.Fatal(http.ListenAndServe(bindAddr, nil))
}
