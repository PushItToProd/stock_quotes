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

Kubernetes deployment with minikube
-----------------------------------

Deployment has been validated with minikube v1.23.2 and kubectl v1.22 on Pop_OS!
20.04 LTS.

### Pre-requisites

- minikube
- kubectl

Install instructions for both tools can be found [here](https://kubernetes.io/docs/tasks/tools/).

### Deployment

The script `kube_e2e.bash` uses Minikube to spin up a new cluster, install the
Ingress addon, configure the secret, deploy the necessary resources for the app,
and validate the app starts successfully.

```
bash kube_e2e.bash
```

This script will prompt for your Alpha Vantage API key while running. You can
also provide it via the environment variable `APIKEY`.

```
APIKEY=demo bash kube_e2e.bash
```

By default, the script tears down the cluster as soon as the deployment is
validated. To prevent this, you can set the environment variable
`SKIP_CLEANUP=1`.

```
SKIP_CLEANUP=1 bash kube_e2e.bash
```

The cluster can be torn down by running the script again without `SKIP_CLEANUP`
set or manually with `minikube delete --profile=stock-quotes-e2e-test`.

### Manual deployment steps

- Run `minikube start --addons ingress` if you don't already have a minikube cluster.
- Run `minikube addons enable ingress` if you have a default minikube cluster and want to use it.
- Create your API key secret: `kubectl create secret generic stock-quotes-secret --from-literal='APIKEY=YOUR API KEY'`
- Deploy the resources using `kubectl apply -f deployment.yml`
  - If you just spun up minikube, you might get an error like `Internal error occurred: failed calling webhook "validate.nginx.ingress.kubernetes.io": Post "https://ingress-nginx-controller-admission.ingress-nginx.svc:443/networking/v1/ingresses?timeout=10s": context deadline exceeded`. In this case, wait and re-run the apply.
- Run `curl -H 'Host: stock-quotes.get' http://$(minikube ip)/` to check the service connects. You might have to retry a few times to validate.
- Remove the created resources using `kubectl delete -f deployment.yml`