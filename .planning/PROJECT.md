# Homelab AKS Platform

## What This Is

A cost-optimized, production-patterned Azure Kubernetes Service homelab running 24x7 within a $100 budget for 6-8 months (~$12-25/month). The platform runs ArgoCD for GitOps-based deployments, Jenkins with ephemeral JNLP build agents, and a full observability stack — all secured via Workload Identity Federation, private networking, and Azure Key Vault for secrets. Provisioned entirely with Terraform.

## Core Value

A real, running Kubernetes platform that teaches production-grade patterns (GitOps, zero-trust secrets, HA networking) without breaking the bank — the cluster must stay alive 24x7 within the $100 total budget.

## Requirements

### Validated

(None yet — ship to validate)

### Active

<!-- Infrastructure Layer -->
- [ ] AKS cluster provisioned via Terraform in Central India region
- [ ] AKS Free tier control plane (free, no SLA premium)
- [ ] 1 reserved Standard_D2as_v5 system node pool in Zone 1 (8 GB RAM, 2 vCPU)
- [ ] Spot Standard_D2as_v5 user node pool spanning Zone 1+2, max 4 nodes, scale-to-zero
- [ ] AKS CNI Overlay networking (pod IPs are overlay, not consuming VNet CIDR)
- [ ] Azure VNet 10.0.0.0/16 with subnets: system (10.0.1.0/24), user/spot (10.0.2.0/24), infra (10.0.4.0/24), ingress (10.0.8.0/24)
- [ ] AKS API server authorized IP ranges (home IP whitelist, no jumpbox required)
- [ ] No public IPs on worker nodes — nodes are fully private
- [ ] Azure Key Vault with Private Endpoint in infra-subnet
- [ ] Workload Identity Federation enabled on AKS (OIDC issuer + workload identity)
- [ ] Secrets Store CSI Driver + Azure Key Vault provider (no static credentials anywhere)
- [ ] Azure Container Registry (optional, or use Docker Hub free tier)

<!-- Platform Services via GitOps/Helm -->
- [ ] ArgoCD (GitOps controller) — manages all platform and workload deployments
- [ ] Jenkins controller with Kubernetes plugin (JNLP ephemeral agents, scale-to-zero)
- [ ] Traefik v3 (replaces ingress-nginx EOL March 2026) — type: LoadBalancer, single public IP, **Gateway API as primary routing API** (HTTPRoute, GatewayClass)
- [ ] Kubernetes Gateway API CRDs v1.5.1 installed (GatewayClass, Gateway, HTTPRoute, ReferenceGrant)
- [ ] cert-manager with Let's Encrypt ClusterIssuer (automatic TLS for all services)
- [x] kube-prometheus-stack (Prometheus + Grafana + Alertmanager, scrape 30s interval) (via grafana/k8s-monitoring Alloy)
- [x] External Secrets Operator or Secrets Store CSI to surface AKV secrets as K8s secrets

<!-- Security & Networking -->
- [ ] Cloudflare DNS proxy → Azure Public Load Balancer (Traefik v3) → cluster services
- [ ] All platform UIs (ArgoCD, Jenkins, Grafana) exposed via **`HTTPRoute`** (Gateway API) with TLS auto-issued by cert-manager via `Gateway` annotation
- [ ] `ReferenceGrant` per namespace enabling cross-namespace Gateway → Service binding
- [ ] NetworkPolicies: default-deny ingress/egress per namespace, allow only required paths
- [ ] Managed Identity per workload via Workload Identity (zero service principal credentials)
- [ ] Azure Key Vault RBAC: least-privilege per identity (Key Vault Secrets User role)

### Out of Scope

- ~~ingress-nginx (`kubernetes/ingress-nginx`)~~ — **EOL March 2026**, repo archived/read-only, no security patches. Replaced by Traefik v3
- Multi-region or cross-subscription setup — budget constraint makes this impractical
- Azure VPN Gateway — adds ~$27/mo (kills budget), replaced by authorized IP ranges
- Azure Bastion / Jumpbox VM — adds ~$4-8/mo; authorized IP ranges is sufficient for homelab
- B-series (burstable) nodes for system pool — Microsoft explicitly warns against this for AKS system pools
- True 3-zone HA — requires minimum 3 nodes (one per AZ), ≥$50/mo base cost, exceeds budget
- HashiCorp Vault — Azure Key Vault CSI with Workload Identity is zero-ops and effectively free at homelab scale
- Azure VPN / Private DNS Zones for AKS — not needed with authorized IP ranges approach
- ARM/Graviton nodes — not yet stable enough for mixed system+user AKS pools in this region

## Context

**Budget reality:**
- $100 total for 6-8 months = $12.50-16.67/month hard ceiling
- Standard_D2as_v5 1-year reserved (Central India): ~$17/mo (slightly over but spans 2 AZs better than B-series)
- Spot node pool (Zone 1+2): ~$3-8/mo when active for Jenkins builds (scale-to-zero when idle)
- Azure Key Vault: <$0.50/mo at homelab operation volumes
- Total estimated: ~$20-25/mo for 4-5 months OR ~$16-18/mo for 6 months if spot pool idle most of the time

**Node sizing rationale (system node 8 GB RAM):**
- OS + kubelet (~600 MB) + ArgoCD (~1.5 GB) + Jenkins controller (~512 MB)
- kube-prometheus-stack (~900 MB) + ingress/cert-manager/CSI (~320 MB) + buffer (~400 MB)
- Total: ~4.5 GB / 8 GB — comfortable headroom ✅
- Spot user nodes: Jenkins JNLP ephemeral build pods only (512 MB–1 Gi per agent, scale-to-zero)

**Workload Identity Federation flow:**
- AKS OIDC issuer → Azure AD Federated Credential → User-Assigned Managed Identity → AKV RBAC
- Zero static credentials: no Service Principal secrets in code, git, or cluster secrets
- Pattern applies to: ArgoCD image updater, External Secrets Operator, application workloads

**GitOps model:**
- ArgoCD is the single deployment operator — Terraform provisions infra, ArgoCD reconciles K8s manifests
- Helm charts stored in Git → ArgoCD ApplicationSets deploy platform services
- Jenkins pipelines build images → push to registry → update image tags in Git → ArgoCD syncs

**Traefik v3 + Gateway API + cert-manager + Cloudflare strategy:**
- Traefik v3 (CNCF, Helm chart 41.0.1) as the Gateway API implementation
- **Kubernetes Gateway API v1.5** (GA, Feb 2026) as primary routing API — `HTTPRoute` not `IngressRoute` CRDs
- cert-manager `v1.15+` with `--set config.enableGatewayAPI=true` — annotate `Gateway` with `cert-manager.io/cluster-issuer` for auto-TLS
- Cloudflare DNS-01 challenge for wildcard cert (`*.yourdomain.com`) — single cert covers all subdomains
- Cloudflare SSL mode: "Full (strict)" — encrypts both edge-to-origin leg
- `ReferenceGrant` in each app namespace grants Traefik Gateway cross-namespace access

**VNet design:**
- VNet: 10.0.0.0/16 | System subnet: 10.0.1.0/24 | User/spot subnet: 10.0.2.0/24
- Infra subnet: 10.0.4.0/24 (AKV private endpoint) | Ingress subnet: 10.0.8.0/24
- Pod CIDR (overlay): 10.244.0.0/16 (does not consume VNet addresses)
- API server access: Authorized IP ranges (home IP) — no VPN Gateway needed

## Constraints

- **Budget**: $100 total / ~$12-25/month — every Azure resource choice must be justified against this
- **Region**: Central India — most cost-effective Azure region for South Asia with AZ support
- **Control plane**: AKS Free tier only — no Standard/Premium uptime SLA needed for homelab
- **Node SKUs**: D-series only for system pool — B-series explicitly not recommended by Microsoft for AKS system node pools
- **Terraform**: All Azure infra via Terraform (azurerm provider) — no ClickOps
- **No exposed credentials**: Workload Identity Federation only — no Service Principal client secrets anywhere
- **K8s version**: Latest stable AKS minor version at time of deploy (currently 1.31.x)

## Key Decisions

| Decision | Rationale | Outcome |
|---|---|---|
| AKS CNI Overlay | Saves VNet address space vs Azure CNI; pod IPs are overlay; supports large pod counts on small subnets | — Pending |
| Authorized IP ranges (not private cluster + VPN) | Private cluster adds ~$27/mo VPN Gateway cost; authorized IP ranges achieves 95% of security benefit at $0 | — Pending |
| Standard_D2as_v5 over B-series | Microsoft explicitly advises against B-series for AKS system pools; D2as_v5 AMD is cheapest D-series with AZ support | — Pending |
| 1-year reserved instance for system node | ~50% discount vs pay-go; system node runs 24x7 so reserved is pure savings | — Pending |
| AKV CSI + Workload Identity over HashiCorp Vault | Zero infra to operate, effectively free, teaches production-grade zero-trust pattern | — Pending |
| 2-zone (Zone 1+2) not 3-zone | 3-zone requires 3 system nodes minimum; doubles/triples base cost beyond budget | — Pending |
| Traefik v3 over ingress-nginx | `kubernetes/ingress-nginx` EOL March 2026 — archived, no patches. Traefik v3 is CNCF-backed, 100% Gateway API conformance | — Pending |
| Gateway API over IngressRoute CRD | Gateway API v1.5 is GA (Feb 2026), standard across all controllers, no Traefik vendor lock-in, cert-manager native support via `Gateway` annotation, `ReferenceGrant` for multi-tenancy | — Pending |
| Cloudflare free tier as edge | Free DDoS protection + CDN + DNS — no Azure Front Door or Application Gateway cost | — Pending |
| ArgoCD as sole deployment operator | Single source of truth for K8s state; Terraform handles infra, ArgoCD handles K8s manifests | — Pending |

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
*Last updated: 2026-07-02 after initialization*
