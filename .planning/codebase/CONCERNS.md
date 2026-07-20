# Technical Concerns & Debt

> Mapped: 2026-07-20

## Technical Debt

| Area | Description | Severity | Location | Evidence |
|------|-------------|----------|----------|----------|
| Incomplete Bicep Module Implementation | Several Bicep modules are stubs or incomplete (e.g., `network/private-endpoints.bicep`, `network/vnet-peering.bicep` are near-empty) | High | `infra/bicep/modules/network/` | `private-endpoints.bicep` has only comments, `vnet-peering.bicep` is ~5 lines |
| Missing Bicep Module Implementations | Core modules like `container-app-env`, `container-app`, `log-analytics`, `key-vault` are referenced but implementations are missing/stubbed | High | `infra/bicep/modules/` | Multiple module directories contain only `README.md` or empty `.bicep` files |
| Incomplete Azure DevOps Pipeline | Pipeline YAML references templates that don't exist (`templates/` directory missing) | High | `infra/pipelines/azure-pipelines.yml` | References `templates/deploy-bicep.yml` which doesn't exist |
| Missing Bicep Module Tests | No test infrastructure for Bicep modules (no `bicep-linter` config, no test framework) | Medium | `infra/bicep/` | No `.bicepconfig.json`, no test directory |
| Incomplete Azure Policy Definitions | Policy definitions directory exists but is empty | Medium | `infra/policies/` | Directory exists but contains no policy definitions |
| Missing GitHub Actions Workflows | `.github/workflows/` directory doesn't exist despite references in docs | Medium | `.github/workflows/` | Directory doesn't exist |
| Incomplete Documentation | Multiple `.md` files in `docs/` are stubs or TODOs | Medium | `docs/` | Multiple files have "TODO" or minimal content |
| Missing Environment Configuration | `.env.example` references variables not documented in `ENVIRONMENT.md` | Medium | `.env.example`, `docs/ENVIRONMENT.md` | Variables like `AZURE_CLIENT_ID` referenced but not documented |
| Missing Bicep Linter Configuration | No `.bicepconfig.json` for linting rules | Medium | `infra/bicep/` | No linting configuration found |
| Missing Test Infrastructure | No test framework for Bicep, no Terratest, no Terratest equivalents for Bicep | Medium | `infra/` | No test directories or configs found |

## Known Bugs / Issues

| Issue | Impact | Location | Status |
|-------|--------|----------|--------|
| Bicep module `container-app-env.bicep` referenced but missing | Pipeline deployment fails | `infra/bicep/modules/app/container-app-env.bicep` | Open |
| Bicep module `container-app.bicep` referenced but missing | Container app deployment fails | `infra/bicep/modules/app/container-app.bicep` | Open |
| Bicep module `log-analytics.bicep` referenced but missing | Monitoring deployment fails | `infra/bicep/modules/monitoring/log-analytics.bicep` | Open |
| Bicep module `key-vault.bicep` referenced but missing | Secrets/key management fails | `infra/bicep/modules/security/key-vault.bicep` | Open |
| Pipeline references non-existent templates | CI/CD pipeline fails | `infra/pipelines/azure-pipelines.yml` | Open |
| `.env.example` references `AZURE_TENANT_ID` but docs don't explain how to obtain | Onboarding friction | `.env.example`, `docs/ENVIRONMENT.md` | Open |
| `docs/ARCHITECTURE.md` references diagrams that don't exist | Documentation incomplete | `docs/ARCHITECTURE.md` | Open |
| `docs/DEPLOYMENT.md` references Azure DevOps pipelines that don't exist | Deployment docs misleading | `docs/DEPLOYMENT.md` | Open |

## Security Concerns

| Concern | Severity | Location | Mitigation |
|---------|----------|----------|------------|
| No Azure Policy definitions for security baselines | High | `infra/policies/` | Implement Azure Policy definitions for: allowed locations, allowed resource types, encryption at rest, network restrictions |
| No Key Vault module implemented | High | `infra/bicep/modules/security/key-vault.bicep` | Implement Key Vault module with RBAC, purge protection, soft delete |
| No Private Endpoint implementations | High | `infra/bicep/modules/network/private-endpoints.bicep` | Implement private endpoints for all PaaS services |
| No Network Security Group (NSG) rules defined | Medium | `infra/bicep/modules/network/nsg.bicep` | Define NSG rules for least-privilege network access |
| No Azure AD / Entra ID integration module | Medium | `infra/bicep/modules/identity/` (missing) | Implement Entra ID integration for workload identity |
| No secret scanning configuration | Medium | `.github/` (missing) | Add GitHub secret scanning, TruffleHog, or similar |
| No dependency scanning for Bicep modules | Medium | `infra/bicep/` | Add Bicep linter, check for deprecated API versions |
| `.env.example` contains placeholder secrets pattern | Low | `.env.example` | Ensure no real secrets committed; add pre-commit hooks |

## Performance Concerns

| Area | Impact | Location | Evidence |
|------|--------|----------|----------|
| No Bicep module compilation caching in CI | Slower CI runs | `infra/pipelines/` (missing) | No caching strategy in pipelines |
| No Bicep module registry usage | Slower module resolution, no versioning | `infra/bicep/` | Modules referenced via relative paths only |
| No Azure Resource Graph queries for inventory | Slow inventory/audit | `infra/` (missing) | No inventory automation |
| No deployment parallelization strategy | Sequential deployments slow | `infra/pipelines/azure-pipelines.yml` | Pipeline runs stages sequentially |
| No cost estimation in CI/CD | Unexpected cost surprises | `infra/pipelines/` | No cost estimation step in pipeline |
| No Bicep what-if analysis in pipeline | Risk of unintended changes | `infra/pipelines/azure-pipelines.yml` | No `az deployment what-if` step |

## Fragile Areas

| Area | Why Fragile | Location | Risk |
|------|-------------|----------|------|
| Bicep module reference chains | Relative paths break easily when moving files | `infra/bicep/main.bicep`, `infra/bicep/modules/` | High - moving any module breaks references |
| Pipeline template references | Missing templates cause pipeline failures | `infra/pipelines/azure-pipelines.yml` | High - pipeline fails completely |
| Environment variable references | `.env.example` vars not validated | `.env.example`, `docs/ENVIRONMENT.md` | Medium - deployment fails silently |
| Bicep module versioning | No versioning, breaking changes affect all | `infra/bicep/modules/` | Medium - no version pinning |
| Hardcoded resource names in modules | Naming collisions in shared environments | `infra/bicep/modules/` | Medium - no unique suffix generation |
| No Bicep linter enforcement | Code quality degrades silently | `infra/bicep/` (no config) | Medium - style drift |
| Documentation drift | Docs reference non-existent resources | `docs/` | Medium - misleading docs |
| Environment variable validation | No validation at deploy time | `infra/pipelines/` | Medium - late failures |

## Missing Tests / Coverage Gaps

| Area | Missing Coverage | Location | Risk |
|------|------------------|----------|------|
| Bicep module unit tests | No test framework for Bicep | `infra/bicep/modules/` | High - no validation of module logic |
| Bicep integration tests | No test deployments | `infra/` | High - no validation of full deployments |
| Pipeline integration tests | No pipeline test runs | `infra/pipelines/` | High - pipeline changes untested |
| Policy compliance tests | No policy tests | `infra/policies/` | Medium - no compliance validation |
| Infrastructure drift detection | No drift detection | `infra/` | Medium - drift undetected |
| Cost regression tests | No cost tracking in CI | `infra/pipelines/` | Medium - cost surprises |
| Security policy tests | No policy-as-code tests | `infra/policies/` | High - security regressions undetected |
| Disaster recovery tests | No DR validation | `infra/` | High - DR untested |

## Deprecated / Legacy Code

| Component | Replacement | Location | Migration Status |
|-----------|-------------|----------|------------------|
| Azure DevOps Pipelines (primary) | GitHub Actions (planned) | `infra/pipelines/azure-pipelines.yml` | Planned - referenced in docs but not implemented |
| Manual Azure CLI deployments | Bicep/IaC | `docs/DEPLOYMENT.md` | Partial - Bicep modules started but incomplete |
| Manual resource creation docs | IaC modules | `docs/` | Incomplete - docs reference manual steps |

## Configuration Drift Risks

| Area | Risk | Location |
|------|------|----------|
| Bicep parameter files per environment | Drift between `dev.bicepparam`, `prod.bicepparam` | `infra/bicep/params/` |
| Environment variables in `.env.example` vs actual deployment | Drift between local dev and deployed env | `.env.example`, `infra/bicep/params/` |
| Azure Policy definitions vs deployed policies | Drift between policy code and assigned policies | `infra/policies/` (empty) |
| Documentation vs actual architecture | Docs reference non-existent components | `docs/ARCHITECTURE.md`, `docs/DEPLOYMENT.md` |
| Pipeline templates vs actual pipeline | Referenced templates don't exist | `infra/pipelines/azure-pipelines.yml` |
| Bicep module versions vs deployed versions | No version pinning | `infra/bicep/modules/` |

## Operational Concerns

| Concern | Impact | Location |
|---------|--------|----------|
| No deployment validation (what-if) in pipeline | Risk of unintended changes | `infra/pipelines/azure-pipelines.yml` |
| No automated rollback strategy | Failed deployments require manual intervention | `infra/pipelines/` |
| No deployment approval gates | No manual approval for prod | `infra/pipelines/azure-pipelines.yml` |
| No infrastructure cost monitoring/alerting | Unexpected cost spikes | `infra/` (missing) |
| No infrastructure drift detection | Configuration drift undetected | `infra/` (missing) |
| No backup/DR validation for infrastructure | DR untested | `infra/` |
| No centralized logging for deployment pipeline | Debugging failures difficult | `infra/pipelines/` |
| No secret rotation automation | Secrets may become stale | `infra/bicep/modules/security/` (missing) |
| No infrastructure inventory/asset management | Unknown resources in subscription | `infra/` (missing) |
| No compliance reporting automation | Manual compliance audits | `infra/policies/` (empty) |

---

*Concerns audit: 2026-07-20*