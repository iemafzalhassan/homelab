# Phase 03: Identities & Secrets - Context

**Gathered:** 2026-07-02
**Status:** Ready for planning

<domain>
## Phase Boundary

Creation of User-Assigned Managed Identities (for ArgoCD, Jenkins, Traefik), federated credentials scoped to their Kubernetes ServiceAccounts, and an Azure Key Vault accessed privately via a Private Endpoint.
</domain>

<decisions>
## Implementation Decisions

### Module Structure
- **D-01:** Use separate `identity` and `keyvault` modules to ensure better separation of concerns as complexity grows.

### Secret Management
- **D-02:** Provision placeholder secrets directly in Terraform so that we can easily test the Secrets Store CSI driver in Phase 4.

### API Server Access (Laptop IP)
- **D-03:** Update the `admin_ip_range` configuration in `terraform.tfvars` (and underlying modules) to support multiple IPs, permanently adding the laptop's public IP alongside the main PC's IP.

### the agent's Discretion
- Key Vault Naming: Append a random string or suffix to ensure the Key Vault name is globally unique.
- RBAC configurations and exact scopes for the Key Vault access.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture & Requirements
- `.planning/PROJECT.md` — Project context and constraints
- `.planning/REQUIREMENTS.md` — Phase 3 requirements (IDENTITY-01 to IDENTITY-05)
</canonical_refs>
