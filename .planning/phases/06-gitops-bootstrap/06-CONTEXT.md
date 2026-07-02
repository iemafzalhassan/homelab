# Phase 6: GitOps Bootstrap - Context

## Context Gathering Notes
*This phase was implemented manually prior to GSD planning to unblock cluster management. This context document retroactively records the implementation decisions made during that process to satisfy GSD workflow requirements.*

## Implementation Decisions

### 1. ArgoCD Installation Method
- **Decision:** Helm bootstrap via custom script (`install.sh`) rather than Terraform.
- **Rationale:** Keeps GitOps tools out of Terraform state. Terraform provisions the cluster, then the cluster connects to GitOps.
- **Details:** `manifests/bootstrap/argocd/values.yaml` disables the default HTTPRoute and we manage it explicitly.

### 2. Networking and Ingress
- **Decision:** Kubernetes Gateway API (`HTTPRoute`)
- **Rationale:** Traefik v3 is installed as the Gateway API controller.
- **Details:** `argocd-httproute.yaml` binds to the `traefik` Gateway to expose ArgoCD at `argocd.smapatticare.com` with TLS terminated by cert-manager.

### 3. Git Authentication
- **Decision:** SSH deploy key or PAT via secret.
- **Details:** A generic secret `argocd-repo-credentials` is applied to allow ArgoCD to fetch the private homelab repository.

### 4. App-of-Apps Pattern
- **Decision:** Self-managed ArgoCD with a root App.
- **Details:** `manifests/apps/root.yaml` deploys child applications for ArgoCD itself, Cert-Manager, Traefik, and Jenkins. Each child app resides in `manifests/apps/`.

## Open Questions
- None. This phase is already implemented and validated.
