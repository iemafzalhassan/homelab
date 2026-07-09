# 11-02 Summary: Demo App & Kargo CRs

## Execution Summary
We have successfully set up the minimal demo application (`homelab-demo`) manifests, registered the target UAT and PROD environments in ArgoCD, and established Kargo Custom Resources for continuous promotion.

### Key Deliverables Completed
1. **Demo Application Manifests**:
   - Created deployment, service, and kustomization manifests for `demo-uat` and `demo-prod` namespaces.
   - Enforced spot pool scheduling (`kubernetes.azure.com/scalesetpriority: spot` tolerations and node selectors) on both deployments so they run on the spot node pool.
2. **ArgoCD Target Applications**:
   - Created `homelab-demo-uat` Application with automated sync.
   - Created `homelab-demo-prod` Application with manual sync only.
   - Annotated both applications with `kargo.akuity.io/authorized-stage` mapping them to Kargo stages.
3. **Kargo Custom Resources**:
   - Created cluster-scoped `Project` resource `homelab-demo`.
   - Created namespace-scoped `ProjectConfig` with promotion policies (auto-promote for UAT, manual for PROD).
   - Created `Warehouse` watching the Git repository path `manifests/apps/homelab-demo/`.
   - Created `Stage` `uat` performing `argocd-update` on the UAT application.
   - Created `Stage` `prod` performing `argocd-update` with Kustomize image overrides on the PROD application.
4. **Resources Deployment**:
   - Deployed the Kargo Custom Resources via a dedicated ArgoCD Application `kargo-resources`.
   - Labeled the namespace `homelab-demo` with `kargo.akuity.io/project=true` to pass Kargo webhook validations.

### Verification Status
- **Kargo CRs**: Warehouse, ProjectConfig, and Stages are active in the `homelab-demo` namespace.
- **ArgoCD Apps**: `homelab-demo-uat`, `homelab-demo-prod`, and `kargo-resources` are synced and healthy.
- **Scheduling**: Pod in `demo-uat` schedules successfully on the spot node pool (`aks-spot-40638830-vmss000006`) and is in pull-wait state.
