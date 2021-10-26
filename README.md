stock_quotes
============

A super simple containerized web service for displaying the last N days of
closing prices for a given stock, written in Go. Prices are retrieved from the
[Alpha Vantage API](https://www.alphavantage.co).

Using the published image with Docker
-------------------------------------

Registry path: `ghcr.io/pushittoprod/stock_quotes`

```sh
docker run \
  --name stock_quotes \
  --env SYMBOL=MSFT \
  --env NDAYS=7 \
  --env APIKEY=[YOUR ALPHAVANTAGE API KEY] \
  --detach \
  ghcr.io/pushittoprod/stock_quotes
```

### Environment variables

Required:

* `SYMBOL`: The stock symbol to retrieve.
* `NDAYS`: The number of days of results to retrieve.
* `APIKEY`: An Alpha Vantage API key.

Development
-----------

### Pre-requisites

* Go 1.17
* GNU Make
* Docker

For local development, update the `APIKEY` variable in `Makefile` with your API
key.

### Run and test

The Makefile provides a few commands for development.

Run `make` to run the tests and start the app using Go. 

* `make` is equivalent to `make default`, which is equivalent to `make test run`.
* `make test` runs just the tests.
* `make run` runs the app.
* `make clean-test` cleans the Go test cache.


### Docker build and run

* `make docker-build` builds a Docker image locally using Docker Compose.
* `make docker-up` builds a Docker image and starts a container using the image in
  Docker.
* `make docker-down` stops the Docker service and cleans up the image builds.
* `make docker-down-all` stops the Docker service and cleans up the 

Kubernetes deployment
---------------------

- install minikube
- `minikube start`
- `minikube addons enable ingress`
- `make kube-apply`
- wait for a bit
- `curl -H 'Host: stock-quotes.get' http://$(minikube ip)/`