# Kyros

> **"The Trusted Software Supply Chain Platform"**

## What This Is

Kyros is a cloud-native, OCI-compliant software supply chain platform that provides:

1. **A trusted OCI registry** — push, pull, and discover container images with full security metadata
2. **A hardened image building pipeline** — submit a Dockerfile or apko YAML, Kyros builds it securely using Kaniko on ephemeral Kubernetes pods
3. **A publisher platform** — organizations publish verified images under their verified publisher badge
4. **A security intelligence layer** — every image has a Trust Score, vulnerability scan results, SBOM, and provenance attestation

The long-term vision: become the **GitHub of trusted software artifacts** — starting with OCI container images, expanding to Helm Charts, WASM modules, AI models packaged as OCI artifacts, and more.

## Core Value

**Developers and organizations deserve to know exactly what they're running.** Every image on Kyros is scanned, signed, and scored. No more blindly trusting `docker pull ubuntu:latest`.

## Tech Stack

| Layer | Technology | Rationale |
|---|---|---|
| **Frontend** | Next.js 15 + TypeScript + Tailwind + ShadCN | Industry standard for SaaS dashboards, App Router, SSR |
| **Backend API** | Go | Dominant in OCI/cloud-native tooling; native OCI library ecosystem |
| **Registry Engine** | `cncf/distribution` v3 | OCI Distribution Spec compliant; Azure Blob backend driver built-in |
| **Auth** | Keycloak (Platform realm) | SSO/OIDC/OAuth2; already running in homelab |
| **Primary DB** | PostgreSQL | Battle-tested metadata store |
| **Search** | Elasticsearch | Full-text + faceted search for image discovery |
| **Cache** | Redis | Token cache, rate limiting, scan result cache, manifest cache |
| **Blob Storage** | Azure Blob Storage | OCI layer/blob store; cncf/distribution native driver |
| **Message Queue** | NATS JetStream | Cloud-native, K8s-native, lightweight async job processing |
| **CDN** | Cloudflare | Global manifest/layer distribution; already managing DNS |
| **Build Engine** | Kaniko | Kubernetes-native, rootless, ephemeral pod builds on spot nodes |
| **Observability** | Full LGTM stack | Prometheus + Grafana + Loki + Tempo + Alloy (all already running) |

## Services Architecture

```
kyros.smapatticare.com (Next.js UI)
        │
        ▼
  [Traefik Gateway]
        │
        ├── /api/v1/*           → kyros-api (Go, REST)
        ├── /v2/*               → kyros-registry (cncf/distribution + auth proxy)
        └── (future: /graphql)

kyros-api (Go microservices):
  ├── Auth service          ← Keycloak OIDC
  ├── Organization service  ← Orgs, Teams, Members
  ├── Repository service    ← Repos, Tags, Manifests
  ├── Scan result service   ← Vulnerability data
  ├── SBOM service          ← SPDX + CycloneDX artifacts
  ├── Trust Score service   ← Weighted composite calculation
  ├── Publisher service     ← Verified publisher workflow
  └── Builder service       ← Kaniko job orchestration

Workers (NATS JetStream consumers):
  ├── kyros-scanner     (Trivy + Grype)
  ├── kyros-sbom        (Syft)
  ├── kyros-signer      (Cosign)
  ├── kyros-builder     (Kaniko job spawner)
  └── kyros-notifier    (webhooks/email)
```

## Monorepo Structure

```
kyros/
├── apps/
│   ├── web/                    # Next.js 15 dashboard
│   └── api/                    # Go REST API (main entrypoint)
├── services/
│   ├── registry/               # cncf/distribution wrapper + auth proxy
│   ├── scanner/                # Trivy + Grype worker
│   ├── sbom/                   # Syft worker (SPDX + CycloneDX)
│   ├── signer/                 # Cosign signing worker
│   ├── builder/                # Kaniko job spawner
│   └── notifier/               # Webhook/email notifications
├── packages/
│   ├── types/                  # Shared TypeScript types/interfaces
│   ├── ui/                     # Shared ShadCN component library
│   └── config/                 # Shared ESLint + TypeScript configs
├── internal/                   # Shared Go packages
│   ├── auth/                   # Keycloak OIDC middleware
│   ├── oci/                    # OCI spec utilities
│   ├── trust/                  # Trust Score engine
│   └── storage/                # Azure Blob abstraction layer
├── database/
│   ├── migrations/             # golang-migrate PostgreSQL migrations
│   └── schema.sql              # Current schema snapshot
├── deployments/
│   ├── helm/kyros/             # Kyros Helm chart
│   └── argocd/                 # ArgoCD Application manifests
├── infrastructure/
│   └── terraform/              # Azure resources for Kyros
├── tests/
│   ├── e2e/                    # Playwright end-to-end tests
│   ├── integration/            # Go integration tests
│   └── load/                   # k6 load tests
├── docs/
│   ├── api/                    # OpenAPI specs
│   ├── architecture/           # Architecture diagrams
│   └── runbooks/               # Operational runbooks
├── scripts/                    # Dev/build/deploy helper scripts
└── .github/
    └── workflows/              # CI: lint, test, build, security scan
```

## Deployment

| Environment | URL | Notes |
|---|---|---|
| Dashboard (MVP) | `https://kyros.smapatticare.com` | AKS homelab cluster |
| Registry (MVP) | `https://registry.kyros.smapatticare.com` | OCI Distribution endpoint |
| Dashboard (Production) | `https://kyros.io` | Future — when production-ready |
| Registry (Production) | `https://r.kyros.io` | Future |

## Trust Score Algorithm

Composite weighted score (0–100) displayed as badge:

| Signal | Weight | Logic |
|---|---|---|
| CVE severity | 35% | 0 Critical = full contribution; each Critical –8pts |
| SBOM completeness | 20% | Both SPDX + CycloneDX present = full contribution |
| SLSA Level | 20% | L4=100%, L3=75%, L2=50%, L1=25%, none=0% |
| Cosign signature | 10% | Signed + Rekor-logged = full contribution |
| Scan freshness | 10% | Scanned ≤24h = 100%, ≤7d = 50%, older = 0% |
| Base image provenance | 5% | Base image is distroless/minimal/trusted = full |

**Badges:** 🟢 Trusted (≥80) · 🟡 Moderate (50–79) · 🔴 At Risk (<50)

## Requirements

### Validated

(None yet — first phase underway)

### Active

**Registry Core**
- [ ] REG-01: OCI Distribution Spec v1.1 compliant push/pull (docker push/pull works)
- [ ] REG-02: Multi-arch manifest support (linux/amd64, linux/arm64)
- [ ] REG-03: OCI Referrers API support (attach SBOMs/signatures to manifests)
- [ ] REG-04: Anonymous pull for public images; bearer token required for private + writes
- [ ] REG-05: Azure Blob Storage backend for all blobs/layers

**Identity & Multi-Tenancy**
- [ ] AUTH-01: Keycloak OIDC integration (Platform realm) for user authentication
- [ ] AUTH-02: Organization + Team + Namespace model (`registry.kyros.smapatticare.com/{org}/{image}:{tag}`)
- [ ] AUTH-03: RBAC roles: org:owner, org:admin, org:devops, org:developer, org:viewer
- [ ] AUTH-04: API token generation + management (personal access tokens)

**Security Pipeline**
- [ ] SEC-01: Trivy vulnerability scan triggered on every push
- [ ] SEC-02: Grype vulnerability scan (dual-scanner, deduplicated results)
- [ ] SEC-03: Syft SBOM generation (SPDX + CycloneDX formats)
- [ ] SEC-04: Cosign signature verification (keyless via Sigstore/Fulcio + Rekor)
- [ ] SEC-05: SLSA L2 provenance attestation for Kyros-built images
- [ ] SEC-06: Trust Score calculation and display on every image page

**Image Building**
- [ ] BUILD-01: Kaniko-based build pipeline (Kubernetes-native, rootless)
- [ ] BUILD-02: Dockerfile upload and build triggering via UI and API
- [ ] BUILD-03: Build status streaming (real-time log output)
- [ ] BUILD-04: Ephemeral build pods on AKS spot node pool

**Dashboard UI**
- [ ] UI-01: Landing page (public-facing, shows featured images + Trust Scores)
- [ ] UI-02: Image search (full-text + faceted: OS, arch, language, Trust Score, signed/distroless)
- [ ] UI-03: Image detail page (all metadata, CVEs, SBOM, pull command, Trust Score badge)
- [ ] UI-04: Organization management (create org, invite members, set roles)
- [ ] UI-05: Publisher dashboard (manage repos, trigger builds, view scan results)
- [ ] UI-06: Vulnerability explorer (browse CVEs with affected images)
- [ ] UI-07: SBOM explorer (browse SPDX/CycloneDX with dependency tree)
- [ ] UI-08: Dark mode + Light mode

**Observability**
- [ ] OBS-01: All services emit OpenTelemetry traces (Tempo)
- [ ] OBS-02: All services expose /metrics, /health, /live, /ready endpoints
- [ ] OBS-03: Grafana dashboards for all Kyros services (RED metrics)
- [ ] OBS-04: Loki log aggregation for all services
- [ ] OBS-05: Alloy agent collecting all telemetry

### Out of Scope (MVP)

- Helm Chart repository hosting — Future Phase 2
- WASM module hosting — Future Phase 3
- AI Model OCI artifact support — Future Phase 4
- OCI Policy distribution — Future Phase 4
- Billing / subscription management — Post-traction
- Marketplace revenue model — Post-traction
- Multi-cloud storage (MinIO, S3, GCS) — Post-MVP, after Azure Blob proven
- Multi-region deployment — Post-MVP
- CockroachDB migration — Only if horizontal scaling needed
- apko/melange declarative build support — Phase 2 builder enhancement
- kyros.io production domain + global CDN — Post-MVP launch
- GitHub Actions integration — Phase 2 CI/CD integration
- Grafana Cloud managed observability — Self-hosted is free and already running

## Key Decisions

| Decision | Rationale | Outcome |
|---|---|---|
| Use `cncf/distribution` as registry engine | Don't build OCI spec from scratch. Industry-standard, Azure Blob native driver, battle-tested. Wrap with Kyros auth proxy + metadata hooks | Locked |
| Go for backend | Dominant in OCI tooling ecosystem. go-containerregistry, oras-go, containerd libs all native Go. Best HTTP streaming performance for blob push/pull | Locked |
| Keycloak for auth | Already running in homelab. OIDC/OAuth2 standard. Multi-realm. No additional infra cost | Locked |
| Trust Score as differentiator | Not a label (Platinum/Silver). A weighted composite (0–100) based on 6 quantifiable security signals. Transparent algorithm published in docs | Locked |
| Kaniko for image builds | Kubernetes-native, rootless, no Docker daemon. Each build = ephemeral pod. Inherently more secure than Docker-in-Docker (DinD) approaches | Locked |
| Cloudflare CDN | Already managing smapatticare.com DNS. Free tier sufficient for MVP. Workers can cache manifests at edge | Locked |
| NATS JetStream over Kafka | Kafka is overkill for MVP. NATS JetStream is purpose-built for K8s, lightweight, zero-ops for homelab scale | Locked |
| Elasticsearch over Typesense | Chosen by user. Note: Heavy ops overhead (1–2GB RAM minimum). Revisit at Phase 3 if resource pressure | Locked (revisit at Phase 3) |
| Public + private images from day 1 | Enterprise expectation. Anonymous pull for public repos. Bearer token for private + all writes | Locked |
| Full LGTM from day 1 | All stack components already running in homelab. Zero additional deployment cost. Correct-by-construction observability | Locked |

## Context

**Deployment target (MVP):**
- AKS homelab cluster (`homelab-aks`, `homelab-demo` namespace)
- Managed by ArgoCD (Kyros deployed as ArgoCD Application, GitOps from `Homelab/kyros/deployments/argocd/`)
- TLS via cert-manager + Let's Encrypt (wildcard `*.smapatticare.com`)
- DNS via Cloudflare (CNAME records to Traefik LoadBalancer IP)
- Ingress via Traefik Gateway API (HTTPRoute + ReferenceGrant)

**Shared homelab resources Kyros reuses:**
- Keycloak `sso.smapatticare.com` — Platform realm for user auth
- Traefik Gateway — HTTPRoute + ReferenceGrant pattern already established
- cert-manager — TLS auto-provisioned
- ArgoCD — GitOps deployment
- LGTM stack — Prometheus + Grafana + Loki + Tempo + Alloy all running
- AKS spot node pool — Available for Kaniko build jobs (spot toleration)
- Azure Blob Storage — Already provisioned

**OCI ecosystem references:**
- OCI Distribution Spec v1.1: https://github.com/opencontainers/distribution-spec
- OCI Image Spec v1.1: https://github.com/opencontainers/image-spec
- cncf/distribution: https://github.com/distribution/distribution
- go-containerregistry: https://github.com/google/go-containerregistry
- oras-go: https://github.com/oras-project/oras-go
- Sigstore/Cosign: https://github.com/sigstore/cosign
- Syft: https://github.com/anchore/syft
- Trivy: https://github.com/aquasecurity/trivy
- Grype: https://github.com/anchore/grype
- Kaniko: https://github.com/GoogleContainerTools/kaniko
- SLSA Framework: https://slsa.dev
- NATS JetStream: https://docs.nats.io/nats-concepts/jetstream

## Constraints

- **Budget:** Homelab AKS budget ($12–25/month total) — Kyros must fit within the existing cluster, no new Azure resources except namespaces
- **No new Kubernetes clusters:** Run within existing `homelab-demo` or new `kyros` namespace on same AKS cluster
- **Kubernetes-first:** Everything runs in K8s. No VMs, no serverless-only components
- **CNCF principles:** Follow CNCF recommendations for all tooling choices
- **OCI compliance:** Registry MUST pass OCI Distribution Spec conformance tests

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-07-12 after initialization*
