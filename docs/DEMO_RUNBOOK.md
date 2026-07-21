# Live Stage Demo Runbook: Grafana & Friends & CNCF Mumbai

This runbook gives you the exact step-by-step actions and prompts for presenting **AI-Powered Incident Investigation using Grafana MCP** live on stage.

---

## 1. Before Presentation (Pre-Flight Check)

- [ ] Azure 2-Node AKS Cluster is healthy.
- [ ] `opentelemetry-demo` ArgoCD App is `Synced` & `Healthy`.
- [ ] `https://shop.iemafzalhassan.tech` opens Astronomy Shop UI.
- [ ] Grafana Cloud dashboard shows green HTTP 200 responses.
- [ ] Grafana MCP configured in Antigravity IDE (`uvx mcp-grafana`).

---

## 2. On Stage Action Plan

### Step 1: The Healthy Reveal
- Show Astronomy Shop UI (`https://shop.iemafzalhassan.tech`).
- Show Grafana Cloud dashboard: *"Here is our live telemetry pipeline running OpenTelemetry Demo + Grafana Alloy."*

### Step 2: Introduce the Fault (The Break)
- Apply configuration drift (misconfigure Redis port to 6380):
  ```bash
  # Apply broken configuration
  kubectl set env deployment/opentelemetry-demo-cartservice -n otel-demo REDIS_ADDR="opentelemetry-demo-redis:6380"
  ```
- Switch to Grafana Cloud Dashboard: Observe 500 error spike & red indicators.

### Step 3: AI Investigation via Grafana MCP
- Open **Antigravity IDE**.
- Prompt the AI:
  > *"Users are reporting checkout failures. Use Grafana MCP to inspect the latest Loki logs and Prometheus metrics for the cartservice pod in otel-demo. Identify the root cause and recommend the fix."*
- Point out to audience: LLM uses Grafana MCP natively to fetch Loki logs (`ConnectionRefusedError: redis-cart:6380`).

### Step 4: Resolution & Recovery
- Allow AI to patch the Redis port back to `6379`.
- Apply patch or push to git.
- Switch back to Grafana Dashboard and watch panels turn **Green**.

---

## 3. Environment Reset (Post-Demo)

Run the reset script between practice runs or after the stage presentation:
```bash
./scripts/demo-reset.sh
```
