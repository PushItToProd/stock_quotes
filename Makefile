# I like capturing common commands with Make or Invoke (pyinvoke.org) to
# streamline development.

### Configuration variables ###

# Export environment variables for use with `go run` and `go test`.
export SYMBOL := MSFT
export NDAYS := 7
export APIKEY := demo


### Commands for running and testing with Go ###

default: test run

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test -v ./alphavantage

# Force tests to run instead of using cached results.
.PHONY: clean-test
clean-test:
	go clean -testcache


### Commands for building and running with Docker ###

# Using docker-compose here is probably redundant with K8s but I like how
# Compose encapsulates the build/run/clean steps.

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


### Commands for K8s ###

.PHONY: kube-apply
kube-apply:
	kubectl apply -f deployment.yml

.PHONY: kube-delete
kube-delete:
	kubectl delete -f deployment.yml
