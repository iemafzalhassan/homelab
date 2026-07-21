#!/usr/bin/env bash
# Verification script for OpenTelemetry Demo & Observability stack
set -euo pipefail

echo "=================================================="
echo "      OpenTelemetry Demo Verification Checklist    "
echo "=================================================="

echo "[1/4] Checking ArgoCD Application Status..."
kubectl get application opentelemetry-demo -n argocd -o jsonpath='{.status.sync.status}{" "}{.status.health.status}{"\n"}'

echo "[2/4] Checking otel-demo Pods..."
kubectl get po -n otel-demo

echo "[3/4] Checking HTTPRoute for shop.iemafzalhassan.tech..."
kubectl get httproute opentelemetry-demo-frontend -n otel-demo

echo "[4/4] Checking Alloy OTLP Receiver Endpoint..."
kubectl get svc -n monitoring | grep alloy-receiver || true

echo "=================================================="
echo "Verification complete!"
