# Phase 8: Observability - Execution Plan

## 1. Requirements Traceability

| Requirement ID | Description | Resolution in Plan |
|---|---|---|
| **OBS-01** | kube-prometheus-stack deployed with 30s scrape interval | **Task 2 & 3:** Replaced with `grafana/k8s-monitoring` Helm chart (Alloy) per `08-CONTEXT.md` to respect 8GB node budget. |
| **OBS-02** | Grafana accessible via HTTPRoute at grafana.yourdomain.com | **Task N/A:** Handled natively by Grafana Cloud UI. |
| **OBS-03** | Dashboards for core platform metrics | **Task 4:** Core metrics shipped to Grafana Cloud where default dashboards exist. |
| **OBS-04** | Alertmanager configured with alert rule | **Task N/A:** Handled natively by Grafana Cloud Alerting (routed to Slack). |
| **OBS-05** | Grafana dashboards stored as ConfigMaps | **Task N/A:** Using Grafana Cloud's hosted dashboards instead of local ConfigMaps. |

## 2. Architecture & Design

Following the `08-CONTEXT.md` decision to avoid out-of-memory crashes on the 8GB node, we are migrating from a local LGTM/Prometheus stack to a **Grafana Cloud** hybrid setup. 

- **Collector:** We will deploy the `grafana/k8s-monitoring` Helm chart, which runs **Grafana Alloy** as a lightweight unified telemetry collector.
- **GitOps:** The chart will be deployed via ArgoCD (added to `manifests/bootstrap/k8s-monitoring`).
- **Secret Management:** Connection details for Grafana Cloud (Prometheus URL, Loki URL, Access Token) will be stored in Azure Key Vault. A `SecretProviderClass` will use Workload Identity to fetch the token and sync it to a standard Kubernetes Secret that the Helm chart will consume.
- **Azure Monitor:** We will configure Grafana Cloud's Azure Monitor integration using the Service Principal created during the context phase.

## 3. Step-by-Step Implementation Tasks

### Task 1: Provision Grafana Cloud Secrets in Azure Key Vault
- **Action:** Obtain the Grafana Cloud Access Token, Prometheus Remote Write URL, and Loki Push URL. Store the Access Token securely in the Azure Key Vault. 
- **Commands:** 
  ```bash
  az keyvault secret set --vault-name <keyvault-name> --name "grafana-cloud-token" --value "<TOKEN>"
  ```
- **Validation:** 
  ```bash
  az keyvault secret show --vault-name <keyvault-name> --name "grafana-cloud-token"
  ```

### Task 2: Create SecretProviderClass for K8s-Monitoring
- **Action:** Create `manifests/bootstrap/k8s-monitoring/secret-provider-class.yaml`. This will use the Azure Key Vault CSI driver to pull `grafana-cloud-token` and sync it into a standard Kubernetes `Secret` named `grafana-cloud-credentials`.
- **Validation:** Apply the manifest and verify it exists via `kubectl get secretproviderclass -n monitoring`.

### Task 3: Deploy `k8s-monitoring` Helm Chart via ArgoCD
- **Action:** 
  1. Create `manifests/bootstrap/k8s-monitoring/values.yaml` containing the `externalServices` config (using the synced secret) and enabling `metrics.cost.enabled=true`.
  2. Create `manifests/bootstrap/k8s-monitoring.yaml` (ArgoCD Application) to deploy the `grafana/k8s-monitoring` Helm chart into a new `monitoring` namespace.
- **Commands:** Git commit and push.
- **Validation:** 
  ```bash
  kubectl get pods -n monitoring
  ```
  Check that the Grafana Alloy pods and `kube-state-metrics` are `Running`.

### Task 4: Verify Telemetry in Grafana Cloud
- **Action:** Log into the Grafana Cloud UI.
- **Validation:** 
  - Navigate to **Kubernetes Monitoring** and verify the `homelab-aks` cluster appears.
  - Navigate to **Logs > Explore** and verify container logs are arriving.
  - Configure the Azure Monitor integration in Grafana Cloud using the Service Principal credentials generated previously.

## 4. Rollback Plan

If the `k8s-monitoring` stack causes unforeseen issues (e.g., unexpectedly high memory usage):
1. **Remove the ArgoCD Application:** Delete the `k8s-monitoring.yaml` file from `manifests/bootstrap/` and commit/push. ArgoCD will prune the resources.
2. **Remove Key Vault Secret:** Delete the `grafana-cloud-token` from Azure Key Vault.
3. **Delete Namespace:** `kubectl delete namespace monitoring` to ensure all Alloy/Prometheus components are forcefully removed.
