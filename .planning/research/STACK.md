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
- **Traefik v3 + Gateway API over ingress-nginx + IngressRoute**: `kubernetes/ingress-nginx` EOL March 2026. Traefik v3 with Kubernetes **Gateway API** (`HTTPRoute`) is the correct modern choice — GA (v1.5 Feb 2026), 100% Traefik conformance, no vendor CRD lock-in, cert-manager native support

## ⛔ ingress-nginx EOL — Decision Log

**Status:** `kubernetes/ingress-nginx` GitHub repository archived and read-only as of **March 2026**.

| What is EOL | What is NOT EOL |
|---|---|
| `kubernetes/ingress-nginx` (community) | Kubernetes Ingress API itself |
| — No security patches | `nginxinc/kubernetes-ingress` (F5, still active) |
| — No K8s version compatibility updates | Traefik, Kong, Envoy Gateway, etc. |

**Chosen routing API: Kubernetes Gateway API** (not legacy Ingress or Traefik IngressRoute CRDs)

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
```bash
# Gateway API v1.5.1 — Standard channel (stable resources only)
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.5.1/standard-install.yaml
```

**Resources installed:** `GatewayClass`, `Gateway`, `HTTPRoute`, `GRPCRoute`, `ReferenceGrant`

### Traefik v3 Helm values (Gateway API first-class)
```yaml
providers:
  kubernetesIngress:
    enabled: false     # Disable legacy Ingress — we use Gateway API exclusively
  kubernetesGateway:
    enabled: true      # Gateway API is PRIMARY routing API
  kubernetesCRD:
    enabled: false     # No IngressRoute CRDs needed

gateway:
  enabled: true        # Auto-create GatewayClass + default Gateway
  listeners:
    web:
      port: 8000
      protocol: HTTP
    websecure:
      port: 8443
      protocol: HTTPS

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
  enabled: true        # Secure via HTTPRoute + auth middleware
```

### cert-manager integration with Gateway API
```yaml
# cert-manager >= v1.15 with Gateway API enabled
# helm upgrade cert-manager ... --set config.enableGatewayAPI=true

# Annotate the Gateway to auto-trigger certificate issuance:
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: homelab-gateway
  namespace: traefik
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod   # ← triggers auto-issuance
spec:
  gatewayClassName: traefik
  listeners:
  - name: https
    port: 443
    protocol: HTTPS
    hostname: "*.yourdomain.com"
    tls:
      mode: Terminate
      certificateRefs:
      - name: wildcard-tls-secret   # cert-manager creates this automatically
```

### Application routing pattern (HTTPRoute)
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: argocd
  namespace: argocd
spec:
  parentRefs:
  - name: homelab-gateway
    namespace: traefik     # cross-namespace: requires ReferenceGrant
  hostnames:
  - "argocd.yourdomain.com"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: argocd-server
      port: 80
```

### Cross-namespace ReferenceGrant (required for Gateway API)
```yaml
# In each app namespace: allows traefik Gateway to bind HTTPRoutes
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: allow-traefik-gateway
  namespace: argocd   # repeat per namespace (argocd, jenkins, monitoring)
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: Gateway
    namespace: traefik
  to:
  - group: ""
    kind: Service
```

---
*Confidence: HIGH — Gateway API v1.5 GA, Traefik v3 conformance verified, cert-manager Gateway API docs confirmed July 2026*

