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
		log.Fatalf("The NDAYS environment variable must be a valid integer - got '%v' instead", ndaysStr)
	}
	args.Ndays = ndays
}

type ApiResponse struct {
	symbol  string
	data    []float64 // in real code, I would not use float64 for currency
	average float64
}

func createApiResponse(symbol string, ndays int) (*ApiResponse, error) {
	data, err := alphavantage.GetClosingData(symbol, ndays)
	if err != nil {
		return nil, err
	}
	average := mean(data)
	resp := &ApiResponse{
		symbol:  symbol,
		data:    data,
		average: average,
	}
	return resp, nil
}

func mean(xs []float64) float64 {
	if len(xs) == 0 {
		panic("Can't take the mean of an empty slice")
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func main() {
	log.Printf("Starting web server on %s", args.BindAddr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// return 404 for requests to other paths, mainly to avoid spamming the alphavantage API
		// when the browser requests favicon.ico
		if r.URL.Path != "/" {
			log.Printf("Not found: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
			fmt.Fprintf(w, "404 not found")
			return
		}

		log.Printf("%s %s", r.Method, r.URL.Path)

		data, err := createApiResponse(args.Symbol, args.Ndays)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error")
			log.Printf("Error getting upstream API response: %v", err)
			return
		}

		fmt.Fprintf(w, "%s data=[%.2f", data.symbol, data.data[0])
		for _, price := range data.data {
			fmt.Fprintf(w, ", %.2f", price)
		}
		fmt.Fprintf(w, "], average=%.2f", data.average)
	})

	log.Fatal(http.ListenAndServe(args.BindAddr, nil))
}
