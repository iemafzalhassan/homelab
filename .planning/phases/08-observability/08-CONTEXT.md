# Phase 8: Observability - Context

**Gathered:** 2026-07-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Deploying a lightweight observability agent (Grafana Alloy) to gather cluster logs, traces, and metrics, and forwarding them to Grafana Cloud. This fulfills the observability requirement without crashing the 8GB system node that a local LGTM (Loki, Grafana, Tempo, Mimir) stack would cause.
</domain>

<decisions>
## Implementation Decisions

### Observability Backend
- **D-01:** We will use **Grafana Cloud (Free Tier)** instead of a local `kube-prometheus-stack` or LGTM stack. This is a critical pivot to respect the strict 8GB memory budget of the `Standard_D2as_v5` node.
- **D-02:** We will deploy **Grafana Alloy** (or OpenTelemetry Collector) as the unified agent inside the AKS cluster to scrape metrics, logs, and traces and send them to Grafana Cloud.

### Alerting
- **D-03:** Alertmanager/Notifications will be handled via Grafana Cloud's alerting system, configured to route alerts to a **Slack Webhook**.

### Dashboards
- **D-04:** For this phase, we will focus strictly on **Core Platform** dashboards (Nodes, Pods, Traefik, Jenkins, ArgoCD).

### the agent's Discretion
- Exact configuration of Grafana Alloy Helm chart values to authenticate with Grafana Cloud.
- Which specific metrics to drop/keep to stay within the 10k active series limit of Grafana Cloud's free tier.
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Constraints
- `PROJECT.md` — Specifically the "Budget reality" and "Node sizing rationale" sections that explicitly prohibit heavy local observability stacks like LGTM.
</canonical_refs>

<deferred>
## Deferred Ideas

- **Advanced Traffic & App Analytics:** The user wants highly granular traffic analysis (Client IPs, request paths, response times, page navigation paths, app uptimes, specific user errors). This requires parsing Traefik access logs and instrumenting applications with OpenTelemetry. We will defer the dashboard creation for this to a later application-specific phase once the base Grafana Cloud connection is solid.
- **Local LGTM Stack:** A full local deployment of Loki, Mimir, and Tempo is deferred indefinitely unless the budget constraint is lifted and the node pool is upgraded to at least 16GB RAM.
</deferred>

---
*Phase: 08-observability*
*Context gathered: 2026-07-03 via gsd-discuss-phase*
