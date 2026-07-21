# GitOps Incident Break-and-Fix Scenario Guide

Detailed guide on reproducing the stage failure scenario for OpenTelemetry Astronomy Shop.

---

## Scenario: Redis Port Misconfiguration (6379 -> 6380)

1. **Break**:
   ```bash
   kubectl set env deployment/opentelemetry-demo-cartservice -n otel-demo REDIS_ADDR="opentelemetry-demo-redis:6380"
   ```

2. **Query via MCP in IDE**:
   *"Investigate why checkout is failing using Grafana MCP."*

3. **Fix**:
   ```bash
   ./scripts/demo-reset.sh
   ```
