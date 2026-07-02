# Stack Research — AKS Homelab Platform (2026)

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

```
homelab/
├── README.md
├── main.tf           # Root: calls all modules, wires outputs
├── variables.tf      # Input vars: subscription_id, location, tags, etc.
├── outputs.tf        # Cluster endpoint, OIDC issuer URL, etc.
├── versions.tf       # terraform >= 1.9, azurerm ~> 4.79
├── terraform.tfvars  # Actual values (gitignored for secrets)
│
├── modules/
│   ├── networking/   # VNet, subnets, NSGs
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── aks/          # AKS cluster, node pools, OIDC, workload identity
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── keyvault/     # Key Vault, private endpoint, RBAC
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── identity/     # User-assigned managed identities, federated credentials
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
│
└── gitops/           # Helm values files for ArgoCD ApplicationSets
    ├── argocd/
    ├── jenkins/
    ├── traefik/
    ├── cert-manager/
    └── monitoring/
```

## Key Terraform Snippets

### AKS with CNI Overlay + Workload Identity
```hcl
resource "azurerm_kubernetes_cluster" "aks" {
  oidc_issuer_enabled       = true
  workload_identity_enabled = true
  
  network_profile {
    network_plugin    = "azure"
    network_plugin_mode = "overlay"
    pod_cidr          = "10.244.0.0/16"
  }
  
  default_node_pool {
    vm_size             = "Standard_D2as_v5"
    zones               = ["1"]
    node_count          = 1
    os_disk_type        = "Managed"
    os_disk_size_gb     = 30
  }
}
```

### Federated Identity Credential
```hcl
resource "azurerm_federated_identity_credential" "argocd" {
  name                = "argocd-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.argocd.id
  issuer              = azurerm_kubernetes_cluster.aks.oidc_issuer_url
  subject             = "system:serviceaccount:argocd:argocd-server"
}
```

## Rationale

- **D2as_v5 not B-series**: Microsoft explicitly advises against B-series for AKS system node pools — CPU credit exhaustion causes CoreDNS and kubelet failures
- **CNI Overlay**: Pod IPs don't consume VNet address space; supports large pod counts on /24 subnets
- **Secrets Store CSI v1.6.0**: Requires K8s 1.30+; pairs with AKV provider for zero-credential secret access
- **ArgoCD 10.x (v3.x)**: Major version with improved ApplicationSet support and multi-source Apps
- **Traefik v3 not ingress-nginx**: `kubernetes/ingress-nginx` reached EOL March 2026 — archived, no patches. Traefik v3 is the CNCF-backed replacement with dual support for legacy `IngressRoute` CRDs and modern `Gateway API` (`HTTPRoute`)

## ⛔ ingress-nginx EOL — Decision Log

**Status:** `kubernetes/ingress-nginx` GitHub repository archived and read-only as of **March 2026**.

| What is EOL | What is NOT EOL |
|---|---|
| `kubernetes/ingress-nginx` (community) | Kubernetes Ingress API itself |
| — No security patches | `nginxinc/kubernetes-ingress` (F5, still active) |
| — No K8s version compatibility updates | Traefik, Kong, Envoy Gateway, etc. |

**Chosen replacement: Traefik v3** (`traefik/traefik` Helm chart `41.0.1`)

**Why Traefik over F5 NGINX Ingress:**
- CNCF project, fully open-source, no commercial strings
- Native **Gateway API** support (`HTTPRoute`, `GatewayClass`) — future-proof
- Legacy `IngressRoute` CRD still supported for gradual migration
- Single `LoadBalancer` service = same single public IP architecture as before
- Built-in dashboard (useful for homelab visibility)
- Clean cert-manager + Cloudflare DNS-01 integration

**Traefik v3 Helm values for AKS (key settings):**
```yaml
providers:
  kubernetesIngress:
    enabled: true      # Backwards compat with Ingress resources
  kubernetesGateway:
    enabled: true      # Forward-compat with Gateway API HTTPRoutes

service:
  type: LoadBalancer
  annotations:
    service.beta.kubernetes.io/azure-load-balancer-resource-group: "homelab-rg"

entryPoints:
  web:
    address: ":80"
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https
  websecure:
    address: ":443"

dashboard:
  enabled: true   # Expose via IngressRoute, auth-protected
```

**cert-manager integration pattern with Traefik:**
- cert-manager issues `Certificate` resources (Cloudflare DNS-01)
- TLS secret referenced explicitly in `IngressRoute.spec.tls.secretName`
- (Unlike ingress-nginx, Traefik does NOT auto-watch cert-manager annotations — explicit Certificate object required)

---
*Confidence: HIGH — versions verified against ArtifactHub, Traefik GitHub, and official docs July 2026*
