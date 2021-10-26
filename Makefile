export SYMBOL := MSFT
export NDAYS := 7
# replace with actual API key
export APIKEY := demo

default: test run

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test -v ./alphavantage