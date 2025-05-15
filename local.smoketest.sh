#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o errtrace
set -u
shopt -s inherit_errexit 2> /dev/null || true

IFS=$'\n\t'

_info() {
  echo >&2 "[info] $*"
}


_wait_for() {
  local status="$1"
  _info "Waiting for generator to be ${status}..."
  while ! kubectl get pods -l app.kubernetes.io/name=joy-generator -o jsonpath="{.items[0].status.containerStatuses[0].${status}}" | grep true; do
    sleep 1
  done

}

_main() {
  kind delete cluster
  kind create cluster

  kubectl config set-context kind-kind

  docker build -t local-generator-test:latest .

  kind load docker-image local-generator-test:latest

  kubectl create secret generic gcp-credentials --from-file=credentials="${GCP_SERVICE_ACCOUNT_CREDENTIALS_FILE}"

  helm install generator ./chart --values - << EOF
env:
  CATALOG_URL: $CATALOG_URL
  CATALOG_REVISION: ${CATALOG_REVISION:-master}
  GH_USER: $GH_USER
  GOOGLE_ARTIFACT_REPOSITORY: ${GOOGLE_ARTIFACT_REPOSITORY:-northamerica-northeast1-docker.pkg.dev}

credentialsSecret:
  name: gcp-credentials
  key: credentials

secretEnv:
  type: secret
  values:
    PLUGIN_TOKEN: token
    GH_TOKEN: $(gh auth token)

image:
  repository: local-generator-test
  tag: latest
EOF

  _wait_for 'started'
  kubectl logs -l app.kubernetes.io/name=joy-generator --follow &
  trap 'kill $(jobs -p)' EXIT

  _wait_for 'ready'
  local port="${PORT:-8080}"
  _info "Service is available. Port forwarding to localhost:$port} ..."
  _info "To test, run: curl -X POST 127.0.0.1:8080/api/v1/getparams.execute -H 'Authorization: Bearer token' -d {}"
  kubectl port-forward svc/generator-joy-generator "${port}:80"

}

_main "$@"
