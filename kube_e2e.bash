#!/usr/bin/env bash
# Validate the k8s deployment on a fresh minikube cluster to ensure this is
# repeatable.

if ! which minikube &>/dev/null; then
  echo "error: minikube is not installed or is not on the path" >&2
  exit 1
fi

if [[ ! -f deployment.yml ]]; then
  echo "error: can't find deployment.yml. ensure you're running in the correct directory" >&2
  exit 1
fi

if [[ ! "$PROFILE" ]]; then
  PROFILE=stock-quotes-e2e-test
  echo "** Using default profile $PROFILE"
else
  echo "** Using provided PROFILE=$PROFILE -- this will not be cleaned up"
  SKIP_CLEANUP=1
fi

echo "** Starting minikube cluster with profile $PROFILE"
minikube start --memory=2G --cpus=1 --addons ingress --profile="$PROFILE"

# Automatically clean up when we're done.
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
trap cleanup EXIT

if minikube --profile="$PROFILE" kubectl -- get secret stock-quotes-secret &>/dev/null; then
  echo "** Using existing secret"
else
  echo "** Creating secret"
  if [[ ! "$APIKEY" ]]; then
    read -ersp "Enter Alpha Vantage API key: " APIKEY
  fi
  minikube --profile="$PROFILE" kubectl -- create secret generic stock-quotes-secret --from-literal="APIKEY=$APIKEY"
fi

# It takes some time for the Ingress controller to become usable, so we need to
# wait a bit and potentially retry a few times before it'll work. Otherwise
# we'll get an internal error about a failure to call the webhook
# validate.nginx.ingress.kuberentes.io.
sleep 10
echo "** Running kubectl apply"
apply_success=''
for (( i=0; i<10; i++ )); do
  if minikube --profile="$PROFILE" kubectl -- apply -f deployment.yml; then
    apply_success=1
    break
  fi
  echo "waiting to retry..."
  sleep 30
done

if [[ ! "$apply_success" ]]; then
  echo "error: kubectl apply failed after retries -- giving up" >&2
  exit 1
fi

ip="$(minikube ip --profile="$PROFILE")"
echo "** Found minikube IP: $ip"

# It takes some time for the Ingress to get configured, so we potentially have
# to retry a few times.
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
echo "** Final output: $(curl --silent --output /dev/stdout -H 'Host: stock-quotes.get' "http://$ip/")"
