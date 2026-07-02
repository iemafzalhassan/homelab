#!/usr/bin/env bash
set -eo pipefail

echo "Adding Traefik Helm repo..."
helm repo add traefik https://helm.traefik.io/traefik --force-update

echo "Creating traefik namespace..."
kubectl create namespace traefik --dry-run=client -o yaml | kubectl apply -f -

echo "Installing Traefik v3..."
helm upgrade --install traefik traefik/traefik \
  --namespace traefik \
  --version 41.0.1 \
  --values manifests/bootstrap/traefik/values.yaml

echo "Traefik installed successfully."
