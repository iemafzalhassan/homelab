# Discussion Log: Phase 08-observability

**Date:** 2026-07-03

## Q1: Prometheus Storage & Full Observability Stack
**Context:** User wanted to know whether to use Ephemeral or Persistent storage for Prometheus, but also requested full end-to-end LGTM (Loki, Mimir, Tempo) + OTEL.
**User Selection:** User requested full LGTM but asked for analysis on node size.
**Agent Analysis:** Analyzed `kubectl top nodes`. 8GB node has 5.7GB allocatable. 4.3GB is used. Only 1.36GB remains. Full LGTM requires ~2GB+. Alerted user that this would cause OOM crashes.
**Resolution:** User agreed to pivot to Option A: **Grafana Cloud Free Tier** with a lightweight agent to save node memory.

## Q2: Azure Monitor Integration
**Context:** User provided Grafana Cloud's Azure Monitor integration UI fields.
**Resolution:** Agent generated an Azure Service Principal `grafana-cloud-azure-monitor` with `Monitoring Reader` access and provided Client ID, Secret, and Tenant ID.

## Q3: Alert Notifications
**Context:** Where should critical cluster alerts go?
**Options:** Discord, Slack, Email, UI Only.
**User Selection:** Slack Webhook.

## Q4: Dashboards
**Context:** Core platform dashboards vs custom app dashboards.
**User Selection:** Core platform only for now. "later we need to create dashbioard that monitors each traffic along with user ip address... what page navigating... uptime of all the apppplications"
**Resolution:** Focus on core platform now. Logged the advanced traffic and app analytics requirement in the Deferred section.
