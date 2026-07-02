# Phase 1 Research: Azure Foundation

## Goal
Establish the foundational Azure infrastructure (VNet, subnets, NSGs, and budget alerts) and the Terraform root module with local state.

## Implementation Details

### Terraform Structure
Since this is a homelab with a single environment, a flat Terraform structure is easiest to manage.
- `main.tf`: Provider configuration (`azurerm`), backend configuration (local state), and module calls.
- `variables.tf`: Input variables (e.g., location, budget_contact_email).
- `outputs.tf`: Important outputs like VNet ID and Subnet IDs.

### Networking Module (`modules/networking`)
- **VNet**: Address space `10.0.0.0/16`.
- **Subnets**:
  - `system-subnet`: `10.0.1.0/24`
  - `user-subnet`: `10.0.2.0/24`
  - `infra-subnet`: `10.0.4.0/24`
  - `ingress-subnet`: `10.0.5.0/24` (Wait, architecture specified only 3 subnets in architecture diagram: system 10.0.1.0/24, user 10.0.2.0/24, infra 10.0.4.0/24. But REQUIREMENTS.md INFRA-02 says: "4 subnets (system, user/spot, infra, ingress)". We will provision all 4.)
- **NSGs**: Create one NSG per subnet, or a shared one if rules are identical, but best practice is per subnet. Default-deny egress is tricky because AKS needs outbound internet access to pull images.
- **NSG Default Deny Egress**: AKS requires outbound access to Azure APIs, MCR, Ubuntu servers. A strict default-deny egress NSG will break AKS unless Azure Firewall or NAT Gateway is used. For a $100 budget homelab, NAT Gateway costs $32/mo and Azure Firewall is $300/mo. Therefore, outbound internet MUST be allowed on ports 80/443/UDP 123 (NTP) at a minimum, or we use standard `AllowInternetOutBound` rule.
*Correction for INFRA-03*: To keep costs low, we will allow standard outbound internet (HTTP/HTTPS) and block everything else.

### Budget Alert
- Resource: `azurerm_consumption_budget_resource_group`
- Thresholds: 80% and 100% of $25/mo (since total budget is $100 for 4 months or $150 for 6 months. Let's set amount to $25/mo).
- Contact email will be passed as a variable.

## Validation Architecture
- Validate that `terraform init` and `terraform plan` succeed.
- Validate that applying the networking module creates the correct VNet and Subnets.
- Validate tags `environment=homelab`, `managed-by=terraform` are present.

## Key Takeaways for Planning
1. Create `modules/networking` with VNet, Subnets, and NSGs.
2. Create `modules/budget` (or put in root) for `azurerm_consumption_budget_resource_group`.
3. Set up the root module using the local backend.
4. Define standard tags as local variables in the root module.
