set -eux

kind delete cluster
kind create cluster

kubectl config set-context kind-kind

docker build -t local-generator-test:latest .

kind load docker-image local-generator-test:latest

helm install generator ./chart --values - <<EOF
env:
  CATALOG_URL: $CATALOG_URL
  CATALOG_REVISION: $CATALOG_REVISION
  GH_USER: $GH_USER

secretEnv:
  type: secret
  values:
    PLUGIN_TOKEN: token
    GH_TOKEN: $GH_TOKEN

image:
  repository: local-generator-test
  tag: latest
EOF
