# Track A: Local Kind / k3d Community Workshop Guide

This guide allows meetup participants to reproduce the OpenTelemetry Demo + Grafana MCP setup locally using zero-cost tooling (`kind` / `k3d`).

---

## Prerequisites
- Docker / Podman
- `kind` or `k3d`
- `kubectl` & `helm`
- Antigravity IDE & `uv` package manager (`curl -LsSf https://astral.sh/uv/install.sh | sh`)

---

## Step 1: Create Local Kubernetes Cluster
```bash
kind create cluster --name otel-workshop
```

## Step 2: Deploy OpenTelemetry Demo
```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update

helm install opentelemetry-demo open-telemetry/opentelemetry-demo \
  --namespace otel-demo --create-namespace
```

## Step 3: Configure Grafana MCP Server in IDE
Add the following to your IDE MCP configuration (`mcp_config.json`):
```json
{
  "mcpServers": {
    "grafana": {
      "command": "uvx",
      "args": ["mcp-grafana"],
      "env": {
        "GRAFANA_URL": "http://localhost:3000",
        "GRAFANA_SERVICE_ACCOUNT_TOKEN": "<your-token>"
      }
    }
  }
}
```
