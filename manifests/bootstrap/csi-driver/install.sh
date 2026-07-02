#!/bin/bash
set -e

cd "$(dirname "$0")"

# Add Helm repos
helm repo add secrets-store-csi-driver https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts
helm repo add csi-secrets-store-provider-azure https://azure.github.io/secrets-store-csi-driver-provider-azure/charts
helm repo update

# Install Secrets Store CSI Driver
# Using version 1.6.0 which is stable for K8s 1.30+
helm upgrade --install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver \
  --namespace kube-system \
  --version 1.6.0 \
  -f values-csi-driver.yaml

# Install Azure Key Vault Provider
helm upgrade --install azure-csi-provider csi-secrets-store-provider-azure/csi-secrets-store-provider-azure \
  --namespace kube-system \
  -f values-akv-provider.yaml

echo "Secrets Store CSI Driver and Azure Provider installed successfully."
