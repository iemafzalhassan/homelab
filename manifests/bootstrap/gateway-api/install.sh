#!/usr/bin/env bash
set -eo pipefail

echo "Installing Kubernetes Gateway API CRDs (v1.5.1 Standard Channel)..."
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.5.1/standard-install.yaml
echo "Gateway API CRDs installed."
