#!/usr/bin/env bash
# Validate the k8s deployment on a fresh minikube cluster to ensure this is
# repeatable.

# Only display output in color if stdout is a terminal and provide an option
# to switch it off.
if [[ -t 1 ]] && [[ ! "$KUBE_E2E_NOCOLOR" ]]; then
  C_RESET="$(tput sgr0)"
  C_ERR="$(tput setaf 9)$(tput bold)"
  C_INFO="$(tput setaf 14)$(tput bold)"
else
  C_RESET=''
  C_ERR=''
  C_INFO=''
fi

fatal() {
  echo "${C_ERR}error: $*${C_RESET}" >&2
  exit 1
}

info() {
  echo "${C_INFO}** $*${C_RESET}"
}

main() {
  if ! which minikube &>/dev/null; then
    fatal "minikube is not installed or is not on the path"
  fi

  if ! which curl &>/dev/null; then
    fatal "curl is not installed or is not on the path"
  fi

  if [[ ! -f deployment.yml ]]; then
    fatal "error: can't find deployment.yml. are you running in the correct" \
          "directory?"
  fi

  if [[ ! "$PROFILE" ]]; then
    PROFILE=stock-quotes-e2e-test
    info "Using default profile $PROFILE"
  else
    info "Using provided PROFILE=$PROFILE -- this will not be cleaned up"
    SKIP_CLEANUP=1
  fi

  info "Starting minikube cluster with profile $PROFILE"
  minikube start \
    --memory=2G \
    --cpus=2 \
    --addons ingress \
    --profile="$PROFILE" \
    --keep-context \
  || fatal "Failed to start minikube cluster"

  # Automatically clean up when we're done.
  cleanup() {
    # only clean up the cluster if it's not the default
    if [[ "$SKIP_CLEANUP" == 1 ]]; then
      info "Skipping cleanup"
      echo "Run 'minikube delete --profile=$PROFILE' to clean up manually"
      return
    fi

    if [[ "$PROFILE" != minikube ]]; then
      info "Cleaning up minikube cluster"
      minikube delete --profile="$PROFILE"
    fi
  }
  trap cleanup EXIT

  if minikube --profile="$PROFILE" kubectl -- \
      get secret stock-quotes-secret &>/dev/null; then
    info "Using existing secret"
  else
    info "Creating secret"
    if [[ ! "$APIKEY" ]]; then
      read -ersp "Enter Alpha Vantage API key: " APIKEY
    fi
    minikube --profile="$PROFILE" kubectl -- \
      create secret generic stock-quotes-secret --from-literal="APIKEY=$APIKEY"
  fi

  # It takes some time for the Ingress controller to become usable, so we need
  # to wait a bit and potentially retry a few times before it'll work. Otherwise
  # we'll get an internal error about a failure to call the webhook
  # validate.nginx.ingress.kuberentes.io.
  sleep 10
  info "Running kubectl apply"
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
    fatal "kubectl apply failed after retries -- giving up"
  fi

  ip="$(minikube ip --profile="$PROFILE")"
  info "Found minikube IP: $ip"

  # It takes some time for the Ingress to get configured, so we potentially have
  # to retry a few times.
  info "Validating service connects"
  success=''
  for (( i=0; i<10; i++ )); do
    status="$(curl --silent -m1 --output /dev/null --write-out '%{http_code}' -H 'Host: stock-quotes.get' "http://$ip/")"
    if [[ "$status" == 200 ]]; then
      success=1
      break
    fi
    echo "got unsuccessful status code $status - retrying in 10 seconds..."
    sleep 10
  done

  if [[ ! "$success" ]]; then
    err_msg="failed to connect"
    if [[ "$status" != 000 ]]; then
      err_msg="got status code: $status"
    fi
    fatal "error: service isn't ready despite waiting repeatedly -- $err_msg"
  fi

  echo "successfully connected -- e2e test complete"
  info "Final output: $(curl --silent --output /dev/stdout -H 'Host: stock-quotes.get' "http://$ip/")"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  main "$@"
fi