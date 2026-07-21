# Grafana Cloud Integration & Telemetry Pipeline

Configuring Grafana Alloy to ship OTLP telemetry (Metrics, Loki Logs, Tempo Traces) to Grafana Cloud.

---

## Secret Configuration
Credentials for Grafana Cloud are managed via Azure Key Vault CSI Provider:
- `GRAFANA_CLOUD_OTLP_ENDPOINT`
- `GRAFANA_CLOUD_SERVICE_ACCOUNT_TOKEN`

Alloy receives OTLP on `http://k8s-monitoring-alloy-receiver.monitoring.svc:4318` and forwards to Grafana Cloud OTLP endpoints.
