# Technology Stack

> Mapped: 2026-07-20

## Languages & Runtimes
- **Primary**: HCL (Terraform) — Infrastructure as Code (`main.tf`, `modules/**/*.tf`)
- **Secondary**: YAML — Kubernetes manifests and Helm values (`manifests/**/*.yaml`)
- **Runtime**: Kubernetes 1.36 (AKS) — Container orchestration (`modules/aks/main.tf`)
- **Container Runtime**: containerd (AKS default) — Container runtime

## Frameworks & Libraries

| Category | Framework/Library | Version | Evidence |
|----------|-------------------|---------|----------|
| Infrastructure | Terraform | ~> 4.0 (azurerm provider) | `main.tf:3-6` |
| Infrastructure | Terraform Azure Provider | ~> 4.0 | `main.tf:3-6` |
| Infrastructure | Terraform Random Provider | 3.9.0 | `.terraform.lock.hcl:25-44` |
| Container Orchestration | Azure Kubernetes Service (AKS) | 1.36 | `modules/aks/main.tf:7` |
| GitOps | ArgoCD | Latest (helm chart) | `manifests/apps/argocd.yaml:13` |
| Progressive Delivery | Kargo | 1.10.8 | `manifests/apps/kargo.yaml:13` |
| CI/CD | Jenkins | 5.9.32 (helm chart) | `manifests/apps/jenkins.yaml:13` |
| Ingress/Gateway | Traefik | 41.0.2 (v3.7.6) | `manifests/apps/traefik.yaml:13`, `manifests/bootstrap/traefik/values.yaml:2` |
| Identity Provider | Keycloak | 25.2.0 (Bitnami) | `manifests/apps/keycloak.yaml:13` |
| Auth Proxy | OAuth2-Proxy | 7.8.1 | `manifests/apps/oauth2-proxy.yaml:13` |
| Database Operator | CloudNativePG (CNPG) | 0.29.0 | `manifests/apps/cnpg.yaml:13` |
| TLS Certificates | cert-manager | Latest | `manifests/bootstrap/cert-manager/values.yaml` |
| Monitoring | Grafana Cloud / Grafana Alloy | Latest | `manifests/bootstrap/k8s-monitoring/values.yaml` |
| Secrets Management | Azure Key Vault + CSI Driver | Standard SKU | `modules/keyvault/main.tf:14`, `manifests/bootstrap/csi-driver/` |
| Container Registry | Azure Container Registry (ACR) | Basic SKU | `modules/acr/main.tf:11` |

## Package Managers & Build Tools
- **Manager**: Terraform Registry (providers), Helm (charts), Container registries (images) — `.terraform.lock.hcl`, `manifests/**/*.yaml`
- **Build**: Terraform (plan/apply), Helm (template/install), Docker/Kaniko (container images) — `Utils/Jenkinsfile`, `manifests/bootstrap/jenkins/values.yaml:155`

## Infrastructure & Cloud

| Component | Technology | Evidence |
|-----------|------------|----------|
| **Cloud Provider** | Microsoft Azure (centralindia region) | `variables.tf:4`, `main.tf:23` |
| **IaC** | Terraform (azurerm, random providers) | `main.tf:1-12`, `.terraform.lock.hcl` |
| **Container Orchestration** | AKS (Azure Kubernetes Service) | `modules/aks/main.tf` |
| **Container Registry** | Azure Container Registry (ACR) | `modules/acr/main.tf` |
| **Networking** | Azure VNet, Subnets, NSGs, Private Endpoints | `modules/networking/main.tf` |
| **Identity** | Azure AD Workload Identity (OIDC Federation) | `modules/identity/main.tf` |
| **Secrets** | Azure Key Vault (RBAC, Private Endpoint) | `modules/keyvault/main.tf` |
| **DNS** | Cloudflare (external DNS management) | `modules/keyvault/main.tf:81` (API token secret) |
| **Cost Management** | Azure Consumption Budget | `main.tf:34-58` |

## Configuration Management

| Category | Tool/Method | Evidence |
|----------|-------------|----------|
| **Environment Variables** | Terraform variables (.tfvars) | `terraform.tfvars`, `variables.tf` |
| **Secrets** | Azure Key Vault + CSI Driver | `modules/keyvault/main.tf`, `manifests/bootstrap/csi-driver/values-akv-provider.yaml` |
| **Kubernetes Config** | Helm values.yaml + Kustomize overlays | `manifests/bootstrap/**/values.yaml`, `manifests/apps/homelab-demo/**/kustomization.yaml` |
| **GitOps** | ArgoCD Application resources | `manifests/apps/*.yaml` |
| **Progressive Delivery** | Kargo Stages/Projects | `kargo/*.yaml` |

## Key Versions & Constraints
- **Terraform**: Azure Provider ~> 4.0, Random Provider 3.9.0 (`.terraform.lock.hcl`)
- **Kubernetes**: 1.36 (AKS) (`modules/aks/main.tf:7`)
- **Node Pools**: System (Standard_D2s_v3, 1 node), Spot (Standard_D2as_v5, 1-4 nodes) (`modules/aks/main.tf:20-75`)
- **Network Plugin**: Azure CNI Overlay (`modules/aks/main.tf:13-14`)
- **OIDC Issuer**: Enabled on AKS (`modules/aks/main.tf:9-10`)
- **Workload Identity**: Enabled on AKS (`modules/aks/main.tf:10`)

---

*Stack analysis: 2026-07-20*