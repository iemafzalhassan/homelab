# Repository Structure

> Mapped: 2026-07-20

## Directory Tree (Key Paths)
```
repo-root/
├── .github/
│   └── workflows/              # CI/CD pipelines (GitHub Actions)
├── .planning/                  # GSD planning artifacts
├── docs/                       # Documentation (architecture diagrams, runbooks)
├── modules/                    # Reusable Terraform modules
│   ├── aks/                    # Azure Kubernetes Service clusters
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── node-pools.tf
│   │   ├── addons.tf
│   │   ├── ingress.tf
│   │   ├── storage.tf
│   │   ├── identity.tf
│   │   └── README.md
│   ├── identity/               # Entra ID, Service Principals, Federated Credentials
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── groups.tf
│   │   ├── service-principals.tf
│   │   ├── federated-credentials.tf
│   │   └── README.md
│   ├── keyvault/               # Azure Key Vault, secrets, keys, certificates
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── access-policies.tf
│   │   ├── secrets.tf
│   │   └── README.md
│   ├── monitoring/             # Log Analytics, Monitoring, Alerts, Diagnostics
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── log-analytics.tf
│   │   ├── action-groups.tf
│   │   ├── alert-rules.tf
│   │   ├── diagnostic-settings.tf
│   │   └── README.md
│   ├── networking/             # Hub-Spoke VNet, Firewall, Bastion, DNS, Peering, Private Endpoints
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── hub-vnet.tf
│   │   ├── spoke-vnets.tf
│   │   ├── firewall.tf
│   │   ├── bastion.tf
│   │   ├── private-endpoints.tf
│   │   ├── private-dns-zones.tf
│   │   ├── peering.tf
│   │   ├── nat-gateway.tf
│   │   ├── application-gateway.tf
│   │   └── README.md
│   ├── policies/               # Azure Policy definitions & assignments
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── definitions/
│   │   │   ├── allowed-locations.tf
│   │   │   ├── allowed-resource-types.tf
│   │   │   ├── required-tags.tf
│   │   │   └── ...
│   │   ├── assignments/
│   │   │   ├── platform-policies.tf
│   │   │   └── ...
│   │   └── README.md
│   ├── storage/                # Storage Accounts, Containers, File Shares
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── storage-accounts.tf
│   │   ├── containers.tf
│   │   ├── file-shares.tf
│   │   └── README.md
│   └── naming/                 # (Optional) Naming convention module
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       └── README.md
├── environments/               # Environment-specific root modules
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── backend.tf
│   │   ├── terraform.tfvars
│   │   └── providers.tf
│   ├── staging/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   ├── versions.tf
│   │   ├── backend.tf
│   │   ├── terraform.tfvars
│   │   └── providers.tf
│   └── prod/
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       ├── versions.tf
│       ├── backend.tf
│       ├── terraform.tfvars
│       └── providers.tf
├── scripts/                    # Automation & bootstrap scripts
│   ├── bootstrap.sh            # Bootstrap backend storage, service principals
│   ├── plan.sh                 # Terraform plan wrapper
│   ├── apply.sh                # Terraform apply wrapper
│   ├── drift-detection.sh      # Drift detection script
│   └── validate.sh             # Validation/linting script
├── tests/                      # Terratest / integration tests
│   ├── fixtures/
│   ├── test_aks.go
│   ├── test_networking.go
│   └── go.mod
├── .github/
│   └── workflows/
│       ├── terraform-plan.yml
│       ├── terraform-apply.yml
│       ├── drift-detection.yml
│       └── validate.yml
├── .gitignore
├── .terraform-version          # Terraform version pin (tfenv/asdf)
├── README.md
└── VERSION                     # Repository version tag
```

## Key Directories & Files
| Path | Purpose | Notes |
|------|---------|-------|
| `modules/` | Reusable Terraform modules (library) | Each module is self-contained with `main.tf`, `variables.tf`, `outputs.tf`, `versions.tf`, `README.md` |
| `environments/` | Environment-specific root modules | One directory per environment (`dev`, `staging`, `prod`); each composes modules with env-specific vars |
| `scripts/` | Automation & bootstrap scripts | Shell scripts for CI/CD, bootstrap, drift detection, validation |
| `tests/` | Integration tests (Terratest) | Go-based tests validating module behavior against real Azure |
| `.github/workflows/` | CI/CD pipelines | GitHub Actions for plan, apply, drift detection, validation |
| `docs/` | Architecture diagrams, runbooks, ADRs | Markdown + diagrams (Mermaid/PlantUML) |

## Naming Conventions
| Pattern | Example | Scope |
|---------|---------|-------|
| Module directory | `modules/<name>/` | All modules under `modules/` |
| Module main config | `modules/<name>/main.tf` | Primary resource definitions |
| Module variables | `modules/<name>/variables.tf` | Input variables with types, descriptions, defaults |
| Module outputs | `modules/<name>/outputs.tf` | Exported values for composition |
| Module versions | `modules/<name>/versions.tf` | Provider version constraints |
| Module docs | `modules/<name>/README.md` | Usage, inputs, outputs, examples |
| Environment root | `environments/<env>/main.tf` | Root module composing child modules |
| Environment vars | `environments/<env>/terraform.tfvars` | Environment-specific variable values |
| Environment backend | `environments/<env>/backend.tf` | Remote state backend config (or via CLI) |
| Variable naming | `snake_case` | All Terraform variables |
| Output naming | `snake_case` with `_id`, `_name`, `_uri` suffixes | Consistent output naming |
| Resource naming | `azurerm_<resource>_<purpose>` | Terraform resource naming convention |
| Tag variable | `var.tags` (map(string)) | Merged with default tags in each module |

## Module Structure Convention
```
modules/<name>/
├── main.tf           # Primary resource definitions
├── variables.tf      # Input variables with types, descriptions, validation
├── outputs.tf        # Output values for composition
├── versions.tf       # Provider version constraints (required_providers)
├── providers.tf      # Provider configuration, aliases (if needed)
├── *.tf              # Additional resource files (logical grouping)
├── README.md         # Module documentation (inputs, outputs, usage, examples)
└── tests/            # (Optional) Module-specific Terratest tests
```

## Configuration Files
| File | Purpose | Location |
|------|---------|----------|
| `.terraform-version` | Terraform version pin for tfenv/asdf | Repo root |
| `.github/workflows/terraform-plan.yml` | PR plan workflow | `.github/workflows/` |
| `.github/workflows/terraform-apply.yml` | Merge/apply workflow | `.github/workflows/` |
| `.github/workflows/drift-detection.yml` | Scheduled drift detection | `.github/workflows/` |
| `.github/workflows/validate.yml` | PR validation (fmt, validate, lint) | `.github/workflows/` |
| `scripts/bootstrap.sh` | Bootstrap backend storage, SP | `scripts/` |
| `environments/*/backend.tf` | Remote state backend config | Per environment |
| `environments/*/providers.tf` | Provider config, aliases | Per environment |

## Where to Add New Code

**New Terraform Module:**
- Create `modules/<new-module>/` with standard structure above
- Add `versions.tf` with `required_providers` for `azurerm`, `azuread`, `random`, `time`, `tls` as needed
- Document inputs/outputs in `README.md`
- Add module tests in `modules/<new-module>/tests/` (optional but recommended)

**New Environment:**
- Create `environments/<new-env>/` with `main.tf`, `variables.tf`, `terraform.tfvars`, `backend.tf`
- Reference modules via `source = "../../modules/<name>"`
- Configure environment-specific variables in `terraform.tfvars`

**New Module Resource File:**
- Add `<feature>.tf` in `modules/<name>/` (e.g., `private-endpoints.tf` in `networking/`)
- Keep related resources grouped; split by logical concern

**New Policy Definition:**
- Add `definitions/<policy-name>.tf` in `modules/policies/definitions/`
- Add assignment in `modules/policies/assignments/<assignment-name>.tf`

**New CI/CD Workflow:**
- Add `.github/workflows/<name>.yml`
- Follow existing patterns: `terraform fmt`, `terraform validate`, `tflint`, `checkov` in validate workflow

**New Script:**
- Add `scripts/<name>.sh` with executable permissions
- Follow existing patterns: `set -euo pipefail`, logging functions, error handling

---

*Structure analysis: 2026-07-20*