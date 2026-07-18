# Kyros

> **"The Trusted Software Supply Chain Platform"**

## Overview

Kyros is a cloud-native, OCI-compliant software supply chain platform that provides hardened container images, vulnerability scanning, SBOM generation, and cryptographic signing. Every image hosted on Kyros has a **Trust Score** — a transparent, weighted composite of security signals.

**MVP URLs:**
- Dashboard: `https://kyros.iemafzalhassan.tech`
- Registry: `https://registry.kyros.iemafzalhassan.tech`

## Stack

| Layer | Technology |
|---|---|
| Frontend | Next.js 15 + TypeScript + Tailwind + ShadCN |
| Backend | Go (go-containerregistry, oras-go) |
| Registry Engine | cncf/distribution v3 + Azure Blob backend |
| Auth | Keycloak OIDC (Platform realm) |
| Database | PostgreSQL + Elasticsearch + Redis |
| Queue | NATS JetStream |
| CDN | Cloudflare |
| Build Engine | Kaniko (K8s-native, ephemeral pods) |
| Observability | LGTM stack (Prometheus + Grafana + Loki + Tempo + Alloy) |

## Quick Start

```bash
# Clone
git clone <repo-url>
cd kyros

# Start dev dependencies
docker-compose up -d postgres redis nats

# Start API
cd apps/api && go run ./cmd/server

# Start web
cd apps/web && npm run dev
```

## Project Structure

```
kyros/
├── apps/           # Next.js UI + Go API
├── services/       # Worker microservices (scanner, sbom, signer, builder, notifier)
├── packages/       # Shared TypeScript packages (types, ui, config)
├── internal/       # Shared Go packages (auth, oci, trust, storage)
├── database/       # PostgreSQL migrations + schema
├── deployments/    # Helm chart + ArgoCD manifests
├── infrastructure/ # Terraform (Azure resources)
├── tests/          # E2E (Playwright), Integration, Load (k6)
├── docs/           # API specs, architecture diagrams, runbooks
└── .planning/      # GSD project planning (ROADMAP.md, phases, etc.)
```

## Planning

This project is managed with GSD. See [.planning/ROADMAP.md](.planning/ROADMAP.md) for phases and progress.

```bash
# Check current status
/gsd-progress

# Plan next phase
/gsd-plan-phase 1
```

## Security Pipeline

Every image push triggers:
1. **Trivy + Grype** — dual vulnerability scanning
2. **Syft** — SBOM generation (SPDX 2.3 + CycloneDX 1.5)
3. **Trust Score** — weighted composite (CVE severity 35%, SBOM 20%, SLSA 20%, Signing 10%, Freshness 10%, Provenance 5%)
4. **Cosign** — signature verification

## License

Apache 2.0
