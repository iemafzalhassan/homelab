# Architecture Research — AKS Homelab Platform (2026)

## Component Map

```
┌─────────────────────────────────────────────────────────────────────────┐
│  AZURE CLOUD (Central India)                                            │
│                                                                         │
│  Resource Group: homelab-rg                                             │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │  VNet: 10.0.0.0/16                                                │ │
│  │                                                                   │ │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────────┐ │ │
│  │  │ system-subnet   │  │ user-subnet      │  │ infra-subnet     │ │ │
│  │  │ 10.0.1.0/24     │  │ 10.0.2.0/24      │  │ 10.0.4.0/24      │ │ │
│  │  │ Zone 1          │  │ Zone 1+2 (spot)  │  │ AKV Private EP   │ │ │
│  │  └────────┬────────┘  └────────┬─────────┘  └──────────────────┘ │ │
│  │           │                    │                                   │ │
│  │  ┌────────▼────────────────────▼─────────────────────────────┐   │ │
│  │  │  AKS CLUSTER (Free tier, API auth via IP allowlist)        │   │ │
│  │  │  K8s v1.35, CNI Overlay, OIDC issuer enabled               │   │ │
│  │  │                                                             │   │ │
│  │  │  ┌─────────────┐    ┌──────────────────────────────────┐   │   │ │
│  │  │  │ system pool │    │ spot user pool (scale 0→4)        │   │   │ │
│  │  │  │ D2as_v5 ×1  │    │ D2as_v5 spot, Zone 1+2           │   │   │ │
│  │  │  │ Zone 1      │    │ Jenkins JNLP pods only            │   │   │ │
│  │  │  └─────────────┘    └──────────────────────────────────┘   │   │ │
│  │  │                                                             │   │ │
│  │  │  NAMESPACES:                                                │   │ │
│  │  │  traefik/     → GatewayClass + Gateway + LoadBalancer svc  │   │ │
│  │  │  cert-manager/→ ClusterIssuer (letsencrypt, Cloudflare)    │   │ │
│  │  │  argocd/      → App-of-Apps root, manages everything below │   │ │
│  │  │  jenkins/     → Controller + ephemeral JNLP agent pods     │   │ │
│  │  │  monitoring/  → Prometheus, Grafana, Alertmanager          │   │ │
│  │  └─────────────────────────────────────────────────────────────┘   │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│  Azure Key Vault (Private Endpoint → infra-subnet)                     │
│  User-Assigned Managed Identities (per workload):                      │
│    • argocd-identity    → AKV: Key Vault Secrets User                  │
│    • jenkins-identity   → AKV: Key Vault Secrets User, ACR: AcrPull   │
│    • traefik-identity   → (future: Azure DNS for DNS-01)               │
└─────────────────────────────────────────────────────────────────────────┘
                              │ Public IP (single)
                   ┌──────────▼──────────┐
                   │   Cloudflare DNS    │
                   │   (Free proxy tier) │
                   │   *.yourdomain.com  │
                   └─────────────────────┘
```

## Data Flow: External Request → Service

```
User → Cloudflare (DNS proxy, DDoS, SSL edge)
     → Azure Public IP (Traefik LoadBalancer svc, port 443)
     → Traefik v3 (GatewayClass implementation, TLS terminates here)
     → HTTPRoute (namespace-scoped, e.g.: argocd ns)
     → Service ClusterIP
     → Pod
```

## Data Flow: CI/CD (Jenkins → ArgoCD)

```
Developer git push
  → GitHub repo (source code)
  → Jenkins webhook trigger
  → Jenkins JNLP pod spins up on spot node
     (Workload Identity → ACR pull, AKV secrets)
  → Build → Test → Docker build → Push to ACR
  → Jenkins updates image tag in GitOps repo
  → ArgoCD detects diff in GitOps repo
  → ArgoCD applies updated manifest to cluster
  → Pod rolls out with new image
```

## Data Flow: Secrets (Zero Static Credentials)

```
Pod startup:
  Kubernetes projects bound ServiceAccount token → /var/run/secrets/...
  Azure AD validates token against OIDC issuer URL
  Azure issues short-lived access token for Managed Identity
  Secrets Store CSI Driver fetches secret from AKV
  Secret mounted as file into Pod (never stored in etcd)
```

## Terraform Module Boundaries

```
modules/
├── networking/          MANAGES: VNet, subnets, NSGs, public IP
│   OUTPUTS: subnet IDs, public_ip_id
│   DEPENDS ON: nothing
│
├── aks/                 MANAGES: AKS cluster, system pool, spot pool
│   OUTPUTS: cluster_id, oidc_issuer_url, kubelet_identity_id
│   DEPENDS ON: networking (subnet IDs)
│
├── keyvault/            MANAGES: Key Vault, private endpoint, RBAC assignments
│   OUTPUTS: vault_uri, vault_id
│   DEPENDS ON: networking (infra-subnet), identity (managed_identity IDs)
│
└── identity/            MANAGES: User-assigned managed identities, federated credentials
    OUTPUTS: client_ids, principal_ids per workload
    DEPENDS ON: aks (oidc_issuer_url)
```

## GitOps Layer: ArgoCD App-of-Apps Pattern

```
argocd/
└── Application: root-app          ← bootstrap this one manually (helm install or kubectl apply)
    ├── Application: traefik        ← Helm chart + values, Gateway API CRDs pre-applied
    ├── Application: cert-manager   ← Helm chart + ClusterIssuer + Cloudflare secret
    ├── Application: jenkins        ← Helm chart + JNLP config + workload identity SA
    ├── Application: monitoring     ← kube-prometheus-stack + dashboards
    └── Application: platform-rbac  ← NetworkPolicies, ReferenceGrants, LimitRanges
```

## Key Interface: Workload Identity Binding

Each workload follows this pattern:
1. Terraform creates `azurerm_user_assigned_identity`
2. Terraform creates `azurerm_federated_identity_credential` (links to K8s ServiceAccount)
3. Terraform assigns AKV RBAC role to the identity
4. ArgoCD deploys ServiceAccount with annotation `azure.workload.identity/client-id: <client_id>`
5. Pod gets label `azure.workload.identity/use: "true"`
6. Secrets Store CSI `SecretProviderClass` references the vault URI + identity

## Suggested Build Order (Dependency-Driven)

| Phase | What | Why First |
|---|---|---|
| 1 | Resource Group, VNet, subnets, NSGs | Foundation for everything |
| 2 | AKS cluster (system pool only) | Needs subnets, outputs OIDC URL |
| 3 | Managed Identities + federated credentials | Needs OIDC issuer URL from AKS |
| 4 | Key Vault + Private Endpoint + RBAC | Needs identity client IDs |
| 5 | Gateway API CRDs (kubectl apply) | Must precede Traefik |
| 6 | Traefik v3 (Helm) | Needs Gateway API CRDs |
| 7 | cert-manager (Helm) + ClusterIssuer | Needs Traefik LB IP for Cloudflare DNS |
| 8 | ArgoCD (Helm bootstrap) | Takes over all subsequent deployments |
| 9 | ArgoCD App-of-Apps (Jenkins, monitoring, rbac) | Managed by ArgoCD from here |
| 10 | Spot node pool (via Terraform) | Add after system is stable |

---
*Research Date: July 2026*
