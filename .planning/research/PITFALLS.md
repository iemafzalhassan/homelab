# Pitfalls Research — AKS Homelab Platform (2026)

## P1 — Critical (Will break the cluster or blow the budget)

### P1.1: Running system components on the spot node pool
- **Warning sign**: CoreDNS, Traefik, ArgoCD pods scheduled on spot nodes
- **What happens**: Spot eviction kills the ingress controller → entire cluster unreachable
- **Prevention**: Always taint the spot pool with `kubernetes.azure.com/scalesetpriority=spot:NoSchedule` and add tolerations only to Jenkins JNLP pod templates
- **Phase**: Phase 2 (AKS cluster provisioning)
- **Terraform guard**: `node_taints = ["kubernetes.azure.com/scalesetpriority=spot:NoSchedule"]` on user pool

### P1.2: Forgetting `config.enableGatewayAPI=true` on cert-manager
- **Warning sign**: `Certificate` resources created but never reconciled; `CertificateRequest` stuck pending
- **What happens**: cert-manager silently ignores `Gateway` annotations — no TLS certs ever issued
- **Prevention**: Always install cert-manager with `--set config.enableGatewayAPI=true` Helm flag. Verify with `kubectl get clusterissuers` showing `Ready: True`
- **Phase**: Phase 5 (cert-manager setup)

### P1.3: Installing Traefik before Gateway API CRDs
- **Warning sign**: Traefik Helm install fails with "no kind `GatewayClass` registered"
- **What happens**: Traefik cannot register its GatewayClass controller → cluster has no ingress
- **Prevention**: Always `kubectl apply -f gateway-api-standard-install.yaml` BEFORE `helm install traefik`
- **Phase**: Phase 4/5 (Traefik + cert-manager)

### P1.4: Workload Identity federated credential namespace/SA mismatch
- **Warning sign**: Pods get `Azure identity token request failed` or `401 Unauthorized` from AKV
- **What happens**: Token exchange silently fails — wrong ServiceAccount name or namespace in federated credential `subject` field
- **Prevention**: Triple-check `subject = "system:serviceaccount:<namespace>:<sa-name>"` in `azurerm_federated_identity_credential`. Add a Terraform output that prints the expected subject for verification
- **Phase**: Phase 3 (identity) and Phase 6+ (when deploying workloads)

### P1.5: No Azure budget alert → unexpected bill
- **Warning sign**: No cost alert configured
- **What happens**: Forgotten resources (orphaned Public IP, undeleted node pool) run up charges silently
- **Prevention**: Create `azurerm_consumption_budget_resource_group` with 80% and 100% threshold alerts via Terraform. Set alert email in `terraform.tfvars`
- **Phase**: Phase 1 (infrastructure baseline)

---

## P2 — Major (Will cause frustrating debugging sessions)

### P2.1: No `ReferenceGrant` → HTTPRoute rejected silently
- **Warning sign**: `HTTPRoute` shows `Reason: NotAllowedByListeners` in status
- **What happens**: Gateway API security model requires explicit cross-namespace permission grants. Without `ReferenceGrant` in the app namespace, Traefik Gateway refuses to bind the route
- **Prevention**: Deploy `ReferenceGrant` in every app namespace (argocd, jenkins, monitoring) as part of the platform RBAC Application in ArgoCD
- **Phase**: Phase 5 (Gateway API routing)

### P2.2: DNS-01 Cloudflare challenge fails due to wrong API token permissions
- **Warning sign**: `CertificateRequest` stuck in `pending`; cert-manager logs show `DNS01 challenge failed: Cloudflare error 9109`
- **What happens**: The Cloudflare API token needs `Zone:DNS:Edit` and `Zone:Zone:Read` permissions scoped to the correct zone
- **Prevention**: Create a dedicated Cloudflare API token (not Global API Key) with minimum permissions. Store it in AKV, surface via Secrets Store CSI
- **Phase**: Phase 5 (cert-manager)

### P2.3: Missing `os_disk_type = "Ephemeral"` → slow node boot + unnecessary cost
- **Warning sign**: Node pool provisioning takes >5 minutes; Azure Storage costs appear on bill
- **What happens**: Default managed disk persists beyond node lifetime; ephemeral OS disk uses node's temp SSD (faster, free)
- **Prevention**: Set `os_disk_type = "Ephemeral"` on both system and spot pools in Terraform. Requires `os_disk_size_gb <= VM temp disk size` (D2as_v5 has 75GB temp disk — 30GB OS disk is safe)
- **Phase**: Phase 2 (AKS node pools)

### P2.4: Prometheus scrape interval too low → OOM on 8GB system node
- **Warning sign**: prometheus pod repeatedly OOMKilled
- **What happens**: kube-prometheus-stack default scrape interval (15s) + many targets → high memory cardinality
- **Prevention**: Set `prometheus.prometheusSpec.scrapeInterval: "30s"` and `retention: "7d"` in Helm values. Disable unnecessary exporters (e.g., etcd scraping — not accessible on AKS managed control plane)
- **Phase**: Phase 7 (observability)

### P2.5: ArgoCD sync fails on CRDs (race condition on first install)
- **Warning sign**: ArgoCD Application shows `ComparisonError: failed to get local live state` on first sync for Traefik or cert-manager
- **What happens**: CRDs deployed in the same sync wave as their dependents — CRDs not yet registered when ArgoCD tries to create CRD instances
- **Prevention**: Use ArgoCD sync waves: `argocd.argoproj.io/sync-wave: "-1"` on CRD-deploying Applications, `"0"` on everything else
- **Phase**: Phase 6 (ArgoCD App-of-Apps)

---

## P3 — Minor (Annoying but not blocking)

### P3.1: `kubectl` access broken after home IP changes
- **Warning sign**: `kubectl get nodes` returns `dial timeout` after IP change
- **What happens**: AKS authorized IP ranges allowlist is static — dynamic home IPs need manual update
- **Prevention**: Create a Makefile target `make update-ip` that runs `terraform apply -var="authorized_ip=$(curl -s ifconfig.me)/32"` to refresh quickly. Accept this as the tradeoff vs. VPN Gateway cost
- **Phase**: Ongoing operations

### P3.2: Jenkins JNLP agents fail to connect back to controller
- **Warning sign**: Agent pods start but immediately show `Connection refused` and terminate
- **What happens**: Jenkins controller URL in pod config doesn't match the actual service/ingress URL, or the controller's agent port (50000) is not accessible from the spot node pool
- **Prevention**: Set `jenkins.controller.jenkinsUrl` Helm value explicitly to the public HTTPS URL (Traefik HTTPRoute hostname). Ensure NetworkPolicy allows ingress on port 50000 from the user namespace
- **Phase**: Phase 7 (Jenkins)

### P3.3: Grafana dashboard not persisting after pod restart
- **Warning sign**: Custom dashboards disappear after Grafana restarts
- **What happens**: Grafana by default stores dashboards in SQLite in the pod — stateless restart loses them
- **Prevention**: Use `grafana.sidecar.dashboards.enabled: true` to load dashboards from ConfigMaps. Store dashboard JSONs in the GitOps repo, managed by ArgoCD
- **Phase**: Phase 7 (observability)

### P3.4: AKS node OS disk fills up with container image layers
- **Warning sign**: Node shows `DiskPressure` taint; pods evicted
- **What happens**: Ephemeral containers, build artifacts, and large images accumulate on the 30GB OS disk
- **Prevention**: Enable `imageGcHighThreshold: 80` and `imageGcLowThreshold: 70` in AKS node config. Use ACR with image retention policy to avoid pulling stale large images
- **Phase**: Phase 2 (AKS) + ongoing

---

## Quick Reference: Pre-Flight Checklist Before Each Phase

```bash
# Before Phase 2 (AKS):
az account show                              # Correct subscription?
az aks get-versions --location centralindia  # Latest K8s version?

# Before Phase 5 (Traefik + cert-manager):
kubectl api-resources | grep gateway         # Gateway API CRDs installed?
kubectl get ns traefik                       # traefik namespace exists?

# Before Phase 6 (ArgoCD):
kubectl get clusterissuers                   # cert-manager ready?
kubectl get gateway -n traefik               # Gateway bound + IP assigned?

# Before Phase 7 (Jenkins):
kubectl get sa -n jenkins                    # ServiceAccount with WI annotation?
kubectl describe secret -n jenkins           # No static secrets?
```

---
*Research Date: July 2026*
