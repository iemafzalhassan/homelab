# Phase 7: CI/CD (Jenkins) - Context

## Implementation Decisions

### 1. Jenkins Configuration Management
- **Decision:** Jenkins Configuration as Code (JCasC) via Helm `values.yaml`.
- **Rationale:** Aligns with GitOps principles. We will pre-install required plugins and configure the Kubernetes cloud provider for dynamic agents directly from the Helm chart values.

### 2. Jenkins Controller Persistence
- **Decision:** Stateful with an 8Gi PVC.
- **Rationale:** While JCasC defines the configuration, the controller needs persistent storage for build history, workspaces, and plugin caching to survive restarts without data loss.

### 3. Docker Build Strategy on Ephemeral Pods
- **Decision:** Kaniko for rootless image building.
- **Rationale:** Running Docker-in-Docker (DinD) in Kubernetes requires privileged pods, which violates security requirements (SEC-02). Kaniko allows building standard Docker images inside unprivileged agent pods without needing a Docker daemon.

### 4. Spot Node Agent Scheduling
- **Decision:** The Jenkins Kubernetes plugin will be configured to inject the specific spot node toleration (`kubernetes.azure.com/scalesetpriority=spot:NoSchedule`) into the JNLP agent pod templates.
- **Rationale:** Ensures Jenkins agents only run on the cheap spot instances, preserving the stability of the system node pool.
