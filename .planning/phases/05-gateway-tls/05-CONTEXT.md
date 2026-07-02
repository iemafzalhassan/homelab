# Phase 5: Gateway & TLS - Context

**Gathered:** 2026-07-02
**Status:** Ready for planning

<domain>
## Phase Boundary

External HTTPS traffic routing into the cluster using Traefik v3 (with Gateway API instead of Ingress), automatic wildcard TLS certificates via cert-manager (with Cloudflare DNS-01 challenges), and cross-namespace routing using ReferenceGrants.
</domain>

<decisions>
## Implementation Decisions

### Deployment Mechanism
- **D-01:** Use Bash bootstrap scripts (using `helm upgrade --install`) for deploying Traefik and cert-manager for now. This keeps infrastructure decoupled from Terraform and makes it easy to migrate these workloads to ArgoCD in Phase 6.

### Gateway API CRDs Channel
- **D-02:** Install the **Standard Channel** of the Kubernetes Gateway API CRDs for maximum stability in the homelab.

### Cloudflare API Token Secrets Sync
- **D-03:** The Cloudflare API Token (stored in Azure Key Vault) will be synced to a Kubernetes Secret **only in the `cert-manager` namespace** (least privilege). A `SecretProviderClass` and a CSI volume mount (using a dummy/daemonset pod or similar workload) will be created in that namespace.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture & Requirements
- `.planning/PROJECT.md` — Project context and constraints
- `.planning/REQUIREMENTS.md` — Phase 5 requirements (GATEWAY-01 to GATEWAY-06, TLS-01 to TLS-05)
- `research/STACK.md` — Technology Stack and Verified Component Versions (Traefik v3, Gateway API, cert-manager versions)
</canonical_refs>
