# Phase 9: Security Hardening - Context

**Gathered:** 2026-07-09
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase focuses on security hardening of the AKS platform workloads. It includes applying NetworkPolicies (default-deny egress/ingress and fine-grained allow rules), namespace-level Pod Security Standards (PSS), workload-level securityContext parameters, and LimitRange resource constraints across the platform namespaces (`argocd`, `jenkins`, `traefik`, `monitoring`, `cert-manager`, `keycloak`).

</domain>

<decisions>
## Implementation Decisions

### NetworkPolicy Scoping & Cross-Namespace Traffic
- **D-01:** Deny all egress by default. Restrict egress to CoreDNS (port 53 in `kube-system`) and allow egress to the internet on ports 80/443 (HTTP/HTTPS) only for necessary external platform connections (e.g., GitHub, cert-manager to Let's Encrypt, Jenkins agents pulling packages).
- **D-02:** Allow ingress from the `monitoring` namespace (Grafana Alloy) strictly to specific metrics ports (e.g., 8080, 9090, 8083, etc.) on pods in targeted namespaces for scraping.
- **D-03:** Allow ingress on target application ports (e.g., 80, 8080, 8443) strictly from pods in the `traefik` namespace to applications in other namespaces.
- **D-04:** Allow all intra-namespace traffic by default (communication within the same namespace is unrestricted).

### Namespace-level Pod Security Standards (PSS)
- **D-05:** Enforce `baseline` PSS labels (`pod-security.kubernetes.io/enforce: baseline`) on all platform namespaces (`argocd`, `jenkins`, `traefik`, `monitoring`, `cert-manager`, `keycloak`) to prevent running privileged or host-access pods.
- **D-06:** Allow exceptions strictly for system/monitoring daemonsets (like `node-exporter` or Grafana Alloy/Beyla if required) that necessitate host-level permissions, confined to the `monitoring` or `kube-system` namespaces.

### LimitRange Resource Profiling
- **D-07:** Establish customized LimitRange default resource profiles tailored to the footprint of the workloads in each namespace (applied to pods that do not declare their own CPU/memory requests/limits):
  - **`jenkins` namespace**: default request `256Mi` RAM / `100m` CPU, limit `2Gi` RAM / `1000m` CPU.
  - **`argocd`, `monitoring`, `keycloak` namespaces**: default request `256Mi` RAM / `100m` CPU, limit `1Gi` RAM / `500m` CPU.
  - **`traefik`, `cert-manager` namespaces**: default request `64Mi` RAM / `50m` CPU, limit `256Mi` RAM / `250m` CPU.

### Claude's Discretion
- Exact YAML structure of NetworkPolicy and LimitRange resources.
- Workload-specific `securityContext` tuning within their respective Helm values files.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Specifications
- [PROJECT.md](file:///Users/iemafzal-mac/Library/Mobile%20Documents/com~apple~CloudDocs/Homelab/.planning/PROJECT.md) — Security & Networking guidelines
- [REQUIREMENTS.md](file:///Users/iemafzal-mac/Library/Mobile%20Documents/com~apple~CloudDocs/Homelab/.planning/REQUIREMENTS.md) §SEC — Requirement specifications SEC-01 to SEC-04
- [ROADMAP.md](file:///Users/iemafzal-mac/Library/Mobile%20Documents/com~apple~CloudDocs/Homelab/.planning/ROADMAP.md) §Phase 9 — Success criteria for Security Hardening

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `manifests/bootstrap/*/values.yaml` — Workload Helm chart configurations (ArgoCD, Traefik, Jenkins, cert-manager).
- Existing manifests in [manifests/apps/](file:///Users/iemafzal-mac/Library/Mobile%20Documents/com~apple~CloudDocs/Homelab/manifests/apps) which deploy platform applications via ArgoCD.

### Established Patterns
- All cluster resources are managed via ArgoCD using GitOps.
- Customized resource configs are located under `manifests/bootstrap/` subdirectories and deployed through ArgoCD Application manifests.

### Integration Points
- Place new NetworkPolicy, LimitRange, and Pod Security Admission label resources in their respective namespace bootstrap folders (e.g., under [manifests/bootstrap/](file:///Users/iemafzal-mac/Library/Mobile%20Documents/com~apple~CloudDocs/Homelab/manifests/bootstrap)).

</code_context>

<specifics>
## Specific Ideas

- None — open to standard Kubernetes security patterns.

</specifics>

<deferred>
## Deferred Ideas

- **PodDisruptionBudget (PDB) Targets (SEC-04)**: Defer implementing PDBs in this phase. PDB planning and execution are deferred to a later stage when additional nodes are added for maintenance and high availability, ensuring evictions do not block cluster operations on single-node setups.

</deferred>

---

*Phase: 9-Security Hardening*
*Context gathered: 2026-07-09*
