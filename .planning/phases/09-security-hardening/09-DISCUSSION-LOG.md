# Phase 9: Security Hardening - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-07-09
**Phase:** 9-Security Hardening
**Areas discussed:** NetworkPolicy Scoping & Cross-Namespace Traffic, Namespace-level Pod Security Standards (PSS), LimitRange Resource Profiling, PodDisruptionBudget (PDB) Targets

---

## NetworkPolicy Scoping & Cross-Namespace Traffic

| Option | Description | Selected |
|--------|-------------|----------|
| Restrict egress to CoreDNS (port 53) and allow egress to the internet on ports 80/443 (HTTP/HTTPS) only | Restricts external/internet egress strictly to ports 80/443 + CoreDNS. | ✓ |
| Allow all egress to the internet | Standard egress allowed, ingress restricted. | |
| Deep egress lockdown | Egress strictly locked to specific destination CIDRs. | |

**User's choice:** Restrict egress to CoreDNS (port 53) and allow egress to the internet on ports 80/443 (HTTP/HTTPS) only.

---

## cross-namespace metrics scraping

| Option | Description | Selected |
|--------|-------------|----------|
| Restrict ingress strictly to specific metrics ports from the 'monitoring' namespace | Least privilege scraping strategy. | ✓ |
| Allow all ingress traffic from the 'monitoring' namespace | Open ingress from monitoring namespace. | |
| Restrict ingress to pods with a specific label | Only scrape specifically labeled pods. | |

**User's choice:** Restrict ingress strictly to specific metrics ports from the 'monitoring' namespace.

---

## ingress routing from 'traefik' ingress Gateway

| Option | Description | Selected |
|--------|-------------|----------|
| Allow ingress on target application ports strictly from pods in the 'traefik' namespace | Restrict target ports (80, 8080, 8443) to traefik namespace. | ✓ |
| Allow all ingress traffic from the 'traefik' namespace | Allow any ingress from Traefik. | |
| Allow ingress only to pods labeled as ingress-accessible | Label-based ingress routing checks. | |

**User's choice:** Allow ingress on target application ports strictly from pods in the 'traefik' namespace.

---

## intra-namespace traffic (communication within same namespace)

| Option | Description | Selected |
|--------|-------------|----------|
| Allow all intra-namespace traffic by default | Pods in same namespace communicate freely. | ✓ |
| Deny all intra-namespace traffic by default | Deny by default, explicit allow per pod. | |

**User's choice:** Allow all intra-namespace traffic by default.

---

## Namespace-level Pod Security Standards (PSS)

| Option | Description | Selected |
|--------|-------------|----------|
| Enforce baseline PSS labels on namespaces | Prevent privileged/host-access pods. | ✓ |
| Workload-only securityContext configuration | No admission controls on namespace. | |
| Enforce restricted PSS labels everywhere | Force all pods to run as non-root/unprivileged. | |

**User's choice:** Enforce baseline PSS labels on namespaces.

---

## PSS Exceptions

| Option | Description | Selected |
|--------|-------------|----------|
| Apply baseline PSS to all platform namespaces; allow exceptions strictly for system/monitoring daemonsets | Node-exporter/Alloy daemonsets allowed host access. | ✓ |
| Apply baseline PSS with zero exceptions | Strictly block any privileged/host access pods. | |
| Exclude the 'jenkins' namespace | Keep jenkins namespace completely free from PSS. | |

**User's choice:** Apply baseline PSS to all platform namespaces; allow exceptions strictly for system/monitoring daemonsets.

---

## workload-level securityContext parameters

| Option | Description | Selected |
|--------|-------------|----------|
| Enforce standard unprivileged settings (runAsNonRoot: true, allowPrivilegeEscalation: false, runAsUser: 1000) | Apply via Helm values/manifests. | ✓ |
| Leave workload securityContext settings optional | Rely purely on PSS. | |

**User's choice:** Enforce standard unprivileged settings.

---

## LimitRange Resource Profiling

**User's feedback:** "analyse the running applications and based on each application you need to set limit and rages"

**Outcome:** Customized LimitRanges based on pod footprints fetched from the running cluster.
- `jenkins` namespace: default request `256Mi` RAM / `100m` CPU, limit `2Gi` RAM / `1000m` CPU.
- `argocd`, `monitoring`, `keycloak`: default request `256Mi` RAM / `100m` CPU, limit `1Gi` RAM / `500m` CPU.
- `traefik`, `cert-manager`: default request `64Mi` RAM / `50m` CPU, limit `256Mi` RAM / `250m` CPU.

---

## PodDisruptionBudget (PDB) Targets

**User's feedback:** "skip this now will plan later document it will leverage new noes ffor maintance and all will do this planning lkater and execution later on after planning skip it for now"

**Outcome:** Deferred. PDB implementation (SEC-04) will be planned and executed later when additional nodes are added to leverage node maintenance without blocking.

---

## Claude's Discretion

- Exact YAML layout of NetworkPolicy and LimitRange resources.
- Workload-specific `securityContext` tuning within their respective Helm values files.

## Deferred Ideas

- **PodDisruptionBudget (PDB) Targets (SEC-04)**: Defer implementing PDBs in this phase. PDB planning and execution are deferred to a later stage when additional nodes are added for maintenance and high availability, ensuring evictions do not block cluster operations on single-node setups.
