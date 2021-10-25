package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type EnvArgs struct {
	bindAddr string
	symbol   string
	ndays    int
	apikey   string
}

func getEnvArgs() *EnvArgs {
	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = ":8080"
	}

	symbol := os.Getenv("SYMBOL")
	if symbol == "" {
		panic("The SYMBOL environment variable must be set")
	}

	ndays_str := os.Getenv("NDAYS")
	if ndays_str == "" {
		panic("The NDAYS environment variable must be set")
	}
	ndays, err := strconv.Atoi(ndays_str)
	if err != nil {
		fmt.Errorf("Invalid value for NDAYS: %v", ndays)
		panic("The NDAYS environment variable must be a valid integer")
	}

	apikey := os.Getenv("APIKEY")
	if apikey == "" {
		panic("The APIKEY environment variable must be set")
	}

	return &EnvArgs{bindAddr: bindAddr, symbol: symbol, ndays: ndays, apikey: apikey}
}

type ApiResponse struct {
	symbol  string
	data    []float64 // in real code, I would not use float64 for currency
	average float64
}

func createApiResponse(symbol string, ndays int) *ApiResponse {
	data := []float64{110.56, 111.25, 115.78}
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
	args := getEnvArgs()

	log.Printf("Starting web server on %s", args.bindAddr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := createApiResponse(args.symbol, args.ndays)
		fmt.Fprintf(w, "%s, data=%v, average=%.2f", data.symbol, data.data, data.average)
	})

	log.Fatal(http.ListenAndServe(args.bindAddr, nil))
}
