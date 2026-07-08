# Homelab Node Utilization & Cost Optimization Research

## Current State Analysis
I analyzed the current AKS cluster utilization (`kubectl top nodes` / `pods`) and the Terraform node pool definitions.

### Node Utilization
- **Single Active Node (`aks-system-25810623-vmss000000`)**: Running at **89% Memory (5.16 GiB)** and 34% CPU (661m).
- **Spot Node Pool (`Standard_D2as_v5`)**: Currently at **0 nodes** because no workloads are forcing a scale-up.

### The Misconfiguration
The **System Node Pool** (`Standard_D2s_v3`) is defined with `only_critical_addons_enabled = true`. In Azure, this adds a taint (`CriticalAddonsOnly=true:NoSchedule`) designed to prevent non-system workloads from running there.

However, your application Helm values (Traefik, ArgoCD, Jenkins, and k8s-monitoring) have all been manually given the `CriticalAddonsOnly` toleration. This completely bypasses the intended node segregation. Heavy workloads like Jenkins and ArgoCD are running on your expensive on-demand system node, pushing it to 89% memory. 

If we deploy Keycloak and Kargo right now, the system node will hit 100% memory and start OOM-killing (crashing) pods.

---

## Cost-Optimized Target Architecture

To maximize your Azure Student credits (~$10-15/mo budget) while adhering to production best practices (as defined in your `GEMINI.md`), we must redesign the scheduling strategy:

### 1. Maintain System Node Integrity (D-Series Only)
- **Current**: `Standard_D2s_v3` (~$70/mo on-demand).
- **Validation**: Per `GEMINI.md`, we must **NOT** switch to a B-series VM (Burstable) because CPU credit exhaustion causes CoreDNS and kubelet failures. 
- **Action**: Keep the `Standard_D2s_v3` system node. Keep `only_critical_addons_enabled = true`. By removing user workloads from it, it will run extremely light and stable.

### 2. Properly Utilize the Spot Node Pool for Workloads
- **Current**: `Standard_D2as_v5` Spot (2 vCPU, 8GB RAM). Spot instances are incredibly cheap (~$10/mo) and perfectly suited for your homelab workloads.
- **Action**: Set `min_count = 1` in Terraform so at least one Spot node is always available for workloads to schedule onto without waiting for cluster-autoscaler timeouts.

### 3. Fix Application Scheduling (Helm/Manifests)
- **Remove** the `CriticalAddonsOnly` toleration from all homelab workloads (Jenkins, ArgoCD, Traefik, Monitoring, and the upcoming Keycloak/Kargo).
- **Add** the Spot toleration and NodeSelector/Affinity to all workloads:
  ```yaml
  nodeSelector:
    kubernetes.azure.com/scalesetpriority: "spot"
  tolerations:
    - key: "kubernetes.azure.com/scalesetpriority"
      operator: "Equal"
      value: "spot"
      effect: "NoSchedule"
  ```

## Conclusion
By adopting this structure, the System Node is strictly reserved for CoreDNS, Metrics Server, and Azure CNI, as originally intended. All user workloads (including the upcoming Keycloak and Kargo) will run exclusively on the much cheaper Spot nodes, giving you 8GB+ of workable RAM at a fraction of the cost, safely supporting the new multi-environment SSO requirements.
