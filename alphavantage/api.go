package alphavantage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var apikey = os.Getenv("APIKEY")

func init() {
	if apikey == "" {
		panic("The APIKEY environment variable must be set")
	}
}

type TimeSeriesEntry struct {
	Open             string `json:"1. open"`
	High             string `json:"2. high"`
	Low              string `json:"3. low"`
	Close            string `json:"4. close"`
	AdjustedClose    string `json:"5. adjusted close"`
	Volume           string `json:"6. volume"`
	DividendAmount   string `json:"7. dividend amount"`
	SplitCoefficient string `json:"8. split coefficient"`
}

type TimeSeriesDailyAdjusted struct {
	MetaData   map[string]string          `json:"Meta Data"`
	TimeSeries map[string]TimeSeriesEntry `json:"Time Series (Daily)"`
}

func getTimeSeriesDailyAdjusted(symbol string) TimeSeriesDailyAdjusted {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s",
		apikey, symbol,
	)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("failed reading response body")
	}

	var data TimeSeriesDailyAdjusted
	if err = json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	return data
}

func getClosingDataFromResponse(data TimeSeriesDailyAdjusted, ndays int) []float64 {
	timeseries := data.TimeSeries

	// get the last ndays of dates
	dateKeys := make([]string, len(data.TimeSeries))
	i := 0
	for k := range timeseries {
		dateKeys[i] = k
		i++
	}
	sort.Strings(dateKeys)

	dateKeys = dateKeys[len(dateKeys)-ndays:]
	closingData := make([]float64, ndays)
	for i, date := range dateKeys {
		close, err := strconv.ParseFloat(timeseries[date].Close, 64)
		if err != nil {
			panic("Failed converting close value")
		}
		closingData[i] = close
	}
	return closingData
}

func GetClosingData(symbol string, ndays int) []float64 {
	resp := getTimeSeriesDailyAdjusted(symbol)
	return getClosingDataFromResponse(resp, ndays)
}
