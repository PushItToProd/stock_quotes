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

.PHONY: docker-build
docker-build:
	docker-compose build

.PHONY: docker-up
docker-up:
	docker-compose up --build -d

.PHONY: docker-down
docker-down:
	docker-compose down --rmi local --volumes --remove-orphans

.PHONY: docker-down-all
docker-down-all:
	docker-compose down --rmi all --volumes --remove-orphans