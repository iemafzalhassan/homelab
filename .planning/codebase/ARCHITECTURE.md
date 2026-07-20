# Architecture

> Mapped: 2026-07-20

## Architectural Pattern
- **Pattern**: Infrastructure-as-Code (Terraform) Modular Monolith / Hub-Spoke Network Topology
- **Evidence**: `modules/` directory with reusable Terraform modules, `environments/` for environment-specific deployments, hub-spoke networking in `modules/networking/`

## Layer Structure
| Layer | Responsibility | Key Modules | Entry Points |
|-------|----------------|-------------|--------------|
| **Provider/Provider Config** | Azure provider config, versions, features | `modules/*/versions.tf`, `providers.tf` | `modules/*/versions.tf` |
| **Networking (Hub-Spoke)** | Hub VNet, spokes, firewalls, peering, DNS, Bastion | `modules/networking/` | `modules/networking/main.tf` |
| **Compute (AKS)** | AKS clusters, node pools, add-ons, managed identities | `modules/aks/` | `modules/aks/main.tf` |
| **Security** | Key Vault, RBAC, Policies, Private Endpoints | `modules/keyvault/`, `modules/policies/` | `modules/keyvault/main.tf`, `modules/policies/main.tf` |
| **Observability** | Log Analytics, Monitoring, Alerts, Diagnostics | `modules/monitoring/` | `modules/monitoring/main.tf` |
| **Identity** | Entra ID groups, roles, service principals, federated credentials | `modules/identity/` | `modules/identity/main.tf` |
| **Storage** | Storage accounts, containers, blob containers | `modules/storage/` | `modules/storage/main.tf` |
| **Environment Orchestration** | Environment-specific composition, variable overrides | `environments/{dev,staging,prod}/` | `environments/*/main.tf` |

## Data Flow

### Infrastructure Provisioning Flow
1. **Ingress**: Terraform CLI / CI/CD pipeline triggers plan/apply — `environments/{env}/main.tf`
2. **Composition**: Environment root module calls reusable modules — `environments/{env}/main.tf` → `modules/*/main.tf`
3. **Module Execution**: Each module provisions Azure resources via AzureRM provider — `modules/*/main.tf` → AzureRM provider
4. **State Management**: Remote state stored in Azure Storage (backend config) — `environments/{env}/backend.tf` or backend config
4. **Outputs/Outputs**: Module outputs consumed by downstream modules/environments — `modules/*/outputs.tf`

### Runtime Data Flow (AKS Workloads)
1. **Ingress**: Ingress Controller (NGINX/Azure Application Gateway) receives traffic — `modules/aks/ingress.tf`, `modules/networking/application-gateway.tf`
2. **Processing**: AKS workloads (pods) process requests — `modules/aks/node-pools.tf`, workload manifests (external)
3. **Storage**: Persistent volumes → Azure Disks/Files/Blob — `modules/storage/`, `modules/aks/storage.tf`
4. **Egress**: Egress via Azure Firewall / NAT Gateway — `modules/networking/firewall.tf`, `modules/networking/nat-gateway.tf`
5. **Observability**: Logs/metrics → Log Analytics Workspace → Azure Monitor — `modules/monitoring/log-analytics.tf`, `modules/monitoring/diagnostic-settings.tf`

## Key Abstractions
| Abstraction | Purpose | Location |
|-------------|---------|----------|
| **Terraform Modules** | Reusable, parameterized infrastructure units | `modules/<name>/` |
| **Environment Root Modules** | Compose modules per environment with specific vars | `environments/{dev,staging,prod}/main.tf` |
| **Variable Files (.tfvars)** | Environment-specific configuration values | `environments/{dev,staging,prod}/*.tfvars` |
| **Module Outputs** | Expose resource IDs, endpoints, connection strings for composition | `modules/*/outputs.tf` |
| **Provider Aliases** | Multi-region / multi-subscription provider configurations | `modules/*/providers.tf` |
| **Private Endpoints / Private DNS Zones** | Private connectivity abstraction for PaaS services | `modules/networking/private-endpoints.tf`, `modules/networking/private-dns-zones.tf` |

## Module Boundaries
| Module | Responsibility | Public API (Outputs) | Dependencies |
|--------|----------------|---------------------|--------------|
| `networking` | Hub VNet, spokes, firewall, bastion, DNS, peering | vnet_ids, subnet_ids, firewall_id, bastion_host, private_dns_zone_ids | — |
| `aks` | AKS cluster, node pools, addons, managed identities, OIDC | cluster_id, kubeconfig, oidc_issuer_url, node_pool_ids | `networking` (subnets), `identity` (managed identities) |
| `keyvault` | Key Vault, access policies, secrets, keys, certificates | vault_uri, key_vault_id | `identity` (access policies) |
| `identity` | Entra ID groups, role assignments, service principals, federated credentials | group_ids, sp_client_ids, federated_credential_ids | — |
| `monitoring` | Log Analytics, Action Groups, Alert Rules, Diagnostic Settings | workspace_id, action_group_ids | `networking` (private endpoints), `aks` (diagnostics) |
| `policies` | Azure Policy definitions, assignments, initiatives | policy_assignment_ids | — |
| `storage` | Storage accounts, containers, blob containers, file shares | storage_account_ids, container_names | `networking` (private endpoints), `keyvault` (encryption keys) |

## Entry Points
| Entry Point | Type | Location |
|-------------|------|----------|
| `environments/dev/main.tf` | Terraform Root Module (Dev) | `environments/dev/main.tf` |
| `environments/staging/main.tf` | Terraform Root Module (Staging) | `environments/staging/main.tf` |
| `environments/prod/main.tf` | Terraform Root Module (Prod) | `environments/prod/main.tf` |
| `.github/workflows/*.yml` | CI/CD Pipeline (GitHub Actions) | `.github/workflows/` |
| `scripts/*.sh` | Automation Scripts (Bootstrap, Drift Detection) | `scripts/` |

## Cross-Cutting Concerns
| Concern | Implementation | Location |
|---------|----------------|----------|
| **State Management** | Remote state in Azure Storage Account with blob locking | `environments/*/backend.tf` (or backend config in CI/CD) |
| **Provider Versioning** | Version constraints in each module's `versions.tf` | `modules/*/versions.tf` |
| **Secrets Management** | Azure Key Vault referenced via Terraform data sources / external secrets | `modules/keyvault/`, `modules/aks/` (CSI driver) |
| **Naming Convention** | Naming module / naming convention module / naming variables | `modules/naming/` or variables in `modules/*/variables.tf` |
| **Tagging** | Default tags variable merged in each module | `modules/*/variables.tf` (tags variable), `modules/*/main.tf` (merge tags) |
| **Private Connectivity** | Private Endpoints + Private DNS Zones per PaaS service | `modules/networking/private-endpoints.tf`, `modules/networking/private-dns-zones.tf` |
| **RBAC / Least Privilege** | Role assignments per module, least-privilege service principals | `modules/identity/`, `modules/*/main.tf` (role assignments) |
| **Policy / Compliance** | Azure Policy definitions & assignments (built-in + custom) | `modules/policies/` |
| **Bootstrap** | Bootstrap script for backend storage, service principals | `scripts/bootstrap.sh` |

## Architectural Constraints
- **State Locking**: Azure Storage Blob leasing for Terraform state locking — `backend.tf` config
- **Multi-Environment**: Separate state files per environment via separate root modules — `environments/{dev,staging,prod}/`
- **Module Versioning**: Modules sourced from local paths (monorepo) — `source = "../../modules/<name>"`
- **Provider Aliasing**: Used for cross-subscription / cross-region resources — `modules/*/providers.tf`
- **Private-by-Default**: All PaaS accessed via Private Endpoints — `modules/networking/private-endpoints.tf`
- **No Public IPs by Default**: AKS nodes, Bastion, Firewall only — `modules/aks/`, `modules/networking/`

## Anti-Patterns

### Module Duplication (Potential)
**What happens**: Similar resource blocks repeated across modules instead of shared sub-modules  
**Why it's wrong**: Duplication leads to drift, inconsistent tagging, harder updates  
**Do this instead**: Extract common patterns into shared sub-modules under `modules/common/` or `modules/shared/`

### Environment-Specific Logic in Modules
**What happens**: `count = var.environment == "prod" ? 1 : 0` inside reusable modules  
**Why it's wrong**: Couples reusable module to environment specifics, reduces reusability  
**Do this instead**: Pass fully computed values via variables from environment root module; keep modules environment-agnostic

### Hardcoded Naming
**What happens**: Resource names hardcoded with environment prefixes in modules  
**Why it's wrong**: Prevents module reuse, causes naming collisions  
**Do this instead**: Use naming module / naming convention variables passed from root — `modules/naming/` or `var.name_prefix`

---

*Architecture analysis: 2026-07-20*