# Phase 12: Kyros — Maintainer Handoff (LLM Build Sheet)

> **Audience:** The next Claude/agent instance that is asked to plan, build, or extend Phase 12 of this homelab project.
> **Purpose:** A condensed, dense, decision-complete reference. No background, no selling — just the locked decisions, the executable plan order, the files that already exist, the files you must create, the traps to avoid, and the verification commands to run.
> **Companion docs:** `12-CONTEXT.md` (human-readable audit + decision log), `12-REVIEWS.md` (cross-AI review, all HIGH items resolved here).

---

## 1. WHAT YOU ARE BUILDING

Kyros — an OCI registry + publisher dashboard + Trust Score engine + baseline security pipeline, deployed on the existing AKS homelab cluster.

- **Dashboard:** `https://kyros.smapatticare.com` (Next.js 15)
- **Registry:** `https://registry.kyros.smapatticare.com` (cncf/distribution v3)
- **Differentiator:** Trust Score (weighted composite, 0-100)
- **MVP scope:** Private-only (no public OCI, no self-serve signup)
- **3-month MVP, 1 engineer, $0 additional compute cost** (rides on homelab)

---

## 2. LOCKED DECISIONS (DO NOT RE-LITIGATE)

| ID | Decision | Why |
|---|---|---|
| **D-01** | **Search backend = Typesense** (not Elasticsearch) | Single Go binary, ~50-100 MB RAM, official Helm chart, cloud-native. Replaces the 1-2 GB Elasticsearch heap. |
| **D-02** | **Private-only OCI for MVP** | No public repos, no self-serve signup. Mitigates $100-budget DoS risk. |
| **D-03** | **CronJob-based registry GC** + **per-org Azure Blob quota enforced in Go API before upload** | cncf/distribution does NOT do background GC. Daily 03:00 IST: lock writes → `registry garbage-collect` → unlock. |
| **D-04** | **Trivy + Syft + Cosign keyless only** (drop Grype + SLSA L3+ for MVP) | 3-feature pipeline, not 5. Cosign needs outbound HTTPS to sigstore hosts (egress allowlist needed). |
| **D-05** | **Postgres RLS for multi-tenancy** | Every multi-tenant table has `tenant_id` + RLS policy. Go API sets `app.current_tenant` GUC per txn. App-layer WHERE kept as defense-in-depth. |

**Stack (locked):**
- Backend: Go 1.23 + chi + viper + zap + prometheus/client_golang + go-containerregistry + oras-go
- Registry: cncf/distribution v3 + Azure Blob
- Frontend: Next.js 15 + TS + Tailwind + ShadCN
- Auth: Keycloak OIDC (Platform realm, `sso.smapatticare.com`)
- DB: Postgres 17 (RLS)
- Cache: Redis 7
- Search: **Typesense** (D-01)
- Queue: NATS JetStream
- Storage: Azure Blob
- Build: Kaniko (on spot node pool)
- Obs: LGTM (Prometheus + Grafana + Loki + Tempo + Alloy — already running from Phase 8)

**Trust Score formula (LOCKED — unchanged):**
- CVE severity: 35% | SBOM completeness: 20% | SLSA Level: 20% | Cosign signature: 10% | Scan freshness: 10% | Base image provenance: 5%
- Badges: Trusted ≥80, Moderate 50-79, At Risk <50

---

## 3. CODE THAT ALREADY EXISTS — AUDIT BEFORE WRITING

The monorepo is at `kyros/`. The scaffolding is non-trivial. Run before planning:

```bash
cd kyros
find . -type f \( -name "*.go" -o -name "*.ts" -o -name "*.tsx" -o -name "*.sql" -o -name "*.yaml" -o -name "*.yml" \) \
  -not -path "./node_modules/*" -not -path "./bun.lock" \
  | head -100
```

**Known existing layout:**

| Path | Purpose | Status |
|---|---|---|
| `kyros/go.mod` | Go module + deps | Exists |
| `kyros/cmd/api/main.go` | API entrypoint | Exists |
| `kyros/cmd/registry/main.go` | Registry service | Exists |
| `kyros/cmd/{builder,scanner,sbom,signer,notifier}/main.go` | 5 security workers | Exists |
| `kyros/internal/registry/{auth.go,auth_test.go,proxy.go}` | Registry auth glue | Exists |
| `kyros/internal/oci/` | OCI types | Exists |
| `kyros/internal/storage/` | Azure Blob adapter | Exists |
| `kyros/internal/auth/` | Keycloak OIDC | Exists |
| `kyros/internal/trust/` | Trust Score engine | Exists |
| `kyros/internal/httpserver/` | chi middleware chain | Exists |
| `kyros/internal/{config,logger,metrics,health}/` | Cross-cutting | Exists |
| `kyros/database/migrations/` | SQL migrations | Directory exists; check file count |
| `kyros/deploy/helm/kyros/` | Helm chart | Directory exists; check file count |
| `kyros/infrastructure/terraform/` | TF for Kyros-specific Azure | Directory exists; check file count |
| `kyros/docker-compose.dev.yml` | Local dev: Postgres + Redis + NATS + API + Web | Exists |
| `kyros/Taskfile.yml` | Task runner | Exists |
| `kyros/apps/web/` | Next.js dashboard | Exists |
| `kyros/packages/{config,types,ui}/` | Shared TS packages | Exists |
| `kyros/docs/{api,architecture,runbooks}/` | Doc dirs (likely empty) | Check |

`git status -s` shows heavy uncommitted work in `kyros/`. **Trust but verify** — read every file before deciding to keep/modify/replace.

---

## 4. EXECUTION ORDER (12-week MVP)

| Wk | Sub-phase | Concrete deliverables | Verification |
|---|---|---|---|
| w1-2 | **12.1 Foundation** | (a) Audit existing `kyros/` monorepo and commit any uncommitted work to a feature branch. (b) Write `12-01-PLAN.md` + `12-01-CONTEXT.md` (decisions scoped to foundation). (c) Provision Azure Blob container for registry + Typesense snapshots. (d) Deploy cncf/distribution via Helm on AKS. (e) HTTPRoute for `registry.kyros.smapatticare.com`. (f) Keycloak OIDC clients for all 6 services + dashboard. (g) Postgres schema with RLS policies. (h) First migration applied. | `docker login registry.kyros.smapatticare.com` succeeds; `docker push` of `hello-world` works; GET `/v2/_catalog` returns the catalog. |
| w3-4 | **12.2 Registry Core** | (a) Tag/manifest CRUD API. (b) Basic Next.js dashboard with org list + repo list. (c) **Typesense indexer** — every push triggers a NATS event consumed by an indexer worker that writes to Typesense. (d) **Per-org quota enforcement** in Go API (D-03). (e) org/team/repo CRUD with RLS. | `docker push` to a 60 GB blob returns 413 when org quota = 50 GB. `psql` as a tenant shows only that tenant's rows even if WHERE is omitted (RLS works). |
| w5-7 | **12.3 Security Pipeline** | (a) NATS JetStream streams + consumers. (b) Trivy scanner worker. (c) Syft SBOM worker. (d) Trust Score engine wired to scanner/sbom/signer outputs. (e) Trust Score UI page (show the breakdown, not just the number). (f) **Enforce private-only access** — every repo is private; no public flag in the API (D-02). | Push an image → NATS message → Trivy scan completes → SBOM emitted → Trust Score appears in dashboard. Push `alpine:3.19` then `apk add curl && commit` → second push updates the score. |
| w8-10 | **12.4 Publisher Dashboard** | (a) Org/team mgmt UI (invite, role change, remove). (b) Image builder UI → Kaniko pod on spot node pool (toleration set). (c) Cosign keyless signing worker. (d) **GC CronJob** (D-03) + manual GC runbook. (e) Webhook subscriptions (deferred to post-MVP — note as such). | Invite a user via Keycloak → they can log in. Build an image via UI → Kaniko completes → signed push → Trust Score = 100 (no CVEs in base alpine). GC CronJob successfully runs and deletes an untagged blob. |
| w11-12 | **12.5 Polish** | (a) LGTM instrumentation for all 6 services (Prometheus metrics, Grafana dashboard for Trust Score freshness/scanner latency/signer queue depth). (b) Typesense search UI in dashboard. (c) Cloudflare cache rules for `/v2/` static catalog. (d) Landing page. (e) Runbooks in `kyros/docs/runbooks/`. | `curl -I https://kyros.smapatticare.com` returns 200 + Cloudflare cache headers. Runbook "GC failed" → on-call can resolve in <15 min. |

---

## 5. EXECUTION RULES (NON-NEGOTIABLE)

1. **Respect D-01 through D-05.** If a user later asks for Elasticsearch, public OCI, Grype, or a non-RLS multi-tenancy model, push back with the rationale in `12-REVIEWS.md` first.
2. **Phase 12 is NOT in `ROADMAP.md`.** If the planner refuses to run because the phase isn't in the roadmap, add it first (the format is: `### Phase 12: Kyros` with goal/depends/requirements/success-criteria blocks, then a Plans bullet list).
3. **No `kyros-` requirement IDs in REQUIREMENTS.md.** If the planner asks for REQ-IDs, generate them (suggested prefix `KYROS-01` through `KYROS-NN`) and add a Traceability row to the table.
4. **Audit before write.** Read every existing file in `kyros/` before creating a new one. The monorepo is non-empty.
5. **Use existing patterns.** chi router, viper config, zap logger, prometheus metrics — they're already in `internal/`. Don't introduce echo/gin/cobra/zap alternatives unless there's a reason.
6. **Helm chart lives in `kyros/deploy/helm/kyros/`.** One values file per environment (`values-dev.yaml`, `values-prod.yaml`).
7. **Terraform for Kyros-specific Azure** lives in `kyros/infrastructure/terraform/`. The homelab platform Terraform is at the repo root (`modules/`, `main.tf`). Don't mix.
8. **Dev loop:** `task dev:all` brings up Postgres + Redis + NATS in Docker Compose + Go API + Next.js. Live-reload on both sides.
9. **RLS testing:** every migration that adds a multi-tenant table must include a test that proves the RLS policy blocks cross-tenant reads. Use `psql -c "SET app.current_tenant = '<other-tenant-uuid>'; SELECT * FROM repositories;"` and assert zero rows.

---

## 6. TRAPS / LANDMINES

- **cncf/distribution does not background-GC.** Without the D-03 CronJob, blob storage grows forever. The first run after a few weeks of pushing will be slow and delete a lot.
- **GC requires read-only lock.** If you run `registry garbage-collect` while a push is in progress, you can corrupt the registry. The lock-then-wait-then-GC-then-unlock dance in D-03 is mandatory.
- **Cosign keyless needs Sigstore egress.** `oauth2.sigstore.dev`, `fulcio.sigstore.dev`, `rekor.sigstore.dev` on 443 must be in the cluster's egress allowlist. Phase 9 (NetworkPolicies) hasn't landed yet, so this is the NSG egress for now. Verify with `curl -I https://oauth2.sigstore.dev` from a debug pod.
- **Trivy DB updates.** Trivy downloads its vulnerability DB on every scanner pod start. Either pin a fresh image (Trivy publishes daily) or run a sidecar that updates a shared PVC. Plan for ~200 MB of CVE data.
- **Kaniko on spot = interruption risk.** Spot nodes can be evicted at any time. Set a `retryStrategy` on the Kaniko Job (e.g. `backoffLimit: 3`) and persist intermediate layers to Azure Blob.
- **Typesense single-node = single point of failure.** Acceptable for MVP. For HA, run 3 nodes with a shared PVC + a sidecar load balancer (out of scope).
- **Postgres RLS + connection pooling.** PgBouncer in transaction-pool mode breaks RLS because the GUC is per-session. Use session-pool mode, or set the GUC inside the transaction (`SET LOCAL app.current_tenant = ...`).
- **NATS JetStream persistence.** Set `storage: file` not `storage: memory` for the JetStream streams. Otherwise a NATS pod restart loses every in-flight scan/SBOM/signing job.
- **Keycloak token expiry on long scans.** A scan can take 5+ minutes. Renew the OIDC token at the worker start, not in the middle of a 10-min upload.
- **Azure Blob cold tier.** Default tier is Hot, which is expensive for rarely-accessed blobs (old tags, old SBOMs). Set a lifecycle policy: blobs older than 30 days → Cool, older than 90 days → Archive.

---

## 7. KEY FILES TO READ FIRST

```
kyros/README.md
kyros/Taskfile.yml
kyros/docker-compose.dev.yml
kyros/go.mod
kyros/cmd/api/main.go
kyros/internal/registry/auth.go
kyros/internal/trust/  (entire dir)
kyros/internal/oci/    (entire dir)
kyros/internal/storage/  (entire dir)
kyros/internal/auth/   (entire dir)
kyros/database/migrations/  (every file)
kyros/deploy/helm/kyros/   (every file)
kyros/apps/web/app/page.tsx
kyros/infrastructure/terraform/  (every file)
.planning/PROJECT.md
.planning/REQUIREMENTS.md
.planning/phases/12-kyros-platform/12-CONTEXT.md
.planning/phases/12-kyros-platform/12-REVIEWS.md
```

---

## 8. VERIFICATION COMMANDS (run at end of each sub-phase)

```bash
# 12.1 — Registry reachable + OCI push works
docker login registry.kyros.smapatticare.com
docker pull hello-world
docker tag hello-world registry.kyros.smapatticare.com/kyrosdev/hello-world:v1
docker push registry.kyros.smapatticare.com/kyrosdev/hello-world:v1
curl -s https://registry.kyros.smapatticare.com/v2/_catalog | jq

# 12.2 — Quota enforcement
# (push a >50 GB blob via dd + docker import + push, expect 413)
psql -h localhost -U kyros -d kyros -c "SET app.current_tenant = '00000000-0000-0000-0000-000000000001'; SELECT count(*) FROM repositories;"
# Should be 0 if the test user has no repos (RLS working)

# 12.3 — Trust Score
# After push + scan, expect dashboard shows Trust Score with breakdown
curl -s https://kyros.smapatticare.com/v1/repositories/kyrosdev/hello-world/trust-score | jq

# 12.4 — Cosign verification
cosign verify --certificate-identity-regexp '.*' --certificate-oidc-issuer-regexp '.*' \
  registry.kyros.smapatticare.com/kyrosdev/hello-world:v1

# 12.4 — GC
kubectl create job --from=cronjob/kyros-registry-gc manual-gc -n kyros
kubectl logs -n kyros job/manual-gc -f
# Expect: blobs deleted, no errors

# 12.5 — Observability
curl -s http://kyros-prometheus:9090/api/v1/query?query=trust_score_freshness_seconds | jq
curl -s http://kyros-grafana:3000/api/dashboards/uid/kyros-trust-score | jq .dashboard.title
```

---

## 9. OPEN QUESTIONS FOR THE NEXT LLM

These are not blocking, but the user may ask. Be ready with a recommendation.

1. **Where does Kyros CI run?** Options: (a) reuse the homelab Jenkins from Phase 7, (b) new Tekton, (c) GitHub Actions. Recommendation: (a) — Jenkins is already running, JNLP agents tolerate spot, and the user has muscle memory.
2. **How are orgs created?** Options: (a) CLI tool + Keycloak admin, (b) admin-only web form, (c) Terraform. Recommendation: (b) for MVP simplicity, behind a separate `org:platform-admin` role.
3. **Backups?** Postgres (daily pg_dump to Azure Blob), Typesense snapshots (nightly, D-01), registry blobs (Azure Blob native soft-delete + GRS replication).
4. **Monitoring alerts?** Node memory > 85% (already from Phase 8), GC failure, Trivy DB older than 48h, scan queue depth > 100, signer queue depth > 50, failed sign rate > 5%.
5. **What happens if Sigstore is down?** Cosign keyless signing will fail. Decide: queue and retry? Skip signing and mark Trust Score with a `signer_unavailable` flag? Recommendation: queue and retry for 1h, then fail with a clear error.

---

*Maintained as part of Phase 12. Update when D-01 through D-05 change or when the execution order shifts.*
*Last updated: 2026-07-12*
