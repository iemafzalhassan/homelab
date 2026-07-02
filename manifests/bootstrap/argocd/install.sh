#!/usr/bin/env bash
set -eo pipefail

echo "Adding ArgoCD Helm repo..."
helm repo add argo https://argoproj.github.io/argo-helm --force-update

echo "Creating argocd namespace..."
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -

echo "Installing ArgoCD v3.4.4 (chart 10.1.0)..."
helm upgrade --install argocd argo/argo-cd \
  --namespace argocd \
  --version 10.1.0 \
  --values manifests/bootstrap/argocd/values.yaml

echo "Applying HTTPRoute for ArgoCD..."
kubectl apply -f manifests/bootstrap/argocd/argocd-httproute.yaml

echo "ArgoCD bootstrap complete."
