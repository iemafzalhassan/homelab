# Phase 12: Kyros — The Trusted Software Supply Chain Platform - Context

**Gathered:** 2026-07-12 (initial 2026-07-12; revised after cross-AI review + discuss-phase)
**Status:** Ready for planning
**Source code:** `kyros/` (Go monorepo + Next.js 15 + Taskfile + Docker Compose dev stack)
**Pre-review:** `12-REVIEWS.md` (Antigravity CLI, 2026-07-12 — 3 HIGH, 2 MED findings)

<domain>
## Phase Boundary

Phase 12 builds **Kyros** as a cloud-native, OCI-compliant software supply chain platform on the existing AKS homelab cluster. MVP scope: a working OCI registry (push/pull), a publisher dashboard, a Trust Score engine, and a baseline security pipeline (SBOM + scan + keyless signing). Deployed via ArgoCD on the same homelab cluster from phases 1-11. MVP URLs: `kyros.smapatticare.com` (dashboard) + `registry.kyros.smapatticare.com` (OCI).

**Out of scope for Phase 12** (deferred to post-MVP phases): public OCI hosting with self-serve org signup, marketplace/billing, multi-cloud storage, multi-region, Helm/WASM/AI-artifact repos, kyros.io production domain.

</domain>

<decisions>
## Implementation Decisions

### Search Backend (resolution of Antigravity review HIGH-1)

- **D-01:** **Typesense** replaces Elasticsearch. Search runs as a single Go binary, ~50-100 MB RAM typical, official Helm chart, official cloud-native image, used in production at multiple companies. Single `typesense-pod` per cluster is sufficient for MVP index size (registry tags, SBOM search, org search).
  - Connects from Go API via `typesense-go` client
  - Persistence: PVC backed by Azure Blob (no separate SSD)
  - Backup: nightly snapshot to Azure Blob via a CronJob calling `typesense-cli --operation=disk-snapshot`
  - Resource limits: `requests: cpu=100m memory=128Mi`, `limits: cpu=500m memory=512Mi`

### Public/Private Rollout (resolution of Antigravity review HIGH-3)

- **D-02:** **Private-only for MVP.** No public OCI repos, no self-serve org signup. The `kyros.smapatticare.com` UI and `registry.kyros.smapatticare.com` are reachable from the public internet BUT every repo is private by default; only invited org members can push. Public repos and self-serve signup are deferred to a post-MVP phase (registered in Deferred Ideas).
  - Mitigates the $100-budget DoS risk (no abuse vector for hosting malware images that drive Azure Blob egress costs)
  - All repos are created via Keycloak admin invitation, not via public signup form

### Registry Storage & Garbage Collection (resolution of Antigravity review MED-1)

- **D-03:** **CronJob-based GC with read-only lock + per-org quota enforcement.**
  - **GC:** Kubernetes `CronJob` runs daily at 03:00 IST, executes `registry garbage-collect` against the cncf/distribution API while a maintenance `HTTPRoute` annotation flips the registry to read-only (`503` on writes). Step:
    1. Patch Traefik `HTTPRoute` for `registry.kyros.smapatticare.com` to set `503` on POST/PUT
    2. Wait 60s for in-flight uploads to drain
    3. Run `docker exec registry registry garbage-collect --delete-untagged /etc/docker/registry/config.yml`
    4. Patch `HTTPRoute` back to normal
  - **Quota:** Per-org Azure Blob storage quota enforced in the Go API *before* the upload-to-registry step. The API calls Azure Blob `GetContainerUsage` and rejects the push if `org_storage_used_bytes + blob_size > org_quota_bytes`. Default quota: 50 GB per org (tunable via Kyros admin UI).
  - **Blob linking:** cncf/distribution configured with `redirect.disable: false` and `storage.delete.enabled: true` to support hard-delete on repo deletion.

### Security Pipeline Scope Trim (resolution of Antigravity review MED-2 + HIGH-2)

- **D-04:** **Trim the "5 features from day 1" to a 3-feature MVP slice.** Defer Grype (dual-scanner) and SLSA L3+ to post-MVP. Day-1 pipeline is:
  1. **Trivy** vulnerability scan (single scanner)
  2. **Syft** SBOM (SPDX + CycloneDX, both emitted)
  3. **Cosign keyless** signing via Sigstore Fulcio + Rekor (outbound HTTPS to `oauth2.sigstore.dev`, `fulcio.sigstore.dev`, `rekor.sigstore.dev` — must be allowlisted in NSG egress when Phase 9 lands)
  - **SLSA L2 only** for MVP. L3+ requires hermetic builds + isolated builders — post-MVP.
  - **Kaniko** is the build engine (locked). Runs as ephemeral pods on the spot node pool (toleration: `kubernetes.azure.com/scalesetpriority=spot:NoSchedule`).

### Multi-Tenancy Isolation (resolution of Antigravity review MED-3)

- **D-05:** **Postgres Row-Level Security (RLS) for tenant isolation.** Every multi-tenant table (`organizations`, `repositories`, `tags`, `digests`, `images`, `sboms`, `scan_results`, `signatures`) has `tenant_id uuid NOT NULL` column + RLS policy `USING (tenant_id = current_setting('app.current_tenant')::uuid)`. The Go API sets the GUC `app.current_tenant` on every transaction. App-layer `WHERE tenant_id = $1` is also kept as defense-in-depth, but RLS is the source of truth. Cross-tenant data leaks become structurally impossible at the database layer.

### Locked Tech Stack (carried from initial CONTEXT.md)

- Backend: Go 1.23, `go-containerregistry`, `oras-go`, `chi` router, `viper` config, `zap` logger, `prometheus/client_golang` metrics
- Registry Engine: `cncf/distribution` v3 with Azure Blob backend
- Frontend: Next.js 15 + TypeScript + Tailwind + ShadCN
- Auth: Keycloak OIDC (Platform realm at `sso.smapatticare.com`)
- Database: PostgreSQL 17 (RLS for multi-tenancy)
- Cache: Redis 7
- Search: **Typesense** (D-01)
- Queue: NATS JetStream
- Storage: Azure Blob Storage
- CDN/Edge: Cloudflare (proxied, DNS for `*.smapatticare.com` already managed)
- Build: Kaniko (K8s-native, spot node pool)
- Observability: Full LGTM (Prometheus + Grafana + Loki + Tempo + Alloy) — already running from Phase 8

### Trust Score Formula (LOCKED — unchanged)

Composite weighted score 0-100:
- CVE severity: 35% (0 Critical = max contribution)
- SBOM completeness: 20% (both SPDX + CycloneDX present)
- SLSA Level: 20% (L1=25%, L2=50%, L3=75%, L4=100%)
- Cosign signature: 10%
- Scan freshness: 10% (scanned in last 24h = max)
- Base image provenance: 5%
- Badges: Trusted ≥80, Moderate 50-79, At Risk <50

### Auth / RBAC Model (LOCKED — unchanged)

- `org:owner`, `org:admin`, `org:devops`, `org:developer`, `org:viewer`
- Keycloak Platform realm at `sso.smapatticare.com`, OIDC client `kyros-{dashboard,registry,api,builder,scanner,sbom,signer,notifier}`
- Federation: Workload Identity per Kyros service (one UAMI per `cmd/*`)

### MVP Roadmap (12 weeks, refined)

| Sub-phase | Weeks | Deliverable | Status |
|---|---|---|---|
| 12.1 | w1-2 | Foundation: monorepo, Postgres schema + RLS, Keycloak SSO clients, cncf/distribution Helm deploy on AKS, HTTPRoute for `registry.kyros.smapatticare.com` | Planned |
| 12.2 | w3-4 | Registry Core: docker push/pull, manifest/tag API, basic Next.js dashboard, Typesense indexer (D-01), per-org quota enforcement (D-03) | Planned |
| 12.3 | w5-7 | Security Pipeline: NATS JetStream workers, Trivy scanner, Syft SBOM, Trust Score engine + UI, **private-only access enforced** (D-02) | Planned |
| 12.4 | w8-10 | Publisher Dashboard: org/team mgmt UI, image builder (Kaniko on spot), Cosign keyless signing, GC CronJob (D-03) | Planned |
| 12.5 | w11-12 | Polish: LGTM instrumentation, search via Typesense UI, CDN caching, landing page, runbooks | Planned |

### Claude's Discretion

- Choice of Typesense Helm chart values (single-node, PVC, resource limits) — set to defaults that fit the 8 GB system node budget
- Exact CronJob schedule for GC (locked to 03:00 IST, daily)
- Default per-org quota (50 GB)
- All public OCI policy decisions — explicit non-goal for MVP

### Folded Todos

None — no todos matched Phase 12 in `cross_reference_todos`.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Source Code (existing monorepo)
- `kyros/go.mod` — Go module + dependency pins
- `kyros/cmd/api/main.go` — API entrypoint (HTTP server bootstrap)
- `kyros/cmd/registry/main.go` — Registry service entrypoint
- `kyros/cmd/builder/main.go`, `kyros/cmd/scanner/main.go`, `kyros/cmd/sbom/main.go`, `kyros/cmd/signer/main.go`, `kyros/cmd/notifier/main.go` — all six security pipeline workers
- `kyros/internal/registry/auth.go`, `auth_test.go`, `proxy.go` — registry auth + proxy glue
- `kyros/internal/oci/` — OCI manifest/tag/digest types
- `kyros/internal/storage/` — Azure Blob adapter
- `kyros/internal/auth/` — Keycloak OIDC + token validation
- `kyros/internal/trust/` — Trust Score engine
- `kyros/internal/httpserver/` — HTTP server middleware (chi)
- `kyros/internal/config/config.go` — viper config loader
- `kyros/internal/logger/` — zap logger
- `kyros/internal/metrics/` — Prometheus metrics
- `kyros/database/migrations/` — SQL migrations
- `kyros/deploy/helm/kyros/` — Helm chart
- `kyros/infrastructure/terraform/` — Terraform for any Kyros-specific Azure resources (Blob container, etc.)
- `kyros/docker-compose.dev.yml` — local dev stack (Postgres, Redis, NATS, API, Web)
- `kyros/Taskfile.yml` — task runner (dev:api, dev:web, dev:all, build:go, build:web, test:go, test:web, lint:go, lint:web, migrate)
- `kyros/apps/web/` — Next.js 15 dashboard
- `kyros/packages/config/`, `kyros/packages/types/`, `kyros/packages/ui/` — shared TS packages
- `kyros/services/{builder,registry,sbom,scanner,signer,notifier}/` — service definitions

### Pre-existing Review
- `.planning/phases/12-kyros-platform/12-REVIEWS.md` — Antigravity CLI cross-AI review (3 HIGH, 2 MED); all HIGH items resolved by decisions D-01 through D-04 in this document

### Homelab Platform Reusables
- `.planning/PROJECT.md` — overall homelab scope, budget, node sizing
- `.planning/REQUIREMENTS.md` — v1 platform requirements (no Kyros REQ-IDs yet — see "Folded Todos" / open question)
- `.planning/research/STACK.md`, `FEATURES.md`, `ARCHITECTURE.md`, `PITFALLS.md`, `SUMMARY.md` — research for the homelab platform; Kyros reuses most of these patterns
- Established Traefik Gateway API pattern (HTTPRoute + ReferenceGrant) — `kyros/deploy/helm/kyros/templates/httproute.yaml` should follow the pattern from earlier phases
- Established cert-manager pattern with `*.smapatticare.com` wildcard via Cloudflare DNS-01
- Established ArgoCD App-of-Apps pattern (Kyros is a single `Application` CR)

### External Specs
- OCI Distribution Spec v1.1 — `https://github.com/opencontainers/distribution-spec`
- cncf/distribution — `https://github.com/distribution/distribution`
- go-containerregistry — `https://github.com/google/go-containerregistry`
- Sigstore/Cosign — `https://github.com/sigstore/cosign`
- Syft — `https://github.com/anchore/syft`
- Trivy — `https://github.com/aquasecurity/trivy`
- Kaniko — `https://github.com/GoogleContainerTools/kaniko`
- NATS JetStream — `https://docs.nats.io`
- Typesense — `https://typesense.org/docs/`
- Typesense Go client — `https://github.com/typesense/typesense-go`

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets

- **`kyros/internal/auth/`** (Keycloak OIDC client + JWT validator) — reuse for all Kyros HTTP services; the Keycloak Platform realm is already running at `sso.smapatticare.com`
- **`kyros/internal/oci/`** (OCI types) — `manifest`, `tag`, `digest` types already defined; reuse in registry proxy + TypeSense indexer
- **`kyros/internal/storage/`** (Azure Blob adapter) — use for both registry blob storage and Typesense snapshot backup
- **`kyros/internal/httpserver/`** (chi HTTP server) — middleware, request ID, logger, metrics — reuse for all Go services
- **`kyros/internal/metrics/`** (Prometheus) — already emits metrics; reuse for Trust Score freshness, scanner latency, signer queue depth
- **`kyros/internal/trust/`** (Trust Score engine) — pure function over CVE/SBOM/SLSA/signature signals; no I/O, easily testable
- **`kyros/packages/ui/`** (ShadCN components) — `avatar`, `badge`, `button`, `card`, `dropdown-menu`, `input` — already built
- **Homelab platform services** — Traefik + cert-manager + Keycloak + LGTM are all already running; Kyros gets HTTPRoutes, TLS, SSO, observability for free

### Established Patterns

- **Chi router** with middleware chain: `RequestID → Logger → Recoverer → CORS → Auth → Handler` — already in `internal/httpserver/`
- **Viper config** from env + file with struct unmarshal — already in `internal/config/`
- **Zap structured logging** with `request_id` field propagated from middleware
- **Monorepo with `cmd/<service>/main.go`** layout — one binary per service, all sharing `internal/` packages
- **Taskfile.yml** as task runner — `dev:api`, `dev:web`, `dev:all`, `build:go`, `test:go`, `lint:go`, `migrate`, etc.
- **Helm chart with `templates/` per resource** — follow `kyros/deploy/helm/kyros/templates/` pattern
- **Database migrations** in `kyros/database/migrations/<timestamp>_<name>.sql` (timestamp-prefixed for ordering)

### Integration Points

- **API server** (`cmd/api/main.go`) listens on `:8080`, exposes `/v1/*` REST routes under chi
- **Registry** (`cmd/registry/main.go`) wraps `cncf/distribution` registry with Kyros auth middleware
- **Scanner/SBOM/Signer** workers consume NATS JetStream subjects `kyros.scan`, `kyros.sbom`, `kyros.sign`
- **Notifier** worker emits Trust Score updates + CVE alerts to the dashboard via WebSocket
- **Builder** submits Kaniko pods via Kubernetes API; tolerates spot taint
- **Next.js dashboard** at `apps/web/` calls API server; ShadCN components from `packages/ui/`

</code_context>

<specifics>
## Specific Ideas

- **URLs (locked)**: `kyros.smapatticare.com` (dashboard) + `registry.kyros.smapatticare.com` (OCI)
- **Production future (deferred)**: `kyros.io` + `r.kyros.io` — keep that in mind when designing any URL field in the data model
- **Source code lives inside the Homelab monorepo** at `kyros/` — co-located with Terraform, ArgoCD manifests, and the platform services
- **Trust Score is a "weighted composite, not a label"** — the user has stated this is the primary differentiator; the UI must show the score breakdown (each factor's contribution) not just the final number
- **No revenue for MVP** — every feature ships ungated; no Stripe / billing / paywall work
- **Familiar Kaniko pattern from CI world** — the user has used Kaniko in Jenkins JNLP agents in earlier phases; reuse the same image: `gcr.io/kaniko-project/executor:debug`

## Maintainer Notes (LLM handoff)

The downstream LLM/agent building the MVP from this document should:

1. **Start with `12-MAINTAINER.md`** (separate file) — it is the condensed, dense, LLM-optimized handoff. This CONTEXT.md is the human-readable audit + decision log; MAINTAINER.md is the build sheet.
2. **Respect the locked decisions D-01 through D-05** — these resolve the 3 HIGH-severity review findings; do not relitigate them
3. **Phase 12 is NOT in `ROADMAP.md`** at the time of writing — the roadmap table ends at Phase 11. If you generate a `gsd-plan-phase` workflow, be aware the planner may need the phase added to the roadmap first
4. **No `kyros-` requirement IDs exist yet** in `REQUIREMENTS.md` — generating REQ-IDs is part of the planning step, not this context step
5. **The Kyros `kyros/` subdirectory already has working monorepo scaffolding** — `cmd/`, `internal/`, `database/migrations/`, `apps/web/`, `deploy/helm/`, `infrastructure/terraform/`, `Taskfile.yml`, `docker-compose.dev.yml`. Plan-mode agents should audit what's already there before proposing new files
6. **`set -euo pipefail` in the dev shell** — see cross-distro curl-pipe installer caveat in `~/.hermes/memory`; applies to bootstrap scripts and `Taskfile.yml` shell-out commands
7. **Public OCI is a deferred idea, NOT a future phase** — only revisit when budget, quota, and abuse-mitigation all have concrete plans

</specifics>

<deferred>
## Deferred Ideas

- **Public OCI hosting + self-serve org signup** — post-MVP, requires abuse-mitigation plan (CAPTCHA, manual review queue, content scanning on upload, Azure Blob egress budget)
- **Grype dual-scanner** — re-evaluate post-MVP if Trivy false-positive rate is unacceptable
- **SLSA L3+ (hermetic builds)** — requires isolated builder infrastructure
- **Marketplace / billing** — post-MVP, requires Stripe + SaaS billing infrastructure
- **Multi-cloud storage (MinIO, S3, GCS)** — post-MVP, requires per-provider storage adapters
- **Multi-region deployment** — post-MVP, requires geo-replicated Postgres
- **Helm Chart repo hosting** — out of scope for Phase 12
- **WASM modules + AI Model artifacts** — out of scope
- **OCI Policies** — out of scope
- **kyros.io production domain** — only after MVP proves traction
- **apko/melange declarative image builds** — alternative to Kaniko for the future
- **Public trust score leaderboard / org rankings** — potential differentiator post-MVP
- **Webhook subscriptions for trust score changes** — useful but not MVP

### Reviewed Todos (not folded)

None.

</deferred>

---

*Phase: 12-kyros-platform*
*Context gathered: 2026-07-12*
*Revised: 2026-07-12 (post-cross-AI review + discuss-phase)*
