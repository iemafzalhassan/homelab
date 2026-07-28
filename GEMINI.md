<!-- GSD:project-start source:PROJECT.md -->
## Project

**Kube Telemetry Stage**

A production-patterned Kubernetes observability and Grafana demo environment for CNCF virtual speaking events. The platform runs ArgoCD for GitOps-based deployments, Jenkins with ephemeral JNLP build agents, and a full observability stack — all secured via Workload Identity Federation, private networking, and Azure Key Vault for secrets. Provisioned entirely with Terraform.

**Core Value:** A real, running Kubernetes platform that teaches production-grade patterns (GitOps, zero-trust secrets, HA networking) without breaking the bank — the cluster must stay alive 24x7 within the $100 total budget.

### Constraints

- **Budget**: $100 total / ~$12-25/month — every Azure resource choice must be justified against this
- **Region**: Central India — most cost-effective Azure region for South Asia with AZ support
- **Control plane**: AKS Free tier only — no Standard/Premium uptime SLA needed for homelab
- **Node SKUs**: D-series only for system pool — B-series explicitly not recommended by Microsoft for AKS system node pools
- **Terraform**: All Azure infra via Terraform (azurerm provider) — no ClickOps
- **No exposed credentials**: Workload Identity Federation only — no Service Principal client secrets anywhere
- **K8s version**: Latest stable AKS minor version at time of deploy (currently 1.31.x)
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Verified Component Versions
| Component | Helm Chart Version | App Version | Notes |
|---|---|---|---|
| **Terraform azurerm provider** | `4.79.0` | — | Latest stable |
| **AKS Kubernetes** | `1.35.x` | — | Stable GA; 1.36 is latest GA |
| **ArgoCD** | `10.1.0` | `v3.4.4` | argo/argo-cd |
| **Jenkins** | `5.9.32` | `2.555.3` | jenkinsci/jenkins |
| ~~**ingress-nginx**~~ | ~~`4.15.1`~~ | ~~`1.15.1`~~ | ⛔ **EOL March 2026** — repo archived, no security patches |
| **Traefik** | `41.0.1` | `v3.6+` | traefik/traefik — replaces ingress-nginx, supports Gateway API |
| **cert-manager** | `1.20.3` | `v1.20.3` | jetstack/cert-manager |
| **kube-prometheus-stack** | `87.5.1` | `v0.92.1` | prometheus-community/kube-prometheus-stack |
| **Secrets Store CSI Driver** | `1.6.0` | `1.6.0` | Requires K8s >= 1.30 |
| **AKV CSI Provider** | latest | — | csi-secrets-store-provider-azure |
## Recommended Terraform Module Structure
## Key Terraform Snippets
### AKS with CNI Overlay + Workload Identity
### Federated Identity Credential
## Rationale
- **D2as_v5 not B-series**: Microsoft explicitly advises against B-series for AKS system node pools — CPU credit exhaustion causes CoreDNS and kubelet failures
- **CNI Overlay**: Pod IPs don't consume VNet address space; supports large pod counts on /24 subnets
- **Secrets Store CSI v1.6.0**: Requires K8s 1.30+; pairs with AKV provider for zero-credential secret access
- **ArgoCD 10.x (v3.x)**: Major version with improved ApplicationSet support and multi-source Apps
- **Traefik v3 + Gateway API over ingress-nginx + IngressRoute**: `kubernetes/ingress-nginx` EOL March 2026. Traefik v3 with Kubernetes **Gateway API** (`HTTPRoute`) is the correct modern choice — GA (v1.5 Feb 2026), 100% Traefik conformance, no vendor CRD lock-in, cert-manager native support
## ⛔ ingress-nginx EOL — Decision Log
| What is EOL | What is NOT EOL |
|---|---|
| `kubernetes/ingress-nginx` (community) | Kubernetes Ingress API itself |
| — No security patches | `nginxinc/kubernetes-ingress` (F5, still active) |
| — No K8s version compatibility updates | Traefik, Kong, Envoy Gateway, etc. |
### Why Gateway API > IngressRoute CRD
| | Traefik IngressRoute CRD | **Kubernetes Gateway API** |
|---|---|---|
| Standard | Traefik-proprietary | ✅ Upstream Kubernetes (GA v1.5) |
| cert-manager | Manual `secretName` | ✅ Annotate `Gateway` → auto-issued |
| Portability | Traefik-only | ✅ Works on any conformant controller |
| Role separation | None | ✅ `GatewayClass` / `Gateway` / `HTTPRoute` |
| Multi-tenancy | Weak | ✅ `ReferenceGrant` cross-namespace |
| Future | Traefik maintains it | ✅ Kubernetes SIG-Network standard |
### Gateway API CRDs Install (required before Traefik)
# Gateway API v1.5.1 — Standard channel (stable resources only)
### Traefik v3 Helm values (Gateway API first-class)
### cert-manager integration with Gateway API
# cert-manager >= v1.15 with Gateway API enabled
# helm upgrade cert-manager ... --set config.enableGatewayAPI=true
# Annotate the Gateway to auto-trigger certificate issuance:
### Application routing pattern (HTTPRoute)
### Cross-namespace ReferenceGrant (required for Gateway API)
# In each app namespace: allows traefik Gateway to bind HTTPRoutes
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->\n## Developer Profile

**Name:** Afzal Hassan
**Role:** DevOps Engineer & CNCF Community Leader
**Goals:** Mastery in Cloud Native Architecture, Platform Engineering, and Community Leadership

### Engineering Philosophy
- **Production-Grade First:** Treat homelab and side projects with enterprise rigor.
- **GitOps & Automation:** Everything as code, no manual ClickOps.
- **Deep Understanding:** Focus on the "Why", "How", and Tradeoffs—not just the "What".
- **Maintainability & Open Source:** Prefer long-term maintainable, open-source solutions.

### Communication Style
- **Tone:** Humble, confident, respectful, concise, and engineering-oriented.
- **Vibe:** Friendly and professional; sounds like a real engineer, not corporate HR.
- **Mentorship:** Prefers learning through diagrams, analogies, step-by-step reasoning, and real production examples.

### Collaboration Protocol
- **Challenge Ideas:** Act as a Senior Staff/Principal reviewer—do not blindly agree.
- **Structured Problem Solving:** Context -> Reasoning -> Architecture -> Tradeoffs -> Best Practice -> Future Improvements.
- **Focus:** Build systems and workflows for productivity; optimize health (sleep, fitness) around long office hours.
\n<!-- GSD:profile-end -->
