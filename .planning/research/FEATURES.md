# Features Research — AKS Homelab Platform (2026)

## Table Stakes (Must Have — Platform Won't Function Without These)

| Feature | Why Table Stakes | Implementation |
|---|---|---|
| AKS cluster (Free tier control plane) | Foundation — everything runs here | Terraform `azurerm_kubernetes_cluster` |
| Workload Identity Federation | No static creds anywhere — security baseline | AKS OIDC + `azurerm_federated_identity_credential` |
| Azure Key Vault + CSI Driver | Secret storage without etcd exposure | `azurerm_key_vault` + Secrets Store CSI |
| Traefik v3 + Gateway API | External traffic routing — without it nothing is reachable | Helm + Gateway API CRDs |
| cert-manager + Let's Encrypt | Automated TLS — manual cert management is a full-time job | Helm + Cloudflare DNS-01 |
| ArgoCD | GitOps reconciler — keeps cluster state in sync with Git | Helm via ArgoCD bootstrap |
| Jenkins + JNLP ephemeral agents | CI pipelines — builds, tests, image pushes | Helm + Kubernetes plugin |
| kube-prometheus-stack | Without observability, debugging is blind | Helm (Prometheus + Grafana + Alertmanager) |
| NetworkPolicies (default-deny) | Zero-trust pod isolation — table stakes for any secure cluster | Applied per-namespace via ArgoCD |

## Differentiators (What Makes This Homelab Stand Out)

| Feature | Value | Priority |
|---|---|---|
| Gateway API (`HTTPRoute`) not legacy Ingress | Teaches the actual 2026 standard, portable across controllers | High |
| Spot node pool scale-to-zero | Real cost optimization pattern used in production | High |
| Per-workload Managed Identity | Least-privilege at pod level — each workload has its own identity | High |
| ArgoCD App-of-Apps pattern | Single ArgoCD Application bootstraps entire platform | Medium |
| Pod Disruption Budgets on system components | Handles spot eviction gracefully | Medium |
| Azure Budget Alerts | Cost guardrails — get notified before overspending | Medium |
| Traefik dashboard via HTTPRoute + auth | Visibility into routing rules, useful for debugging | Low |

## Anti-Features (Deliberately NOT Building)

| Anti-Feature | Reason |
|---|---|
| Multi-cluster ArgoCD | One cluster is the right scope for a $100 budget homelab |
| Service Mesh (Istio/Linkerd) | Memory overhead (~500MB+) would blow the node budget; NetworkPolicies achieve the necessary isolation |
| Full LGTM stack (Loki + Tempo + Mimir) | Too heavy for 8GB system node; kube-prometheus-stack covers 80% of value at 20% of cost |
| Azure Firewall for egress | ~$300/month — completely out of scope |
| Kyverno/OPA policy engine | Adds complexity; focus on core platform first |
| Horizontal Pod Autoscaler on platform services | ArgoCD/Jenkins are stateful — HPA doesn't apply; VPA in recommendation mode is sufficient |
| Blue/Green deployments | ArgoCD rolling updates sufficient for homelab |

## Feature Dependencies (Build Order Implications)

```
1. Azure infrastructure (VNet, subnets, resource groups)
   └── 2. AKS cluster (needs VNet subnet IDs)
       ├── 3a. Workload Identity + Managed Identities (needs OIDC issuer URL)
       │    └── 4a. Key Vault + RBAC (needs managed identity client IDs)
       │         └── 5a. Secrets Store CSI (deployed to cluster)
       └── 3b. Gateway API CRDs (kubectl apply — no Helm)
            └── 4b. Traefik v3 (needs Gateway API CRDs pre-installed)
                 └── 5b. cert-manager (needs Traefik LB IP for DNS-01 challenge config)
                      └── 6. ArgoCD (manages all subsequent deployments)
                           ├── 7a. Jenkins + JNLP
                           └── 7b. kube-prometheus-stack
```

## Observability Scope (kube-prometheus-stack)

What we're scraping and why:

| Metric Source | Tool | Value |
|---|---|---|
| Node metrics | node-exporter | CPU/memory/disk on system node — spot eviction triggers |
| K8s object metrics | kube-state-metrics | Pod restarts, deployment health, PVC usage |
| Traefik metrics | Traefik built-in `/metrics` | Request rates, error rates, latency per HTTPRoute |
| Jenkins metrics | Jenkins Prometheus plugin | Build queue depth, build duration, agent utilization |
| ArgoCD metrics | ArgoCD built-in metrics | Sync failures, app health, reconciliation latency |
| cert-manager metrics | cert-manager built-in | Cert expiry, ACME challenge success rate |
| Azure cost (alert) | Azure Cost Management | Budget threshold alerts via email — not in-cluster |

---
*Research Date: July 2026*
