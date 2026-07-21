#!/usr/bin/env bash
# Reset OpenTelemetry Demo to healthy state for rehearsal/stage use
set -euo pipefail

echo "==> Resetting OpenTelemetry Demo environment to Healthy state (Redis 6379)..."
kubectl annotate application opentelemetry-demo -n argocd argocd.argoproj.io/refresh=hard --overwrite
kubectl rollout restart deployment -n otel-demo opentelemetry-demo-cartservice || true
echo "==> Environment reset complete. Dashboard returning to Green."
