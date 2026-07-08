# Phase 11: GitOps Promotion (Kargo) — Discussion Log

**Date:** 2026-07-09
**Phase:** 11 — GitOps Promotion (Kargo)
**Outcome:** Context captured, ready for planning

---

## Discussion Summary

### Area 1: Promotion Gate Strategy

| | |
|---|---|
| **Question** | Should Kargo auto-promote to PROD if tests pass, or require a human click always? |
| **Options** | Auto-promote if tests pass / Manual only |
| **Decision** | **Manual only — no auto-promote to PROD under any condition** |
| **Rationale** | Dev team must manually test UAT. Once satisfied they notify DevOps. DevOps team reviews and approves PROD promotion via Kargo UI. |

---

### Area 2: Freight Source

| | |
|---|---|
| **Question** | What does Kargo watch to know a new version is ready? |
| **Context explained** | Dev pushes code → Jenkins CI builds → pushes to ACR → Jenkins commits updated image tag to GitOps repo → ArgoCD currently auto-syncs everywhere |
| **Decision** | **Kargo Warehouse watches the GitOps repo for image tag commits** (Git-based freight source) |
| **Rationale** | The existing Jenkins CI already commits image tags to Git. Kargo intercepts the PROD ArgoCD Application's sync (set to manual), UAT remains auto. |

---

### Area 3: Demo Application

| | |
|---|---|
| **Question** | Which application gets the full UAT→PROD Kargo pipeline first? |
| **Options** | Existing sample app / New minimal demo app / Platform tools |
| **Decision** | **Create a new minimal demo app (`homelab-demo`)** — nginx or hello-world container |
| **Rationale** | Clean, standalone demo story for the Grafana talk. Doesn't touch platform tools. Shows the full pipeline clearly. |

---

### Area 4: Kargo RBAC

| | |
|---|---|
| **Question** | Who can approve PROD promotions? Who can only view? |
| **Options** | (A) homelab-admins+devops approve PROD, developers view UAT only, viewers read-only / (B) homelab-admins only approve / (C) homelab-admins+devops approve, developers can trigger UAT |
| **Decision** | **Option A**: `homelab-admins` + `devops` approve PROD; `developers` view UAT status only; `viewers` read-only |
| **Rationale** | Mirrors the production pattern: devops owns release gates, developers own code but not production keys. |

---

## Deferred Ideas

| Idea | Disposition |
|------|-------------|
| Automated test gates before PROD promotion (Kargo `VerificationTemplate`) | Future phase — adds complexity, not needed for demo |
| Progressive delivery / canary (Argo Rollouts) | Separate concern, separate phase |
| Promote platform tools (ArgoCD/Jenkins) between realms | Not a Kargo use case |
