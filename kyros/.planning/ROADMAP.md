# Kyros — Roadmap
<!-- GSD:ROADMAP v1 -->

**Project:** Kyros — The Trusted Software Supply Chain Platform  
**Goal:** MVP in 3 months — OCI registry + scanning + publisher dashboard on AKS homelab  
**Version:** v1.0 (MVP)

---

## Milestone 1: MVP — "Foundation & Registry" (3 months)

---

### Phase 1: Monorepo Foundation
**Goal:** Kyros monorepo is scaffolded, builds cleanly, and all services have working skeletons with health endpoints
**Requirements:** REG-05 (partial), AUTH-01 (partial), OBS-02, OBS-03 (scaffolding)
**Depends on:** Nothing (greenfield start)
**Success Criteria:**
1. `kyros/` monorepo directory structure matches PROJECT.md spec
2. `go mod tidy` runs without errors in all Go packages
3. `npm install` runs without errors in `apps/web/` and `packages/`
4. `docker build` succeeds for each Go service producing a valid OCI image
5. All Go services expose `GET /health`, `GET /live`, `GET /ready` returning 200
6. Next.js dev server starts at `localhost:3000` with no errors
7. GitHub Actions CI workflow passes (lint + build + test)

**Plans:**
- [ ] 01-01: Monorepo scaffold — directory structure, Go workspace, Turborepo config, shared packages
- [ ] 01-02: Go service skeletons — kyros-api, kyros-registry proxy, worker stubs with health endpoints
- [ ] 01-03: Next.js app scaffold — App Router layout, ShadCN setup, Tailwind config, dark mode default
- [ ] 01-04: CI pipeline — GitHub Actions (lint, go test, npm test, docker build, OCI conformance smoke test)
- [ ] 01-05: PostgreSQL schema v1 — golang-migrate setup, initial schema (users, orgs, repos, tags, scans)

---

### Phase 2: OCI Registry Core
**Goal:** Full OCI Distribution Spec v1.1 compliant push/pull working end-to-end with Azure Blob Storage
**Requirements:** REG-01, REG-02, REG-03, REG-04, REG-05
**Depends on:** Phase 1
**Success Criteria:**
1. `docker push registry.kyros.smapatticare.com/{org}/{image}:{tag}` succeeds from local machine
2. `docker pull registry.kyros.smapatticare.com/{org}/{image}:{tag}` succeeds and image runs
3. OCI Distribution Spec conformance test suite passes (all required endpoints)
4. Multi-arch push (`linux/amd64` + `linux/arm64`) works via `docker buildx`
5. Pushed blobs are stored in Azure Blob Storage container `kyros-registry-blobs`
6. Referrers API `GET /v2/{name}/referrers/{digest}` returns empty manifest list (endpoint exists)
7. Anonymous pull for public repos works without auth header

**Plans:**
- [ ] 02-01: cncf/distribution deployment on AKS — Helm chart, Azure Blob storage driver config, HTTPRoute `registry.kyros.smapatticare.com`
- [ ] 02-02: Keycloak auth proxy — Bearer token validation middleware, token challenge flow, `www-authenticate` header
- [ ] 02-03: Registry metadata API — image push webhook → PostgreSQL (capture digest, size, media type, push timestamp)
- [ ] 02-04: OCI Referrers API support — enable in distribution config, test with `cosign attach`
- [ ] 02-05: End-to-end test — push real image, verify blobs in Azure Blob, conformance test suite pass

---

### Phase 3: Identity & Multi-Tenancy
**Goal:** Organizations, Teams, and RBAC fully functional; API tokens working for CI/CD authentication
**Requirements:** AUTH-01, AUTH-02, AUTH-03, AUTH-04
**Depends on:** Phase 2 (registry must be running for access control to be meaningful)
**Success Criteria:**
1. User can sign in at `kyros.smapatticare.com` via Keycloak SSO
2. User can create an organization with a unique slug
3. User can invite a second user to their org and assign a role
4. `org:viewer` user can pull public images but cannot push
5. `org:developer` user can push images to their org's namespace
6. API token generated from UI can be used as `docker login` password successfully
7. JWT validation works on all protected endpoints (401 without token, 403 with insufficient role)

**Plans:**
- [ ] 03-01: Next.js auth flow — Keycloak OIDC integration (next-auth or custom), session management, protected routes
- [ ] 03-02: Organization service — CRUD for orgs, teams, members; PostgreSQL tables + REST API
- [ ] 03-03: RBAC middleware — role extraction from JWT groups claim, per-endpoint permission check
- [ ] 03-04: API token service — token generation (random 64-byte hex), scopes, storage, revocation
- [ ] 03-05: Registry access control — org namespace authorization in distribution auth proxy

---

### Phase 4: Security Pipeline
**Goal:** Every pushed image is automatically scanned, SBOM generated, and Trust Score calculated within 90 seconds
**Requirements:** SEC-01, SEC-02, SEC-03, SEC-04 (verify only), SEC-05, SEC-06
**Depends on:** Phase 2 (registry), Phase 1 (NATS JetStream)
**Success Criteria:**
1. NATS JetStream deployed on AKS with `kyros.scans` and `kyros.sboms` subjects
2. Trivy scan triggered automatically within 60s of `docker push` completing
3. Grype scan runs in parallel with Trivy, results deduplicated by CVE ID
4. Trust Score calculated and stored within 90s of image push
5. SBOM generated in SPDX 2.3 JSON and CycloneDX 1.5 JSON formats
6. SBOMs attached to image manifest via Referrers API (OCI artifact)
7. Cosign `verify` command works for images signed before push

**Plans:**
- [ ] 04-01: NATS JetStream deployment — Helm chart, `kyros.scans` stream config, push event publisher in registry hook
- [ ] 04-02: kyros-scanner worker — Trivy + Grype runners, result deduplication, PostgreSQL persistence
- [ ] 04-03: kyros-sbom worker — Syft runner (SPDX + CycloneDX), OCI artifact attachment via Referrers API
- [ ] 04-04: Trust Score engine — Go package `internal/trust`, 6-signal weighted calculation, auto-recalculate triggers
- [ ] 04-05: Cosign verification integration — verify signatures on pull (policy: warn if unsigned), link Rekor log entry

---

### Phase 5: Image Builder
**Goal:** Developer can submit a Dockerfile via UI or API and get a scanned, signed image pushed to their org namespace
**Requirements:** BUILD-01, BUILD-02, BUILD-03
**Depends on:** Phase 3 (auth for build authorization), Phase 4 (scanning auto-runs post-build)
**Success Criteria:**
1. `POST /api/v1/orgs/{slug}/builds` accepts Dockerfile + context, returns build ID
2. Build pod spawned on AKS spot node pool with spot toleration within 30s of submission
3. Real-time build logs stream to browser via WebSocket
4. Completed image pushed to `registry.kyros.smapatticare.com/{org}/{image}:{tag}`
5. Scan triggered automatically after build completes (reuses Phase 4 pipeline)
6. Build pod cleaned up (deleted) within 60s of completion
7. Build history visible on publisher dashboard with status + duration

**Plans:**
- [ ] 05-01: kyros-builder worker — NATS consumer, Kubernetes Job spawner for Kaniko pods, spot toleration config
- [ ] 05-02: Kaniko pod template — resource limits (2 CPU, 4Gi RAM), spot toleration, Azure Blob credentials via Workload Identity
- [ ] 05-03: Build API — `POST /builds`, status tracking, WebSocket log streaming (`GET /builds/{id}/logs`)
- [ ] 05-04: Build UI — Dockerfile upload, build status page, real-time log viewer component

---

### Phase 6: Dashboard UI
**Goal:** World-class web UI — landing page, image search, image detail, org management, publisher dashboard, vulnerability explorer, SBOM explorer
**Requirements:** UI-01, UI-02, UI-03, UI-04, UI-05, UI-06, UI-07, UI-08
**Depends on:** Phase 2 (registry data), Phase 4 (scan results + SBOM), Phase 5 (build history)
**Success Criteria:**
1. `kyros.smapatticare.com` loads landing page in <2s (LCP)
2. Search returns results for any image pushed in the registry
3. Image detail page shows Trust Score badge + CVE table + SBOM package list for any scanned image
4. Organization creation flow works end-to-end (create org → invite member → member can pull)
5. SBOM explorer shows dependency tree for at least one test image
6. Dark mode is default; light mode toggle works and persists
7. Lighthouse score ≥ 90 on Performance, Accessibility, Best Practices, SEO

**Plans:**
- [ ] 06-01: Design system — ShadCN component customization, color tokens, typography, dark/light mode CSS vars
- [ ] 06-02: Landing page + Search — hero section, search input (Elasticsearch API), results grid with Trust Score badges
- [ ] 06-03: Image detail page — all metadata sections, Trust Score breakdown, CVE table, SBOM tab, Cosign tab
- [ ] 06-04: Organization management UI — create org, invite flow, member list, team management, settings
- [ ] 06-05: Publisher dashboard + Build UI — repo list, build trigger, build log viewer, scan history
- [ ] 06-06: Vulnerability explorer + SBOM explorer — CVE browser, SBOM dependency tree, license breakdown

---

### Phase 7: Observability & Search
**Goal:** All services fully instrumented with OpenTelemetry; Elasticsearch search index live; Grafana dashboards for all Kyros services
**Requirements:** OBS-01, OBS-02, OBS-03, OBS-04, OBS-05
**Depends on:** Phases 1–6 (all services must exist before instrumenting)
**Success Criteria:**
1. Distributed trace visible in Grafana Tempo spanning: HTTP request → API handler → NATS publish → worker → DB query
2. All services' `/metrics` endpoints scraped by Prometheus
3. Grafana shows Registry throughput dashboard (req/s, error rate, P95 latency)
4. Grafana shows Build pipeline dashboard (queue depth, active builds, success rate)
5. Elasticsearch index populated for all pushed images (searchable within 30s of push)
6. Log search in Grafana (Loki) returns structured logs for any service by trace_id

**Plans:**
- [ ] 07-01: OTLP instrumentation — add OpenTelemetry SDK to all Go services and Next.js app; configure Alloy OTLP receiver
- [ ] 07-02: Prometheus metrics — add RED metrics to all Go HTTP handlers; NATS consumer lag metrics; build pipeline metrics
- [ ] 07-03: Elasticsearch deployment + indexing — deploy OpenSearch/ES on AKS; index writer subscribes to push events; search API
- [ ] 07-04: Grafana dashboards — 4 dashboards as ConfigMaps in GitOps repo; Kyros-specific alert rules

---

### Phase 8: Deployment & ArgoCD Integration
**Goal:** Kyros fully GitOps-managed by ArgoCD on the homelab cluster; all services deployed via Helm; HTTPRoutes configured and TLS valid
**Requirements:** All requirements from Phases 1–7 deployed to AKS
**Depends on:** Phase 7 (all services complete)
**Success Criteria:**
1. ArgoCD Application `kyros` shows `Healthy + Synced`
2. All Kyros services running with 0 CrashLoopBackOff pods
3. `https://kyros.smapatticare.com` loads with valid TLS
4. `https://registry.kyros.smapatticare.com` serves OCI registry with valid TLS
5. `docker push registry.kyros.smapatticare.com/` works from the internet (not just localhost)
6. Cloudflare CDN serving manifests from cache (CF-Cache-Status: HIT on second request)
7. ArgoCD auto-sync enabled with 3-minute poll interval

**Plans:**
- [ ] 08-01: Kyros Helm chart — values.yaml, templates for all Deployments, Services, ConfigMaps, HPA
- [ ] 08-02: ArgoCD Application manifest + HTTPRoutes — `kyros.smapatticare.com` + `registry.kyros.smapatticare.com` HTTPRoutes with TLS
- [ ] 08-03: Cloudflare CDN config — cache rules for `/v2/{name}/manifests/` (Cache-Control headers), Workers route for blob redirect
- [ ] 08-04: End-to-end validation — push real image from CI, verify scan completes, Trust Score displayed, pull works globally

---

## Progress

| Phase | Plans | Status | Completed |
|---|---|---|---|
| 1. Monorepo Foundation | 0/5 | Not started | — |
| 2. OCI Registry Core | 0/5 | Not started | — |
| 3. Identity & Multi-Tenancy | 0/5 | Not started | — |
| 4. Security Pipeline | 0/5 | Not started | — |
| 5. Image Builder | 0/4 | Not started | — |
| 6. Dashboard UI | 0/6 | Not started | — |
| 7. Observability & Search | 0/4 | Not started | — |
| 8. Deployment & ArgoCD | 0/4 | Not started | — |

---

## Future Milestones (Post-MVP)

### Milestone 2: Helm Charts + Artifact Expansion
- Helm Chart repository (OCI-based, via `helm push`)
- WASM module hosting as OCI artifacts
- CLI tool: `kyros push`, `kyros pull`, `kyros verify`

### Milestone 3: Platform Growth
- SLSA L3 (hosted build platform with tamper-resistant provenance)
- apko + melange support (declarative minimal image builds)
- Organization-level policy engine (OPA rules for admission)
- kyros.io production domain + global CDN

### Milestone 4: Enterprise
- Billing + subscription management
- SSO federation (connect enterprise IdP to Kyros)
- Audit log export (CSV, SIEM integration)
- Private registry mirroring (cache Docker Hub, GHCR, ECR)

### Milestone 5: Marketplace & SaaS
- AI Model OCI artifact support
- OCI Policy distribution
- Verified Publisher marketplace
- Usage-based billing (per pull, per build minute)
