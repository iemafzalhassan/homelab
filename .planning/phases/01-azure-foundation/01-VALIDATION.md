# Phase 1: azure-foundation - Validation

**Created:** 2026-07-02

## Nyquist Validation Dimensions

*This checklist is used by `gsd-verify-phase` to audit the completed work. The executing agent must ensure all criteria are met before completing the phase.*

### 1. Requirements Completeness
- [ ] ALL phase requirements (INFRA-01 to INFRA-05) are fully implemented.
- [ ] No feature was "skipped for now" or left partially complete.

### 2. Functional Acceptance
- [ ] Root module works with `terraform apply` (INFRA-01).
- [ ] VNet with 4 subnets is created (INFRA-02).
- [ ] NSGs are attached and configured (INFRA-03).
- [ ] Budget Alert is active at 80% and 100% (INFRA-04).
- [ ] Tags are applied to resources (INFRA-05).

### 3. Edge Cases & Error Handling
- [ ] Terraform state backend (local) correctly creates state file.
- [ ] Missing variables produce clear Terraform errors.

### 4. Performance & Scale
- [ ] N/A for Phase 1.

### 5. Security & Privacy
- [ ] NSGs correctly restrict egress to essential ports, denying unauthorized traffic.

### 6. Architectural Consistency
- [ ] Terraform modules follow standard structure (`main.tf`, `variables.tf`, `outputs.tf`).
- [ ] Naming convention matches homelab standard (`homelab-*`).

### 7. Code Quality & Maintainability
- [ ] Terraform code is formatted (`terraform fmt`).
- [ ] Variables are well-documented.

### 8. Observability & Telemetry
- [ ] Cost budget alert is appropriately configured to notify the user.

## Manual UAT Script
```bash
terraform init
terraform validate
terraform plan
terraform apply -auto-approve
az network vnet show --name homelab-vnet --resource-group homelab-rg
az consumption budget show --resource-group homelab-rg --name homelab-budget
```
