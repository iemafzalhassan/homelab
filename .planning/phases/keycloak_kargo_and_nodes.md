# Phase: Cost Optimization, Keycloak SSO, and Kargo

## Overview
Before rolling out Keycloak and Kargo, the cluster infrastructure must be optimized to ensure it can handle the additional memory footprint while remaining within the Azure Student budget constraints and adhering to `GEMINI.md` project rules.

## Step 1: Infrastructure Node Optimization (Terraform)
1. **Retain System Node SKU**: Do NOT change the system node from `Standard_D2s_v3` to B-series (explicitly forbidden by `GEMINI.md` due to credit exhaustion).
2. **Enable Spot Nodes**: Modify the Spot Node pool in `modules/aks/main.tf` to have a `min_count = 1` so that a Spot node is always provisioned and ready for workloads.
3. Apply Terraform changes.

## Step 2: Workload Scheduling Adjustments
1. Update existing Helm values (ArgoCD, Jenkins, Traefik, Monitoring) to **remove** the `CriticalAddonsOnly` toleration.
2. Add the `kubernetes.azure.com/scalesetpriority: "spot"` nodeSelector and toleration to all these apps so they schedule exclusively on the Spot node.

## Step 3: Keycloak & Kargo Rollout
1. Deploy Keycloak with an embedded PostgreSQL database onto the Spot node.
2. Configure OIDC Clients in Keycloak for `argocd`, `jenkins`, `kargo`, and `grafana`.
3. Set up Keycloak Groups (`platform-admins`, `developers`, `viewers`) and users.
4. Deploy Kargo and configure it to use Keycloak OIDC.
5. Update ArgoCD, Jenkins, and Grafana configurations to authenticate via Keycloak OIDC and enforce RBAC mappings based on the Keycloak groups.

## Success Criteria
- Cluster system node is strictly running system pods (CoreDNS, Metrics Server) and utilizes < 50% Memory.
- Spot node is actively running all homelab applications (ArgoCD, Jenkins, Keycloak).
- A user logging in via Keycloak with the `viewers` role can access Grafana, but is explicitly denied access to Jenkins and ArgoCD.
