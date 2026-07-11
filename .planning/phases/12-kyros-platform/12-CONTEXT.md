# Phase 12: Kyros — The Trusted Software Supply Chain Platform
<!-- GSD:CONTEXT v1 -->
Date: 2026-07-12

## Domain
Kyros is a cloud-native, OCI-compliant software supply chain platform. The MVP delivers:
a fully functional OCI registry (push/pull) + vulnerability scanning UI + publisher dashboard,
deployed on the existing AKS homelab cluster.
- Dashboard: kyros.smapatticare.com
- Registry: registry.kyros.smapatticare.com
Target: 3-month MVP with 5 security pipeline features fully integrated from day 1.

## Decisions

### Product Identity
- Kyros is simultaneously: image BUILDER + trusted REGISTRY + publisher PLATFORM
- Primary differentiator: the Trust Score (weighted composite, not a label)
- Long-term positioning: "The GitHub of trusted software artifacts"

### Business Model
- No revenue for MVP — build traction first
- Future: Free public tier + paid private repos + enterprise SaaS

### Deployment
- Cloud-native SaaS on existing AKS homelab cluster
- MVP URLs: kyros.smapatticare.com (dashboard) + registry.kyros.smapatticare.com (OCI)
- Production future: kyros.io + r.kyros.io

## Tech Stack (LOCKED)
- Backend: Go (go-containerregistry, oras-go)
- Registry Engine: cncf/distribution (Docker Distribution v3) + Azure Blob backend
- Frontend: Next.js 15 + TypeScript + Tailwind + ShadCN
- Auth: Keycloak (reuse homelab Platform realm)
- Database: PostgreSQL (metadata) + Elasticsearch (search) + Redis (cache)
- Storage: Azure Blob Storage
- Queue: NATS JetStream
- CDN: Cloudflare
- Build: Kaniko (K8s-native, ephemeral pods on spot nodes)
- Observability: Full LGTM stack (Prometheus + Grafana + Loki + Tempo + Alloy)

## Security Pipeline (ALL in MVP)
1. Vulnerability Scanning: Trivy + Grype (dual-scanner)
2. SBOM Generation: Syft -> SPDX + CycloneDX
3. Signing: Cosign (keyless via Sigstore/Fulcio + Rekor)
4. SLSA Provenance: L2 provenance for Kyros-built images
5. Image Building: Kaniko-based pipeline

## Trust Score (LOCKED)
Composite weighted score 0-100:
- CVE severity: 35% (0 Critical = max contribution)
- SBOM completeness: 20% (both SPDX + CycloneDX present)
- SLSA Level: 20% (L1=25%, L2=50%, L3=75%, L4=100%)
- Cosign signature: 10%
- Scan freshness: 10% (scanned in last 24h = max)
- Base image provenance: 5%
Badges: Trusted >=80, Moderate 50-79, At Risk <50

## Multi-Tenancy
- Structure: registry.kyros.smapatticare.com/{org}/{image}:{tag}
- Entities: User, Organization, Team, Repository, Tag, Digest
- RBAC: org:owner, org:admin, org:devops, org:developer, org:viewer
- Public + private images from day 1

## Source Code Location
- Inside Homelab repo: Homelab/kyros/ (monorepo sub-project)
- GSD-managed project: separate /gsd-new-project for Kyros

## MVP Roadmap (12 weeks)
- Phase 12.1 (w1-2): Foundation - monorepo, Postgres schema, Keycloak SSO, distribution deployment
- Phase 12.2 (w3-4): Registry Core - docker push/pull, manifest/tag API, basic UI
- Phase 12.3 (w5-7): Security Pipeline - NATS, scanner, SBOM, Trust Score engine + UI
- Phase 12.4 (w8-10): Publisher Dashboard - org/team mgmt, image builder, Cosign signing
- Phase 12.5 (w11-12): Polish - LGTM instrumentation, search, CDN, landing page, docs

## Deferred Ideas
- Helm Chart repo, WASM modules, AI Model artifacts, OCI Policies
- Marketplace/billing
- Multi-cloud storage (MinIO, S3, GCS)
- Multi-region deployment
- apko/melange declarative minimal image builds
- kyros.io production domain

## Canonical Refs
- OCI Distribution Spec v1.1: https://github.com/opencontainers/distribution-spec
- cncf/distribution: https://github.com/distribution/distribution
- go-containerregistry: https://github.com/google/go-containerregistry
- Sigstore/Cosign: https://github.com/sigstore/cosign
- Syft: https://github.com/anchore/syft
- Trivy: https://github.com/aquasecurity/trivy
- Kaniko: https://github.com/GoogleContainerTools/kaniko
- NATS JetStream: https://docs.nats.io

## Code Context (Homelab Reusables)
- Keycloak: sso.smapatticare.com (Platform realm) - reuse for Kyros user auth
- Traefik Gateway: HTTPRoute + ReferenceGrant pattern established
- cert-manager: TLS auto-provisioned for *.smapatticare.com
- ArgoCD: Kyros as ArgoCD Application (GitOps)
- LGTM stack: Prometheus + Grafana + Loki + Tempo + Alloy already running
- AKS spot node pool: available for Kaniko build jobs
- Azure Blob Storage: already provisioned
- Cloudflare DNS: manages smapatticare.com
