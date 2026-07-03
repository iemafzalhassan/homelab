# Phase 8: Observability - Research

## Standard Stack
- **Helm Chart:** `grafana/k8s-monitoring` (replaces `kube-prometheus-stack` to save memory)
- **Core Agent:** Grafana Alloy (deployed automatically by the chart)
- **Dependencies:** `kube-state-metrics`, `prometheus-node-exporter` (included in the chart)
- **Backend:** Grafana Cloud (Metrics, Logs, Traces)

## Architecture Patterns
- **Unified Collector Pattern:** Instead of running separate Prometheus (for metrics), Promtail (for logs), and Tempo (for traces), Grafana Alloy acts as a single, lightweight collector for all telemetry data.
- **Push vs Pull:** Alloy scrapes local Kubernetes components (pull) and forwards them via `remote_write` (push) to Grafana Cloud.
- **GitOps Deployment:** The `k8s-monitoring` Helm chart will be managed via ArgoCD as a standard platform application in `manifests/bootstrap` or `manifests/apps`.
- **Zero-Trust Secrets:** Connection details (Grafana Cloud Token) must be stored in Azure Key Vault and surfaced via Secrets Store CSI as a Kubernetes Secret, rather than hardcoded in `values.yaml`.

## Don't Hand-Roll
- Do not hand-roll custom `.river` (Alloy config) files for basic Kubernetes metrics. The `k8s-monitoring` chart provides highly optimized defaults out-of-the-box for Grafana Cloud.
- Do not deploy a local Prometheus instance. The `k8s-monitoring` chart entirely replaces `kube-prometheus-stack`.

## Common Pitfalls
- **Active Series Limits:** Grafana Cloud Free Tier limits metrics to 10k active series. Enable `metrics.cost.enabled = true` in the chart values to drop high-cardinality metrics (like raw cAdvisor container stats) that easily blow past 10k series.
- **Secret Injection:** The `k8s-monitoring` Helm chart expects a standard Kubernetes Secret for `externalServices.prometheus.basicAuth.password`. Because of our AKV CSI architecture, we must sync the secret from Key Vault to a Kubernetes Secret using `secretObjects` in the `SecretProviderClass` before the chart can mount it.
- **Node Memory:** Ensure the Alloy pods have appropriate resource limits set (e.g., ~250-500MB) to protect the remaining 1.36GB free space on the `Standard_D2as_v5` system node.

## Code Examples
### k8s-monitoring values.yaml baseline
```yaml
cluster:
  name: homelab-aks
externalServices:
  prometheus:
    host: "https://prometheus-prod-XX.grafana.net/api/prom/push"
    basicAuth:
      username: "123456"
      password: "<reference-to-synced-secret>"
  loki:
    host: "https://logs-prod-XX.grafana.net/loki/api/v1/push"
    basicAuth:
      username: "654321"
      password: "<reference-to-synced-secret>"
metrics:
  enabled: true
  cost:
    enabled: true # Keep series count low for Free Tier
logs:
  enabled: true
  pod_logs:
    enabled: true
```
