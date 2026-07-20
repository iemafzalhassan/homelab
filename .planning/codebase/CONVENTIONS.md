# Code Conventions

> Mapped: 2026-07-20

## Code Style
- **Formatter**: prettier — `.prettierrc`
- **Linter**: eslint — `eslint.config.mjs`
- **Line Length**: 100 chars — `.prettierrc`

## Naming Conventions
| Element | Convention | Example |
|---------|------------|---------|
| Variables | camelCase | `clusterName`, `resourceGroupName` |
| Functions | camelCase | `createResourceGroup`, `getVnetId` |
| Classes/Types | PascalCase | `ResourceGroupModule`, `NetworkConfig` |
| Constants | UPPER_SNAKE_CASE | `DEFAULT_LOCATION`, `MAX_RETRIES` |
| Files | kebab-case | `main.tf`, `variables.tf`, `network-module.tf` |
| Terraform Resources | snake_case | `azurerm_resource_group`, `azurerm_virtual_network` |
| Terraform Modules | kebab-case | `modules/network`, `modules/compute` |
| Kubernetes Resources | kebab-case | `deployment.yaml`, `service-account.yaml` |
| Kubernetes Resources (kind) | PascalCase | `Deployment`, `ServiceAccount`, `ClusterRoleBinding` |

## Terraform Patterns
| Pattern | Usage | Example Location |
|---------|-------|------------------|
| Module composition | Root modules compose child modules | `infra/main.tf` → `modules/` |
| Variable validation | Validation blocks in variables | `modules/*/variables.tf` |
| Dynamic blocks | Dynamic resource configuration | `modules/network/main.tf` |
| For_each loops | Resource iteration | `modules/compute/main.tf` |
| Locals for computed values | Computed values and transformations | `modules/*/locals.tf` |
| Data sources for references | Reference existing resources | `modules/*/data.tf` |
| Outputs for module outputs | Module output values | `modules/*/outputs.tf` |
| Variable validation blocks | Input validation with regex/conditions | `modules/*/variables.tf` |

## Terraform Module Structure
```
modules/
├── <module-name>/
│   ├── main.tf          # Main resource definitions
│   ├── variables.tf     # Input variables with validation
│   ├── outputs.tf       # Output values
│   ├── locals.tf        # Computed values (optional)
│   ├── data.tf          # Data sources (optional)
│   ├── versions.tf      # Provider versions (optional)
│   └── README.md        # Module documentation (optional)
```

## Error Handling Patterns
| Pattern | Example Location |
|---------|------------------|
| Variable validation blocks | `modules/*/variables.tf` — regex, length, custom conditions |
| Precondition/Postcondition blocks | `modules/*/main.tf` — resource lifecycle rules |
| Try/catch in locals | `modules/*/locals.tf` — try() for optional values |
| Default values with fallback | `modules/*/variables.tf` — default = null + coalesce |

## Documentation Patterns
| Artifact | Location | Format |
|----------|----------|--------|
| Module README | `modules/*/README.md` | Markdown with inputs/outputs tables |
| Root module README | `infra/README.md`, `bootstrap/README.md` | Markdown with usage examples |
| Architecture Decision Records | `docs/adr/` | Markdown (ADR format) |
| Runbooks | `docs/runbooks/` | Markdown with procedures |

## Git Practices
| Practice | Convention |
|----------|------------|
| Commit format | Conventional Commits: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:` |
| Branch naming | `feature/<name>`, `fix/<name>`, `chore/<name>`, `docs/<name>` |
| PR template | `.github/pull_request_template.md` |
| Commit signing | Signed commits required (GPG/SSH) |
| Branch protection | Main branch protected, PR required |

## Terraform Version Constraints
- **Provider**: `azurerm` ~> 4.0
- **Terraform**: >= 1.5.0
- **Backend**: Azure Storage (remote state)

## Naming Standards for Azure Resources
| Resource Type | Naming Pattern |
|---------------|----------------|
| Resource Group | `rg-<env>-<workload>-<region>` |
| Virtual Network | `vnet-<env>-<workload>-<region>` |
| Subnet | `snet-<purpose>` |
| Storage Account | `st<env><workload><region>` (lowercase, no hyphens) |
| Key Vault | `kv-<env>-<workload>-<region>` |
| Log Analytics | `log-<env>-<workload>-<region>` |
| Key Vault Secret | `<service>-<purpose>-<env>` |
| Key Vault Key | `<service>-<purpose>-<env>` |
| Key Vault Certificate | `<service>-<purpose>-<env>` |

## Terraform Code Organization
- **Root modules**: `infra/`, `bootstrap/` — environment-specific compositions
- **Child modules**: `modules/` — reusable, parameterized components
- **Bootstrap**: `bootstrap/` — bootstrap resources (storage account for state, key vault)
- **Environments**: Separate state files per environment via workspaces or separate state files

## Variable Patterns
- Use `variable` blocks with `type`, `description`, `validation` blocks
- Use `default = null` for optional variables, handle with `coalesce()` in locals
- Use `validation` blocks with `condition` and `error_message` for input validation
- Group related variables with `type = object({ ... })`

## Output Patterns
- All modules have `outputs.tf` with meaningful descriptions
- Output sensitive values with `sensitive = true`
- Use `depends_on` only when implicit dependencies aren't sufficient

## Provider Configuration
- Provider versions pinned in `versions.tf` or root `versions.tf`
- Azure provider configured with `features {}` block
- Use `azuread` provider for Azure AD resources
- Use `random` provider for random strings/integers
- Use `tls` provider for TLS certificates/keys

---

*Convention analysis: 2026-07-20*