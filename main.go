package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/pushittoprod/stock-quotes/alphavantage"
)

type EnvArgs struct {
	BindAddr string
	Symbol   string
	Ndays    int
}

var args EnvArgs

func init() {
	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = ":8080"
	}
	args.BindAddr = bindAddr

	symbol := os.Getenv("SYMBOL")
	if symbol == "" {
		panic("The SYMBOL environment variable must be set")
	}
	args.Symbol = symbol

	ndaysStr := os.Getenv("NDAYS")
	if ndaysStr == "" {
		panic("The NDAYS environment variable must be set")
	}

	ndays, err := strconv.Atoi(ndaysStr)
	if err != nil {
		panic("The NDAYS environment variable must be a valid integer")
	}
	args.Ndays = ndays
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
	log.Printf("Starting web server on %s", args.BindAddr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling request")

		data := createApiResponse(args.Symbol, args.Ndays)

		fmt.Fprintf(w, "%s data=[%.2f", data.symbol, data.data[0])
		for _, price := range data.data {
			fmt.Fprintf(w, ", %.2f", price)
		}
		fmt.Fprintf(w, "], average=%.2f", data.average)
	})

	log.Fatal(http.ListenAndServe(args.BindAddr, nil))
}
