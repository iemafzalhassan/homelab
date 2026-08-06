# Kube Telemetry Stage 🚀

A production-patterned Kubernetes observability and Grafana demo environment for CNCF virtual speaking events. The platform runs ArgoCD for GitOps-based deployments, Jenkins with ephemeral JNLP build agents, and a full observability stack — all secured via Workload Identity Federation, private networking, and Azure Key Vault for secrets.

## ✨ Core Value

A real, running Kubernetes platform that teaches production-grade patterns (GitOps, zero-trust secrets, HA networking) without breaking the bank — the cluster stays alive 24x7 within a strict **$100 total budget**.

## 🏗️ Architecture & Technology Stack

- **Infrastructure as Code**: Terraform (Azure Provider `~> 4.0`)
- **Cloud Provider**: Microsoft Azure (Central India Region)
- **Kubernetes**: Azure Kubernetes Service (AKS) - Free Tier Control Plane
  - System Node Pool: 1 reserved `Standard_D2as_v5` node
  - User Node Pool: Spot `Standard_D2as_v5` nodes (scale-to-zero)
- **Networking**: AKS CNI Overlay, fully private worker nodes, Traefik v3 via Kubernetes Gateway API (GA v1.5)
- **Security & Identity**: 
  - Azure Key Vault with Private Endpoints
  - Workload Identity Federation (Zero static credentials)
  - Secrets Store CSI Driver
- **GitOps & CI/CD**:
  - **ArgoCD**: Single deployment operator for the platform.
  - **Jenkins**: JNLP ephemeral build agents.
- **Observability**: `kube-prometheus-stack` (Prometheus, Grafana, Alertmanager)

## 📁 Repository Structure

```text
├── modules/                # Terraform modules
│   ├── acr/                # Azure Container Registry
│   ├── aks/                # Azure Kubernetes Service
│   ├── identity/           # Workload Identity Federation setup
│   ├── keyvault/           # Azure Key Vault and Secrets
│   └── networking/         # VNet, Subnets, and Network Security
├── manifests/              # Kubernetes Manifests (Helm / GitOps)
│   ├── apps/               # Application deployments (e.g., opentelemetry-demo)
│   ├── bootstrap/          # Core platform services (ArgoCD, Jenkins, Traefik)
│   └── portfolio/          # Portfolio related apps
├── scripts/                # Helper utilities and automation
├── main.tf                 # Primary Terraform entrypoint
├── variables.tf            # Input variables
├── outputs.tf              # Terraform outputs
└── terraform.tfvars        # (Gitignored) Environment-specific variables
```

## 🚀 Getting Started

### Prerequisites

1. [Terraform](https://www.terraform.io/downloads.html) (>= 1.5.0)
2. [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) (`az`)
3. `kubectl` and `helm`

### Deployment

1. **Login to Azure:**
   ```bash
   az login
   az account set --subscription "<YOUR_SUBSCRIPTION_ID>"
   ```

2. **Initialize Terraform:**
   ```bash
   terraform init
   ```

3. **Configure Variables:**
   Create a `terraform.tfvars` file (this is ignored by Git) and populate it with required values:
   ```hcl
   location             = "centralindia"
   rg_name              = "homelab-rg"
   budget_contact_email = "your.email@example.com"
   admin_ip_ranges      = ["YOUR.PUBLIC.IP.ADDRESS/32"]
   ```

4. **Plan and Apply:**
   ```bash
   terraform plan -out=tfplan
   terraform apply tfplan
   ```

5. **Connect to AKS:**
   ```bash
   az aks get-credentials --resource-group homelab-rg --name <CLUSTER_NAME>
   ```

## 🔒 Security Posture

- **Zero Trust Identity**: No Service Principal secrets are used. All K8s to Azure authentication goes through Workload Identity Federation.
- **Private by Default**: Worker nodes lack public IPs. The Key Vault is only accessible internally via Private Endpoints. The API server is secured by an IP whitelist.

## 🐛 Debugging & Runbooks

Active debug sessions and historical troubleshooting context are kept in `.planning/debug/`. The cluster utilizes GSD workflows for stateful session persistence.

---
*Created for presentation and production-grade demonstration purposes.*
