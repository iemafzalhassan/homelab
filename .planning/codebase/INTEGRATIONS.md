# External Integrations

> Mapped: 2026-07-20

## Cloud Provider Services

| Service | Purpose | Configuration Location |
|---------|---------|------------------------|
| **Azure Resource Manager** | Core infrastructure provisioning (RG, VNet, AKS, ACR, Key Vault, etc.) | `main.tf`, `modules/**/main.tf` |
| **Azure Kubernetes Service (AKS)** | Managed Kubernetes cluster (v1.36, Azure CNI Overlay, OIDC issuer, Workload Identity) | `modules/aks/main.tf` |
| **Azure Container Registry (ACR)** | Private container image registry (Basic SKU, admin disabled) | `modules/acr/main.tf` |
| **Azure Key Vault** | Secrets store (Standard SKU, RBAC, Private Endpoint, purge protection off) | `modules/keyvault/main.tf` |
| **Azure Virtual Network** | Network isolation (10.0.0.0/16 with system/user/infra/ingress subnets) | `modules/networking/main.tf` |
| **Azure Network Security Groups** | Subnet-level firewall rules (allow VNet, allow 80/443 in/out, NTP, DNS) | `modules/networking/main.tf:37-334` |
| **Azure Private Endpoints** | Private connectivity to Key Vault (and future PaaS) | `modules/keyvault/main.tf:38-55` |
| **Azure Private DNS Zones** | Private DNS resolution for Key Vault (`privatelink.vaultcore.azure.net`) | `modules/keyvault/main.tf:58-69` |
| **Azure AD Workload Identity** | OIDC federation for Kubernetes service accounts (ArgoCD, Jenkins, Traefik, Cert-Manager, Monitoring) | `modules/identity/main.tf`, `main.tf:71-98` |
| **Azure Role Assignments** | RBAC for Key Vault (Secrets User), ACR (AcrPush for Jenkins) | `main.tf:110-128` |
| **Azure Cost Management** | Monthly budget ($25) with 80%/100% alerts | `main.tf:34-58` |
| **Azure Monitor / Log Analytics** | (Referenced via Grafana Cloud integration) | `manifests/bootstrap/k8s-monitoring/values.yaml` |

## Databases

| Database | Purpose | Connection Config |
|----------|---------|-------------------|
| **PostgreSQL (CloudNativePG)** | Keycloak backing database (managed by CNPG operator) | `manifests/bootstrap/keycloak/values.yaml:22-30` (`externalDatabase` config), `manifests/bootstrap/keycloak/cnpg-cluster.yaml` |
| **Keycloak Internal DB** | Not used (external PostgreSQL via CNPG) | `manifests/bootstrap/keycloak/values.yaml:22` (`postgresql.enabled: false`) |

## Message Queues & Event Streaming

| System | Purpose | Config Location |
|--------|---------|-----------------|
| *None detected* | No message queue or event streaming system currently provisioned | — |

## Authentication & Authorization

| Provider | Purpose | Config Location |
|----------|---------|-----------------|
| **Keycloak (OIDC/SAML)** | Central identity provider (SSO, realm: Platform, clients: argocd, jenkins, kargo, oauth2-proxy) | `manifests/bootstrap/keycloak/values.yaml`, `manifests/apps/keycloak.yaml` |
| **ArgoCD OIDC** | ArgoCD authentication via Keycloak (client: argocd, groups: homelab-admins, devops, developers, viewers) | `manifests/bootstrap/argocd/values.yaml:14-19` |
| **Jenkins OIDC** | Jenkins authentication via Keycloak (client: jenkins, matrix auth with group mapping) | `manifests/bootstrap/jenkins/values.yaml:49-80` |
| **Kargo OIDC** | Kargo authentication via Keycloak (client: kargo, admin/viewer group mapping) | `manifests/bootstrap/kargo/values.yaml:11-21` |
| **OAuth2-Proxy** | Reverse proxy authentication for protected apps (provider: keycloak-oidc, client: oauth2-proxy) | `manifests/bootstrap/oauth2-proxy/values.yaml:2-17` |
| **Traefik Gateway API** | TLS termination and routing (letsencrypt-prod ClusterIssuer via cert-manager) | `manifests/bootstrap/traefik/values.yaml:12-71` |
| **Azure AD Workload Identity** | Keyless auth for Kubernetes workloads to Azure (Key Vault CSI, ACR, etc.) | `modules/identity/main.tf`, various `values.yaml` with `azure.workload.identity/*` annotations |

## Monitoring & Observability

| Tool | Purpose | Config Location |
|------|---------|-----------------|
| **Grafana Cloud (Hosted)** | Metrics (Prometheus), Logs (Loki), Traces (OTLP), Profiles (Pyroscope) | `manifests/bootstrap/k8s-monitoring/values.yaml:4-47` |
| **Grafana Alloy** | Telemetry collector (metrics, logs, traces, profiles) replacing Prometheus Agent | `manifests/bootstrap/k8s-monitoring/values.yaml:118-161` |
| **Kube State Metrics** | Kubernetes object metrics | `manifests/bootstrap/k8s-monitoring/values.yaml:163-188` |
| **Node Exporter** | Host-level metrics | `manifests/bootstrap/k8s-monitoring/values.yaml:189-195` |
| **OpenCost** | Kubernetes cost monitoring | `manifests/bootstrap/k8s-monitoring/values.yaml:198-213` |
| **Kepler** | Energy/power metrics | `manifests/bootstrap/k8s-monitoring/values.yaml:216-222` |
| **Beyla (eBPF)** | Auto-instrumentation for application observability | `manifests/bootstrap/k8s-monitoring/values.yaml:297-308` |
| **Cert-Manager** | TLS certificate management (Let's Encrypt via Gateway API) | `manifests/bootstrap/cert-manager/values.yaml`, `manifests/apps/cert-manager.yaml` |

## CI/CD & GitOps

| Platform | Purpose | Config Location |
|----------|---------|-----------------|
| **ArgoCD** | GitOps continuous deployment (Application, ApplicationSet, AppProject) | `manifests/apps/argocd.yaml`, `manifests/bootstrap/argocd/values.yaml` |
| **Kargo** | Progressive delivery / promotion (Warehouse, Stage, Promotion, Project) | `manifests/apps/kargo.yaml`, `kargo/*.yaml`, `manifests/bootstrap/kargo/values.yaml` |
| **Jenkins** | CI pipeline execution (Kaniko for builds, Job DSL for config-as-code) | `manifests/apps/jenkins.yaml`, `manifests/bootstrap/jenkins/values.yaml`, `Utils/Jenkinsfile` |
| **GitHub** | Source control + webhook trigger for Kargo/ArgoCD | `kargo/warehouse.yaml:9`, `manifests/bootstrap/jenkins/values.yaml:43-48` |
| **Kaniko** | Container image builder (in Jenkins agent pod, pushes to ACR) | `manifests/bootstrap/jenkins/values.yaml:142-156`, `Utils/Jenkinsfile` |

## External APIs & Webhooks

| API | Purpose | Config Location |
|-----|---------|-----------------|
| **Cloudflare API** | DNS management (API token stored in Key Vault) | `modules/keyvault/main.tf:79-88` (secret: `cloudflare-api-token`) |
| **Let's Encrypt (ACME)** | TLS certificate issuance via cert-manager (HTTP-01 via Gateway API) | `manifests/bootstrap/traefik/values.yaml:15` (`letsencrypt-prod` ClusterIssuer ref) |
| **Grafana Cloud APIs** | Remote write for metrics/logs/traces/profiles (basic auth credentials in Key Vault) | `manifests/bootstrap/k8s-monitoring/values.yaml:5-47`, `manifests/bootstrap/k8s-monitoring/values.yaml:241-295` (SecretProviderClass) |

## Third-Party Services

| Service | Purpose | Config Location |
|---------|---------|-----------------|
| **Cloudflare** | Public DNS zone (`iemafzalhassan.tech`), DNS records for ingress hostnames | `modules/keyvault/main.tf:81` (token), Traefik listeners reference hostnames |
| **GitHub** | Source repository (`https://github.com/iemafzalhassan/homelab.git`), GitHub Token for Jenkins | `kargo/warehouse.yaml:9`, `manifests/bootstrap/jenkins/values.yaml:14-18`, `manifests/bootstrap/jenkins/values.yaml:86-99` |
| **Grafana Cloud** | Hosted observability stack (Prometheus, Loki, Tempo, Pyroscope) | `manifests/bootstrap/k8s-monitoring/values.yaml:4-47` |
| **Bitnami/Charts** | Helm charts for Keycloak, Jenkins, Traefik, CNPG, ArgoCD, Kargo, OAuth2-Proxy | `manifests/apps/*.yaml` (repoURL references) |

## Network & Security Integrations

| Integration | Purpose | Config Location |
|-------------|---------|-----------------|
| **Gateway API (Kubernetes)** | Ingress routing via Traefik (HTTPRoutes, GRPCRoutes, TLS) | `manifests/bootstrap/traefik/values.yaml:9-10`, `manifests/bootstrap/gateway-api/install.sh` |
| **CSI Secrets Store Driver** | Mount Key Vault secrets as Kubernetes volumes | `manifests/bootstrap/csi-driver/values-csi-driver.yaml`, `manifests/bootstrap/csi-driver/values-akv-provider.yaml` |
| **Azure Key Vault Provider for CSI** | Sync Key Vault secrets to K8s secrets | `manifests/bootstrap/csi-driver/values-akv-provider.yaml` |
| **Network Policies (NSG + CNI)** | Pod/Subnet level network segmentation | `modules/networking/main.tf:37-334`, AKS Azure CNI Overlay |

## Secrets & Credentials Management

| Secret | Stored In | Consumed By |
|--------|-----------|-------------|
| Cloudflare API Token | Azure Key Vault (`cloudflare-api-token`) | External DNS / Cert-Manager (future) |
| Grafana Cloud Token | Azure Key Vault (`grafana-cloud-token`) | Grafana Alloy (via SecretProviderClass) |
| Grafana Metrics Username | Azure Key Vault (`grafana-metrics-username`) | Grafana Alloy |
| Grafana Logs Username | Azure Key Vault (`grafana-logs-username`) | Grafana Alloy |
| Grafana OTLP Username | Azure Key Vault (`grafana-otlp-username`) | Grafana Alloy |
| ArgoCD OIDC Client Secret | Kubernetes Secret (`argocd-oidc-secret`) | ArgoCD |
| Jenkins OIDC Client Secret | Kubernetes Secret (`jenkins-oidc-secret`) | Jenkins |
| Jenkins GitHub Token | Kubernetes Secret (`jenkins-github-token`) | Jenkins |
| Kargo Admin Password Hash | Helm values (`manifests/bootstrap/kargo/values.yaml:3`) | Kargo API |
| Kargo Token Signing Key | Helm values (`manifests/bootstrap/kargo/values.yaml:4`) | Kargo API |

---

*Integration audit: 2026-07-20*