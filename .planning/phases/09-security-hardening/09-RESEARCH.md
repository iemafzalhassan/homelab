# Phase 9: Security Hardening — Research

## 1. Overview & Objectives

This phase focuses on hardening the security posture of the AKS platform workloads. It ensures zero-trust network boundaries, restricts pod privileges to prevent host compromises, enforces resource boundaries to prevent denial-of-service (DoS) from resource-exhaustion, and provides a clear security baseline.

### Requirements Addressed
*   **SEC-01**: `NetworkPolicy` default-deny applied to all platform namespaces; explicit allow rules for required traffic only.
*   **SEC-02**: Pod security: no privileged containers, no `hostNetwork`, no `hostPID` on any platform component (except strictly necessary monitoring daemonsets).
*   **SEC-03**: `LimitRange` applied to all namespaces — no unbounded pods can consume the entire node.
*   **SEC-04**: `PodDisruptionBudget` on ArgoCD and Prometheus (deferred to a later phase as per architectural decisions).

---

## 2. NetworkPolicy Design (SEC-01)

### Core Behavioral Model
Kubernetes NetworkPolicies are **additive** (allow-only). If no policies apply to a pod, all ingress/egress is allowed. Once any policy selects a pod, that pod is isolated according to the rules in the policy: anything not explicitly allowed is denied.
To implement a true default-deny posture:
1.  Apply a default-deny-all policy (ingress and egress) to isolate the namespace.
2.  Add specific policies to permit required traffic.

For simplicity and maintainability, each platform namespace will have a **base policy** defining:
*   Default deny for all ingress and egress.
*   Explicit allow-all for same-namespace traffic (D-04).
*   Explicit allow-egress to CoreDNS (port 53 UDP/TCP in `kube-system`) to ensure name resolution.
*   Explicit allow-egress to Kubernetes API Server (port 443) for control-plane communication.

### Cross-Namespace Traffic Matrix

The table below outlines all permitted cross-namespace flows:

| Source Namespace | Source Pods | Destination Namespace | Destination Pods | Port | Protocol | Purpose |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `traefik` | Traefik Controller | `argocd` | `argocd-server` | `8080` / `443` | TCP | Ingress Routing |
| `traefik` | Traefik Controller | `jenkins` | `jenkins` (Controller) | `8080` | TCP | Ingress Routing |
| `traefik` | Traefik Controller | `keycloak` | `keycloak` | `8080` | TCP | Ingress Routing |
| `traefik` | Traefik Controller | `oauth2-proxy` | `oauth2-proxy` | `4180` | TCP | Ingress Routing |
| `oauth2-proxy` | `oauth2-proxy` | `keycloak` | `keycloak` | `8080` | TCP | OIDC Auth / Token Verification |
| `cnpg-system` | CNPG Operator | `keycloak` | `keycloak-db-*` | `8000` | TCP | PostgreSQL Instance Monitoring |
| `monitoring` | Grafana Alloy | `argocd` | All argocd metrics pods | `8080`, `8082`, `8083`, `8084` | TCP | Metrics Scraping |
| `monitoring` | Grafana Alloy | `jenkins` | `jenkins` (Controller) | `8080` | TCP | Metrics Scraping |
| `monitoring` | Grafana Alloy | `traefik` | `traefik-metrics` | `9100` | TCP | Metrics Scraping |
| `monitoring` | Grafana Alloy | `cert-manager` | `cert-manager` | `9402` | TCP | Metrics Scraping |
| `argocd`, `jenkins`, `traefik`, `keycloak`, `oauth2-proxy`, `cert-manager` | All Pods | `monitoring` | Grafana Alloy (Receiver) | `4317` (gRPC) / `4318` (HTTP) | TCP | OTLP Telemetry Push |
| `kube-system` / External | API Server | `cert-manager` | `cert-manager-webhook` | `10250` | TCP | Validation/Mutation Webhooks |
| `kube-system` / External | API Server | `cnpg-system` | `cloudnative-pg` | `9443` | TCP | Validation/Mutation Webhooks |

### Internet Egress Authorization (Ports 80 / 443)
The following workloads must be granted egress access to the external internet:
*   **`cert-manager`**: Port 443 to reach Let's Encrypt servers and Cloudflare API.
*   **`argocd`**: Port 443 to pull Helm charts and git repositories from external sources (e.g., GitHub).
*   **`jenkins` (Controller + Agents)**: Port 443 to pull build dependencies (Maven, npm, NuGet) and communicate with GitHub and Azure Container Registry (ACR).
*   **`monitoring` (Grafana Alloy)**: Port 443 to push metrics, logs, and traces to Grafana Cloud endpoints.

---

## 3. Pod Security Standards (SEC-02)

Kubernetes Pod Security Admission (PSA) enforces Pod Security Standards (PSS) at the namespace admission boundary using three levels: `privileged`, `baseline`, and `restricted`.

### Namespace Hardening Strategy
We will label the namespaces to enforce these restrictions. Because Kubernetes PSA does not support pod-level exemptions within a namespace:
1.  **`baseline` Enforced namespaces**: `argocd`, `jenkins`, `traefik`, `cert-manager`, `keycloak`, `oauth2-proxy`, `cnpg-system`. These namespaces will block any pods with privileged modes, host namespaces (`hostNetwork`, `hostPID`, `hostIPC`), or host path mounts.
2.  **`privileged` Enforced namespaces**: `monitoring` and `kube-system`.
    *   *Why?* The observability stack deploys `node-exporter` (requires `hostPID: true`, `hostNetwork: true`, `hostPath: /` to gather node statistics), `kepler` (requires privileged mode for eBPF power tracking), and Grafana Alloy logging daemonsets (requires `hostPath` to scrape `/var/log/pods`).
    *   *Mitigation*: We will set the `enforce` mode to `privileged` for `monitoring`, but apply `warn: baseline` and `audit: baseline` to ensure other non-privileged pods in the monitoring namespace are flagged if they violate baseline standards.

### Namespace Label Specifications
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: <namespace-name>
  labels:
    pod-security.kubernetes.io/enforce: baseline
    pod-security.kubernetes.io/enforce-version: latest
    pod-security.kubernetes.io/warn: baseline
    pod-security.kubernetes.io/warn-version: latest
```

---

## 4. Workload-level Security Contexts (SEC-02)

To satisfy the zero-privilege requirement, workloads must run with hardened `securityContext` parameters in their Helm charts or manifests:

### Standard Security Profile
Where supported by the charts, we will inject:
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000 # Or chart-specific unprivileged user (e.g., 65532 for distroless/Traefik)
  runAsGroup: 1000
  fsGroup: 1000
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop:
      - ALL
  seccompProfile:
    type: RuntimeDefault
```

### Component Adjustments
*   **Jenkins Controller**: Cannot use `readOnlyRootFilesystem: true` globally because the home directory (`/var/jenkins_home`) is active and writable. If `readOnlyRootFilesystem` is enabled, `/tmp` and other cache directories must be mapped to `emptyDir` volumes.
*   **Keycloak (Bitnami)**: Runs as user `1001` by default. Can set `fsGroup: 1001` and `runAsUser: 1001`.
*   **Traefik**: Binds to ports `8000` (web) and `8443` (websecure). Since these are non-privileged ports (>1024), Traefik does not need root privileges and runs as user `65532` by default.

---

## 5. LimitRange Profiling (SEC-03)

LimitRanges define default requests and limits for containers that do not declare them, preventing runaway pods from consuming node resources.

### Namespace Profile Grid
As per user implementation decisions (D-07), we configure the following values:

| Namespace | Default Request CPU | Default Request Memory | Default Limit CPU | Default Limit Memory |
| :--- | :--- | :--- | :--- | :--- |
| `jenkins` | `100m` | `256Mi` | `1000m` (1 Core) | `2Gi` |
| `argocd` | `100m` | `256Mi` | `500m` (0.5 Core) | `1Gi` |
| `keycloak` | `100m` | `256Mi` | `500m` (0.5 Core) | `1Gi` |
| `monitoring` | `100m` | `256Mi` | `500m` (0.5 Core) | `1Gi` |
| `traefik` | `50m` | `64Mi` | `250m` (0.25 Core) | `256Mi` |
| `cert-manager` | `50m` | `64Mi` | `250m` (0.25 Core) | `256Mi` |
| `oauth2-proxy` *(Proposed)* | `50m` | `64Mi` | `250m` (0.25 Core) | `256Mi` |
| `cnpg-system` *(Proposed)* | `100m` | `128Mi` | `500m` | `512Mi` |

---

## 6. PodDisruptionBudgets (SEC-04) — Deferral Rationale

### The Single-Node Eviction Trap
The Homelab AKS platform is provisioned with a single system node pool containing 1 node (`Standard_D2as_v5`).
If a PodDisruptionBudget (PDB) is applied to critical deployments (e.g., `argocd-server` or `prometheus`) with a rule such as `minAvailable: 1`, then:
*   When executing node maintenance or AKS upgrades, the eviction controller tries to drain the node.
*   To drain the node, it must evict the pods.
*   However, evicting the single pod replica would bring the availability of the deployment to `0`, violating the PDB rule (`minAvailable: 1`).
*   Because there are no other nodes in the system node pool to schedule the replacement pod *before* the current pod is evicted, the eviction will fail.
*   **Result**: The node drain will hang indefinitely, failing the cluster upgrade or node restart.

### Decision
Consequently, PDB enforcement (SEC-04) is explicitly **deferred** until a multi-node system pool is provisioned.

---

## 7. Critical Pitfalls & Mitigation Strategies

### Pitfall 1: The API Server Egress Block (Controllers break)
*   **The Trap**: Setting a strict default-deny egress NetworkPolicy without allowing access to the K8s API server blocks controller pods (ArgoCD, cert-manager, CNPG, Traefik, Alloy) from talking to the API server. This breaks state reconciliations, webhook calls, and dynamic discovery.
*   **Mitigation**: Include an explicit egress rule to port `443` in the default network policy for all controller namespaces. In AKS, the API server endpoint maps to the service `kubernetes.default.svc.cluster.local`, which resolves to the Service CIDR API IP (usually `10.0.0.1` or `10.96.0.1` or similar). Allowing TCP egress to port `443` globally or to the Service CIDR is critical.

### Pitfall 2: Webhook Admission Failure (Mutations fail)
*   **The Trap**: When the API server tries to call `cert-manager-webhook` (port 10250) or `cloudnative-pg` webhook (port 9443), and there is an ingress NetworkPolicy that blocks traffic from the control plane. This results in errors like `failed calling webhook` and blocks all resource creations.
*   **Mitigation**: Create specific ingress policies in `cert-manager` and `cnpg-system` namespaces to allow incoming TCP traffic on ports `10250` and `9443` respectively from any source IP (represented by `{}` in `ingress.from` rules) to ensure API Server webhook calls succeed.

### Pitfall 3: DNS Egress Block (Pods cannot resolve hostnames)
*   **The Trap**: Default-deny egress blocks DNS resolution (UDP/TCP port 53). Pods fail to resolve external endpoints (GitHub, Let's Encrypt) and internal services.
*   **Mitigation**: Inject an egress allowance targeting the `kube-system` namespace on port 53 (UDP/TCP) matching the label `k8s-app: kube-dns` in every NetworkPolicy.

### Pitfall 4: Read-Only Root Filesystem Crash
*   **The Trap**: Setting `readOnlyRootFilesystem: true` in workload securityContexts will cause some containers (e.g., Jenkins controllers, Keycloak) to crash because they attempt to write logs, caches, or plugins to local paths.
*   **Mitigation**: Carefully verify write paths for each workload. Use `emptyDir` volumes for temporary locations like `/tmp`, `/var/cache`, or `/run` when `readOnlyRootFilesystem: true` is enforced.

---

## 8. ArgoCD GitOps Integration Strategy

To manage these security configurations systematically:
1.  **Namespace Resources**: Define explicit `Namespace` YAML manifests in the repository for each namespace and include the PSA/PSS labels in them.
2.  **Manifest Placement**: Place `namespace.yaml`, `networkpolicy.yaml`, and `limitrange.yaml` in the respective bootstrap folders:
    *   `manifests/bootstrap/argocd/`
    *   `manifests/bootstrap/cert-manager/`
    *   `manifests/bootstrap/jenkins/`
    *   `manifests/bootstrap/keycloak/`
    *   `manifests/bootstrap/oauth2-proxy-extras/`
    *   `manifests/bootstrap/k8s-monitoring/`
    *   `manifests/bootstrap/traefik/` (New files)
    *   `manifests/bootstrap/cnpg/` (New files)
3.  **Bootstrap App Updates**:
    *   Update `argocd-extras.yaml` to include `namespace.yaml`, `networkpolicy.yaml`, and `limitrange.yaml` in its directory includes.
    *   Update `cert-manager-extras.yaml` similarly.
    *   Update `jenkins-extras.yaml` similarly.
    *   `keycloak-extras.yaml` already excludes only `values.yaml`, so adding files to its folder will auto-sync them.
    *   `oauth2-proxy-extras.yaml` uses `recurse: true`, so it will auto-sync the new manifests.
    *   Create **`traefik-extras`**, **`monitoring-extras`**, and **`cnpg-extras`** Application manifests in `manifests/apps/` to sync the new security resources for these namespaces.
