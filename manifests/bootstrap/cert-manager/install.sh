#!/usr/bin/env bash
set -eo pipefail

echo "Adding Jetstack Helm repo..."
helm repo add jetstack https://charts.jetstack.io --force-update

echo "Creating cert-manager namespace..."
kubectl create namespace cert-manager --dry-run=client -o yaml | kubectl apply -f -

echo "Installing cert-manager..."
helm upgrade --install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.15.3 \
  --set crds.enabled=true \
  --values manifests/bootstrap/cert-manager/values.yaml

echo "cert-manager installed successfully."
