#!/usr/bin/env bash
# Validate the k8s deployment on a fresh minikube cluster.

if ! which minikube &>/dev/null; then
  echo "error: minikube is not installed or is not on the path" >&2
  exit 1
fi

if [[ ! -f deployment.yml ]]; then
  echo "error: can't find deployment.yml. ensure you're running in the correct directory" >&2
  exit 1
fi

: "${PROFILE:=stock-quotes-e2e-test}"

echo "** Starting minikube cluster"
minikube start --addons ingress --profile="$PROFILE"

# echo "** Enabling ingress"
# minikube --profile="$PROFILE" addons enable ingress

cleanup() {
  # only clean up the cluster if it's not the default
  if [[ "$SKIP_CLEANUP" == 1 ]]; then
    echo "** Skipping cleanup"
    echo "Run 'minikube delete --profile=$PROFILE' to clean up manually"
    return
  fi

  if [[ "$PROFILE" != minikube ]]; then
    echo "** Cleaning up minikube cluster"
    minikube delete --profile="$PROFILE"
  fi
}
trap cleanup INT EXIT

if minikube --profile="$PROFILE" kubectl -- get secret stock-quotes-secret &>/dev/null; then
  echo "** Using existing secret"
else
  echo "** Creating secret"
  if [[ ! "$APIKEY" ]]; then
    read -ersp "Enter Alpha Vantage API key: " APIKEY
  fi
  minikube --profile="$PROFILE" kubectl -- create secret generic stock-quotes-secret --from-literal="APIKEY=$APIKEY"
fi

echo "** Running kubectl apply"

minikube --profile="$PROFILE" kubectl -- apply -f deployment.yml

ip="$(minikube ip --profile="$PROFILE")"
echo "** Found minikube IP: $ip"

echo "** Validating service connects"
success=''
for (( i=0; i<10; i++ )); do
  status="$(curl --silent -m1 --output /dev/null --write-out '%{http_code}' -H 'Host: stock-quotes.get' "http://$ip/")"
  if [[ "$status" == 200 ]]; then
    success=1
    break
  fi
  echo "got unsuccessful status code $status - retrying in 30 seconds..."
  sleep 30
done

if [[ ! "$success" ]]; then
  echo "error: service isn't ready despite waiting repeatedly" >&2
  if [[ "$status" != 000 ]]; then
    echo "got status code: $status"
  else
    echo "failed to connect"
  fi
  exit 1
fi

echo "successfully connected -- e2e test complete"
echo "** Final output: $(curl -H 'Host: stock-quotes.get' "http://$ip/")"
