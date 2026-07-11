# Kyros — Requirements
<!-- GSD:REQUIREMENTS v1 -->

## Project: Kyros — The Trusted Software Supply Chain Platform
**Code:** KYROS  
**Created:** 2026-07-12

---

## Epic 1: OCI Registry Core

### REG-01 — OCI Distribution Spec Compliance
**As a** developer,  
**I want** to `docker push` and `docker pull` images from `registry.kyros.smapatticare.com`,  
**So that** I can use Kyros as a drop-in replacement for Docker Hub or GHCR.

**Acceptance Criteria:**
- [ ] `docker push registry.kyros.smapatticare.com/{org}/{image}:{tag}` succeeds
- [ ] `docker pull registry.kyros.smapatticare.com/{org}/{image}:{tag}` succeeds
- [ ] All OCI Distribution Spec v1.1 conformance tests pass
- [ ] Manifest `GET /v2/{name}/manifests/{reference}` returns correct Content-Type
- [ ] Blob upload via chunked PUT succeeds for layers >100MB

### REG-02 — Multi-Architecture Manifests
**As a** developer,  
**I want** to push and pull multi-arch image manifests (OCI Image Index),  
**So that** my images work on `linux/amd64` and `linux/arm64` without separate tags.

**Acceptance Criteria:**
- [ ] `docker buildx push --platform linux/amd64,linux/arm64` succeeds
- [ ] Manifest list (OCI Image Index) stored and retrieved correctly
- [ ] Platform-specific pull selects correct variant

### REG-03 — OCI Referrers API
**As a** security tool (Cosign, Syft),  
**I want** to attach attestations and SBOMs to image manifests via the Referrers API,  
**So that** supply chain metadata is co-located with the artifact it describes.

**Acceptance Criteria:**
- [ ] `GET /v2/{name}/referrers/{digest}` returns referrers manifest list
- [ ] Cosign can attach signatures via `cosign attach signature`
- [ ] SBOM attachments appear in Referrers response
- [ ] Referrers API filters by `artifactType` query param

### REG-04 — Access Control for Push/Pull
**Acceptance Criteria:**
- [ ] Anonymous pull works for public repositories (no auth header required)
- [ ] Authenticated pull required for private repositories (returns 401 without token)
- [ ] Push always requires authenticated user with `org:developer` role or higher
- [ ] `www-authenticate` challenge returned on 401 with correct Keycloak token URL

### REG-05 — Azure Blob Storage Backend
**Acceptance Criteria:**
- [ ] All blob layers stored in Azure Blob Storage (not local disk)
- [ ] `cncf/distribution` configured with Azure storage driver
- [ ] Blob deletion cascades correctly when image is deleted
- [ ] Storage bucket/container named `kyros-registry-blobs`

---

## Epic 2: Identity & Multi-Tenancy

### AUTH-01 — Keycloak OIDC Authentication
**As a** user,  
**I want** to sign in to `kyros.smapatticare.com` using my Keycloak account,  
**So that** I have one identity across all platform tools.

**Acceptance Criteria:**
- [ ] Sign-in redirects to `sso.smapatticare.com/realms/Platform` for authentication
- [ ] After auth, JWT returned contains `sub`, `email`, `groups` claims
- [ ] Token validated on every API request via Keycloak JWKS endpoint
- [ ] Token refresh handled automatically; session expires after 24h

### AUTH-02 — Organizations & Teams
**As an** organization admin,  
**I want** to create an organization with a unique slug (e.g., `acme`),  
**So that** my team's images are namespaced under `registry.kyros.smapatticare.com/acme/`.

**Acceptance Criteria:**
- [ ] `POST /api/v1/orgs` creates organization with unique slug
- [ ] `POST /api/v1/orgs/{slug}/teams` creates teams within org
- [ ] `POST /api/v1/orgs/{slug}/members` invites users to org
- [ ] Org slug is globally unique and immutable after creation

### AUTH-03 — RBAC Roles
**Acceptance Criteria:**
- [ ] `org:owner` — full control, can delete org
- [ ] `org:admin` — manage members, repos, settings; cannot delete org
- [ ] `org:devops` — push, delete images; cannot manage members
- [ ] `org:developer` — push images only
- [ ] `org:viewer` — pull private images, read-only access

### AUTH-04 — API Tokens
**As a** CI/CD pipeline,  
**I want** to generate a personal access token (PAT) for machine authentication,  
**So that** I can push images from GitHub Actions / Jenkins without using my password.

**Acceptance Criteria:**
- [ ] `POST /api/v1/tokens` generates token with configurable expiry
- [ ] Token scopes: `read:packages`, `write:packages`, `delete:packages`
- [ ] Token usable for `docker login` (as password)
- [ ] Token revocable from UI

---

## Epic 3: Security Pipeline

### SEC-01 — Trivy Vulnerability Scanning
**As a** developer,  
**I want** every pushed image automatically scanned with Trivy,  
**So that** I see CVE counts without having to run Trivy myself.

**Acceptance Criteria:**
- [ ] Scan triggered automatically within 60s of image push
- [ ] Scan results stored in PostgreSQL per digest
- [ ] CVE counts by severity (Critical, High, Medium, Low, Unknown) displayed on image page
- [ ] Affected packages + CVE IDs + CVSS scores surfaced
- [ ] Scan status: `pending` → `scanning` → `complete` / `failed`

### SEC-02 — Grype Dual Scanning
**Acceptance Criteria:**
- [ ] Grype scan runs in parallel with Trivy for same image
- [ ] Results deduplicated by CVE ID
- [ ] Source attribution shown (Trivy-only, Grype-only, or both)
- [ ] Discrepancies between scanners flagged with note

### SEC-03 — SBOM Generation (Syft)
**As a** security engineer,  
**I want** every image to have a machine-readable SBOM in both SPDX and CycloneDX formats,  
**So that** I can audit dependencies for license compliance and vulnerability exposure.

**Acceptance Criteria:**
- [ ] Syft generates SBOM in SPDX 2.3 JSON format
- [ ] Syft generates SBOM in CycloneDX 1.5 JSON format
- [ ] SBOMs stored as OCI artifacts attached via Referrers API
- [ ] SBOMs downloadable from image detail page
- [ ] Package count and detected languages shown on image page

### SEC-04 — Cosign Signature Verification
**Acceptance Criteria:**
- [ ] `cosign verify` command shown on image page with correct registry endpoint
- [ ] Signature presence shown as badge (Signed ✓ / Unsigned ✗)
- [ ] Rekor transparency log entry linked from image page
- [ ] Keyless signing via OIDC (Sigstore Fulcio) supported
- [ ] Key-based signing (traditional `cosign.key`) also accepted

### SEC-05 — SLSA Provenance
**Acceptance Criteria:**
- [ ] Kyros-built images include SLSA L2 provenance attestation
- [ ] Provenance includes: builder ID (Kyros), build timestamp, source repo, commit SHA
- [ ] Provenance stored as OCI attestation via Referrers API
- [ ] SLSA Level displayed on image detail page

### SEC-06 — Trust Score
**Acceptance Criteria:**
- [ ] Every image has a Trust Score (0–100) calculated from the 6 signals
- [ ] Trust Score updates automatically when scan completes or signature changes
- [ ] Score badge (🟢/🟡/🔴) visible on image cards in search results
- [ ] Score breakdown (per-signal) shown on image detail page
- [ ] Algorithm publicly documented at `/docs/trust-score`

---

## Epic 4: Image Builder

### BUILD-01 — Kaniko Build Pipeline
**As a** developer,  
**I want** to submit a Dockerfile to Kyros and have it build the image securely,  
**So that** I don't need a local Docker daemon or insecure DinD setup.

**Acceptance Criteria:**
- [ ] `POST /api/v1/builds` accepts Dockerfile + context (tarball or Git URL)
- [ ] Build executes in ephemeral Kaniko pod on AKS spot node pool
- [ ] Build pod has spot toleration + resource limits (2 CPU, 4Gi RAM)
- [ ] Completed image pushed to `registry.kyros.smapatticare.com/{org}/{image}:{tag}`
- [ ] Build pod deleted after completion (no persistent build artifacts)

### BUILD-02 — Build Triggering
**Acceptance Criteria:**
- [ ] Trigger build from UI (upload Dockerfile + build context)
- [ ] Trigger build from API with JSON payload
- [ ] Build queued via NATS JetStream (`kyros.builds` subject)
- [ ] Builder worker dequeues and spawns Kaniko pod

### BUILD-03 — Real-Time Build Logs
**Acceptance Criteria:**
- [ ] Build logs streamed in real-time via WebSocket on build detail page
- [ ] Log lines include timestamps and step labels (FROM, RUN, COPY, etc.)
- [ ] Final status (success/failure) shown with elapsed time
- [ ] Failed builds show error message + exit code prominently

---

## Epic 5: Dashboard UI

### UI-01 — Landing Page
**Acceptance Criteria:**
- [ ] Hero section with Kyros tagline + search bar
- [ ] Featured trusted images section (curated by admins)
- [ ] Statistics banner (total images, total scans, publishers)
- [ ] "Get Started" CTA linking to sign-up / docs
- [ ] Dark mode default; light mode toggle

### UI-02 — Image Search
**Acceptance Criteria:**
- [ ] Full-text search across image names, publishers, descriptions
- [ ] Facet filters: OS, Architecture, Language/Runtime, Trust Score tier, Signed (yes/no), Distroless (yes/no), Verified Publisher
- [ ] Results sorted by: relevance, Trust Score, pull count, last updated
- [ ] `docker pull` command pre-filled on hover/expand

### UI-03 — Image Detail Page
**Acceptance Criteria:**
- [ ] All metadata displayed: name, org, description, digest, tags, size, OS/arch, created date, last scan
- [ ] Trust Score badge + breakdown (per-signal)
- [ ] Pull command with one-click copy
- [ ] CVE table: severity, ID, package, version, CVSS, link to NVD
- [ ] SBOM tab: package list with name, version, license, type
- [ ] Signatures tab: Cosign signature status + Rekor link
- [ ] Provenance tab: SLSA level + build metadata
- [ ] All tags listed with digest + push date

### UI-04 — Organization Management
**Acceptance Criteria:**
- [ ] Create org with slug, display name, description, avatar
- [ ] Invite members by email with role assignment
- [ ] Member list with role + joined date + remove action
- [ ] Team management (create teams, assign members)
- [ ] Org settings (delete org, transfer ownership)

### UI-05 — Publisher Dashboard
**Acceptance Criteria:**
- [ ] List all repositories with: image count, last push, avg Trust Score
- [ ] Per-repo: tag list, build history, scan history, SBOM download
- [ ] Trigger new build from dashboard
- [ ] Verified publisher badge application flow

### UI-06 — Vulnerability Explorer
**Acceptance Criteria:**
- [ ] Browse all CVEs across all scanned images
- [ ] Filter by severity, CVE ID, package name, affected images
- [ ] CVE detail page with description, CVSS score, affected images list

### UI-07 — SBOM Explorer
**Acceptance Criteria:**
- [ ] Browse SBOM for any image (SPDX or CycloneDX view toggle)
- [ ] Package dependency tree visualization
- [ ] License breakdown pie chart
- [ ] Download SBOM as JSON or XML

### UI-08 — Dark / Light Mode
**Acceptance Criteria:**
- [ ] Dark mode is default
- [ ] Toggle in navigation bar persists preference
- [ ] System preference respected on first visit

---

## Epic 6: Observability

### OBS-01 — OpenTelemetry Traces
**Acceptance Criteria:**
- [ ] All Go services emit OTLP traces to Alloy collector
- [ ] All Next.js API routes emit OTLP traces
- [ ] Distributed traces span entire request lifecycle (UI → API → Worker → Registry)
- [ ] Trace sampling: 100% in dev, 10% in prod

### OBS-02 — Health Endpoints (all services)
**Acceptance Criteria:**
- [ ] `GET /health` returns `{"status":"ok"}` when healthy
- [ ] `GET /live` returns 200 when process is alive
- [ ] `GET /ready` returns 200 when dependencies (DB, Redis, NATS) are reachable
- [ ] Kubernetes liveness + readiness probes configured for all pods

### OBS-03 — Prometheus Metrics
**Acceptance Criteria:**
- [ ] All services expose `GET /metrics` in Prometheus text format
- [ ] RED metrics (Rate, Errors, Duration) for every HTTP handler
- [ ] Queue metrics: NATS subject lag, consumer pending count
- [ ] Build metrics: queue depth, build duration, success/failure rate
- [ ] Registry metrics: push/pull rate, blob cache hit rate, storage bytes

### OBS-04 — Loki Logging
**Acceptance Criteria:**
- [ ] All services emit structured JSON logs
- [ ] Log fields: level, timestamp, service, trace_id, span_id, message
- [ ] Alloy agent scrapes logs and ships to Loki
- [ ] Grafana dashboard includes log panel per service

### OBS-05 — Grafana Dashboards
**Acceptance Criteria:**
- [ ] Dashboard: Registry throughput (push/pull rate by org, error rate, latency P95)
- [ ] Dashboard: Build pipeline (queue depth, active builds, success rate)
- [ ] Dashboard: Security pipeline (scan queue lag, scans/min, CVE distribution)
- [ ] Dashboard: API gateway (req/s, error rate, latency by endpoint)

---

## Definition of Done (per Phase)

A phase is complete when:
1. All acceptance criteria checked off
2. Unit tests passing (≥80% coverage for new code)
3. Integration tests passing (happy path + error cases)
4. All services have `/health`, `/live`, `/ready` endpoints
5. All services emit OTLP traces + Prometheus metrics
6. ArgoCD Application shows `Healthy + Synced`
7. Manual smoke test: `docker push` + `docker pull` work end-to-end
8. Security scan shows no Critical CVEs in Kyros own images

---
*Last updated: 2026-07-12 after initialization*
