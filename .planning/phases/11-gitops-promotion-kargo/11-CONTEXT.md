# Phase 11: GitOps Promotion (Kargo) — Context

**Gathered:** 2026-07-09
**Status:** Ready for planning

<domain>
## Phase Boundary

Deploy and configure Kargo as a promotion engine sitting between Jenkins CI and ArgoCD CD. Kargo introduces a manual, human-approved gate between the UAT and PROD environments. This phase does NOT change the CI pipeline (Jenkins) or the CD engine (ArgoCD) — it adds a promotion control plane on top of them.

The deliverable is a live, end-to-end promotion pipeline demonstrated with a minimal demo application:
`Jenkins build → ACR push + Git tag commit → Kargo Warehouse detects → auto-deploys UAT → DevOps clicks Promote → PROD deploys`
</domain>

<decisions>
## Implementation Decisions

### D-01: Promotion Gate Strategy — MANUAL ONLY, No Auto-Promote
- PROD promotion requires an explicit human click in the Kargo UI. There is no automatic promotion to PROD under any condition.
- UAT auto-deploys when a new `Freight` is created (Kargo detects new image tag commit from Jenkins in the GitOps repo).
- Dev team tests UAT manually. When satisfied, they notify the DevOps team.
- DevOps team reviews and clicks "Promote" in the Kargo UI for the PROD Stage.
- Kargo `PromotionPolicy` for PROD Stage: `Manual` (never `Automatic`).

### D-02: Freight Source — Git Commit Watching (image tag in GitOps repo)
- Jenkins CI flow: Developer pushes code → Jenkins builds → pushes image to ACR → Jenkins commits updated image tag to the GitOps repo → ArgoCD currently auto-syncs.
- With Kargo: the `Warehouse` watches the GitOps repo for commits that change the image tag in the demo app's Helm values / kustomize overlay.
- Kargo intercepts before ArgoCD syncs PROD — ArgoCD's PROD Application will be set to `syncPolicy: {}` (manual sync), controlled by Kargo.
- ArgoCD's UAT Application remains auto-sync; Kargo triggers it by updating the UAT Stage's Git path.

### D-03: Demo Application — Minimal "hello-world" Container
- Create a new minimal demo app (`homelab-demo`) — a simple nginx or custom hello-world container.
- The app will have:
  - A `Dockerfile` in `apps/homelab-demo/` in the homelab repo
  - A Jenkins pipeline (`Jenkinsfile.demo`) that builds, tags, and pushes to ACR, then commits the image tag update to `manifests/apps/homelab-demo/uat/values.yaml`
  - Two ArgoCD Applications: `homelab-demo-uat` (namespace: `demo-uat`) and `homelab-demo-prod` (namespace: `demo-prod`)
  - Two Kargo Stages: `uat` and `prod` under a Kargo `Project: homelab-demo`
- This gives a clean, standalone demo story for the Grafana talk without touching platform tools.

### D-04: Kargo RBAC — Keycloak Group Mapping
- `homelab-admins` → Kargo `admin` role (full control, approve PROD, manage Projects)
- `devops` → Kargo `admin` role on the `homelab-demo` Project (approve PROD promotions, view all Stages)
- `developers` → Kargo `viewer` role (can see UAT Stage status + Freight; cannot trigger or approve promotions)
- `viewers` → Kargo `viewer` role (read-only everywhere, same as developers in this context)

Kargo OIDC will use Keycloak Platform realm. Groups claim mapper already configured on all Platform realm clients.

### D-05: Infrastructure Placement
- Kargo runs on the **spot node pool** (toleration: `kubernetes.azure.com/scalesetpriority=spot:NoSchedule`)
- Resource limits: api `100m CPU / 128Mi RAM`, controller `200m CPU / 256Mi RAM` — well within headroom
- Exposed at `https://kargo.smapatticare.com` via Traefik HTTPRoute + cert-manager TLS (wildcard cert already covers this)
- Installed via ArgoCD Application (Helm chart: `kargo` from `ghcr.io/akuity/kargo-charts`)
</decisions>

<canonical_refs>
## Files Downstream Agents Must Read

- `.planning/ROADMAP.md` — Phase 11 goals and success criteria
- `.planning/REQUIREMENTS.md` — PROMO-01 through PROMO-05 (to be added)
- `manifests/bootstrap/argocd/values.yaml` — ArgoCD OIDC + RBAC config (Kargo integrates with ArgoCD)
- `manifests/bootstrap/jenkins/values.yaml` — Jenkins JCasC + OIDC config (Jenkins is the CI trigger)
- `manifests/apps/` — App-of-Apps pattern, where new Kargo ArgoCD Application will live
- `docs/ADR-001-platform-auth-rbac.md` — Keycloak group structure, OIDC clients, RBAC decisions
</canonical_refs>

<code_context>
## Reusable Assets & Patterns

### Existing ArgoCD Application pattern
- All platform apps follow the App-of-Apps pattern in `manifests/apps/`
- Kargo's ArgoCD Application follows the same structure

### Existing Traefik HTTPRoute pattern
- Every platform service has an HTTPRoute in its namespace + ReferenceGrant
- `kargo.smapatticare.com` follows the same pattern as `argocd.smapatticare.com`

### Existing OIDC client pattern
- All platform tools (ArgoCD, Jenkins) use Keycloak Platform realm with `groups` claim mapper
- Kargo OIDC client follows the exact same Keycloak client setup

### Spot node toleration (all user workloads)
```yaml
tolerations:
  - key: "kubernetes.azure.com/scalesetpriority"
    operator: "Equal"
    value: "spot"
    effect: "NoSchedule"
nodeSelector:
  kubernetes.azure.com/scalesetpriority: "spot"
```

### ArgoCD Application pointing to Helm chart
```yaml
source:
  repoURL: https://charts.kargo.io
  chart: kargo
  targetRevision: "1.x.x"  # latest stable
  helm:
    valuesObject: { ... }
```
</code_context>

<deferred>
## Deferred Ideas (out of scope for this phase)

- Progressive delivery / canary promotions (Argo Rollouts integration) — separate phase if needed
- Automated test gates between UAT and PROD (e.g. run a smoke test before Kargo allows promotion) — can be added as a Kargo `VerificationTemplate` in a future phase
- Promote platform tools (ArgoCD/Jenkins/Keycloak) between UAT/PROD realms — separate concern, not a Kargo use case
</deferred>
