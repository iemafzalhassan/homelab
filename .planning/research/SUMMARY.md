# Research Summary — AKS Homelab Platform (2026)

Synthesized from: STACK.md, FEATURES.md, ARCHITECTURE.md, PITFALLS.md

## TL;DR

Build a 2-zone AKS homelab on Azure (Central India) running 24x7 on ~$20-25/month. One reserved D2as_v5 system node hosts ArgoCD, Jenkins controller, Traefik v3, cert-manager, and Prometheus. Spot D2as_v5 nodes (scale 0→4) handle ephemeral Jenkins JNLP builds. Everything uses Kubernetes Gateway API (`HTTPRoute`) and Workload Identity Federation — no static credentials, no legacy ingress, no deprecated components. Terraform provisions all Azure infra; ArgoCD reconciles all Kubernetes state.

---

## Stack Decision Summary

| Component | Choice | Why |
|---|---|---|
| **IaC** | Terraform + azurerm `4.79.0` | Industry standard, modullar |
| **Cluster** | AKS Free tier, K8s `1.35.x`, Central India | Free control plane, cheapest D-series region |
| **Networking** | CNI Overlay (`10.244.0.0/16` pod CIDR) | Saves VNet address space |
| **Node: system** | Standard_D2as_v5 reserved Zone 1 | ~$17/mo 1yr reserved, enough RAM (8GB), no B-series |
| **Node: spot** | Standard_D2as_v5 spot Zone 1+2, 0→4 | Scale-to-zero, Jenkins build workloads only |
| **Secrets** | AKV CSI + Workload Identity | Zero static creds, no etcd secret exposure |
| **Ingress** | Traefik v3 `41.0.1` + **Gateway API v1.5** | ingress-nginx EOL March 2026; Gateway API is GA standard |
| **TLS** | cert-manager `1.20.3` + Cloudflare DNS-01 | Auto wildcard cert via `Gateway` annotation |
| **GitOps** | ArgoCD `10.1.0` (v3.4.4) App-of-Apps | Manages all K8s state after infra bootstrap |
| **CI/CD** | Jenkins `5.9.32` + JNLP ephemeral pods | Scale-to-zero, spot node tolerations |
| **Observability** | kube-prometheus-stack `87.5.1` (30s scrape, 7d retention) | Full stack without LGTM overhead |

---

## Critical Architecture Decisions

### 1. Gateway API over IngressRoute
Kubernetes Gateway API v1.5 is GA (Feb 2026). Traefik v3 passes 100% conformance. cert-manager v1.15+ natively manages TLS via `Gateway` annotation. Using `HTTPRoute` instead of Traefik's proprietary `IngressRoute` CRDs means zero vendor lock-in and teaches the actual production standard.

### 2. Spot nodes only for Jenkins JNLP (never system components)
System components (ArgoCD, Traefik, Prometheus) live on the reserved system node exclusively. Spot pool is tainted — only Jenkins JNLP pods tolerate the taint. This prevents spot eviction from taking down the entire cluster.

### 3. Authorized IP ranges over private cluster + VPN Gateway
Private cluster adds ~$27/mo for VPN Gateway. Authorized IP ranges achieves equivalent access control for the homelab context. `make update-ip` refreshes the allowlist when home IP changes.

### 4. One Managed Identity per workload
Not a shared identity for the cluster. Each workload (ArgoCD, Jenkins, Traefik) gets its own User-Assigned Managed Identity with a federated credential scoped exactly to its ServiceAccount and namespace.

---

## Phase-Critical Pitfalls (Must Read Before Executing)

| Risk | When | Fix |
|---|---|---|
| Traefik installed before Gateway API CRDs | Phase 4/5 | Always `kubectl apply` CRDs first |
| cert-manager missing `enableGatewayAPI=true` | Phase 5 | Add `--set config.enableGatewayAPI=true` to helm install |
| No `ReferenceGrant` → HTTPRoutes rejected | Phase 5/6 | Deploy ReferenceGrant in every app namespace |
| Workload Identity SA/namespace mismatch | Phase 3+ | Triple-check `subject` field in federated credential |
| Spot eviction kills system components | Phase 2 | Add taint to spot pool; no tolerations on system workloads |
| No Azure budget alert | Phase 1 | Always deploy `azurerm_consumption_budget_resource_group` |

---

## Recommended Build Order (10 Phases)

```
Phase 1: Azure foundation (RG, VNet, subnets, NSGs, budget alert)
Phase 2: AKS cluster + system node pool (free tier, CNI overlay, OIDC on)
Phase 3: Managed Identities + federated credentials (per workload)
Phase 4: Key Vault + Private Endpoint + RBAC
Phase 5: Gateway API CRDs → Traefik v3 → cert-manager → ClusterIssuer
Phase 6: ArgoCD bootstrap (Helm) + App-of-Apps root application
Phase 7: Jenkins + JNLP (via ArgoCD) + spot node pool
Phase 8: kube-prometheus-stack + dashboards (via ArgoCD)
Phase 9: Platform hardening (NetworkPolicies, PDBs, LimitRanges, ReferenceGrants)
Phase 10: Validation + documentation + cost review
```

---
*Synthesized: July 2026*
