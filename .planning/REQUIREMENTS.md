# Requirements — Homelab AKS Platform v1

## v1 Requirements

### INFRA — Azure Infrastructure

- [ ] **INFRA-01**: Platform engineer can create all Azure resources via `terraform apply` from a single root module with no manual portal steps
- [ ] **INFRA-02**: Azure VNet (10.0.0.0/16) with 4 subnets (system, user/spot, infra, ingress) provisioned by networking Terraform module
- [ ] **INFRA-03**: NSGs applied to all subnets with default-deny egress except required ports
- [ ] **INFRA-04**: Azure Budget Alert at 80% and 100% of monthly threshold fires email notification
- [ ] **INFRA-05**: All resources tagged with `environment=homelab`, `managed-by=terraform`

### CLUSTER — AKS Cluster

- [ ] **CLUSTER-01**: AKS cluster provisioned on Free tier (no control plane cost) in Central India, K8s `1.35.x`
- [ ] **CLUSTER-02**: System node pool: 1× Standard_D2as_v5, Zone 1, reserved instance, ephemeral OS disk (30GB), AzureLinux3 OS SKU
- [ ] **CLUSTER-03**: Spot node pool: Standard_D2as_v5, Zone 1+2, max 4 nodes, scale-to-zero, tainted `kubernetes.azure.com/scalesetpriority=spot:NoSchedule`
- [ ] **CLUSTER-04**: AKS CNI Overlay networking with pod CIDR `10.244.0.0/16` (pods do not consume VNet addresses)
- [ ] **CLUSTER-05**: AKS API server access restricted to authorized IP ranges (home IP) — no public exposure without IP allowlist
- [ ] **CLUSTER-06**: OIDC issuer enabled and Workload Identity enabled on cluster
- [ ] **CLUSTER-07**: Worker nodes have no public IPs — all traffic through Load Balancer only

### IDENTITY — Workload Identity Federation

- [ ] **IDENTITY-01**: No Service Principal client secrets exist anywhere in the codebase, cluster secrets, or Azure
- [ ] **IDENTITY-02**: Separate User-Assigned Managed Identity created per workload (argocd, jenkins, traefik)
- [ ] **IDENTITY-03**: Each Managed Identity has an `azurerm_federated_identity_credential` scoped to its exact ServiceAccount + namespace
- [ ] **IDENTITY-04**: Each Managed Identity has AKV RBAC role `Key Vault Secrets User` scoped only to required secrets (least privilege)
- [ ] **IDENTITY-05**: Azure Key Vault provisioned with Private Endpoint in infra-subnet (no public network access)
- [ ] **IDENTITY-06**: Secrets Store CSI Driver deployed; secrets mounted as files into pods (not stored as K8s Secret objects in etcd)

### GATEWAY — Traffic Routing (Traefik v3 + Gateway API)

- [ ] **GATEWAY-01**: Kubernetes Gateway API CRDs v1.5.1 (Standard channel) installed before Traefik
- [ ] **GATEWAY-02**: Traefik v3 installed with `kubernetesGateway.enabled=true` and `kubernetesIngress.enabled=false` — Gateway API is the exclusive routing API
- [ ] **GATEWAY-03**: Single Azure Public IP assigned to Traefik `LoadBalancer` service — all external traffic enters through one IP
- [ ] **GATEWAY-04**: HTTP (port 80) automatically redirected to HTTPS (port 443) at Gateway entrypoint level
- [ ] **GATEWAY-05**: `ReferenceGrant` deployed in every app namespace (argocd, jenkins, monitoring) to allow cross-namespace Gateway → Service binding
- [ ] **GATEWAY-06**: Cloudflare DNS `A` record for `*.yourdomain.com` points to Traefik public IP with proxy enabled

### TLS — Certificate Management

- [ ] **TLS-01**: cert-manager installed with `config.enableGatewayAPI=true` — watches Gateway API resources
- [ ] **TLS-02**: Cloudflare DNS-01 `ClusterIssuer` configured for Let's Encrypt production (wildcard cert support)
- [ ] **TLS-03**: Cloudflare API token stored in Azure Key Vault, surfaced via Secrets Store CSI — never hardcoded
- [ ] **TLS-04**: `Gateway` annotated with `cert-manager.io/cluster-issuer` — TLS certificates auto-issued and renewed without manual steps
- [ ] **TLS-05**: All platform UIs (ArgoCD, Jenkins, Grafana) accessible over HTTPS with valid cert — no browser warnings

### GITOPS — ArgoCD

- [ ] **GITOPS-01**: ArgoCD deployed via Helm bootstrap; subsequent ArgoCD upgrades managed by ArgoCD itself (self-managed)
- [ ] **GITOPS-02**: App-of-Apps pattern: single root `Application` in ArgoCD manages all platform Applications
- [ ] **GITOPS-03**: ArgoCD server accessible via `HTTPRoute` at `argocd.yourdomain.com` with TLS
- [ ] **GITOPS-04**: ArgoCD sync automatically reconciles any manual `kubectl` changes (auto-heal enabled on platform apps)
- [ ] **GITOPS-05**: ArgoCD ApplicationSet used for workload apps (directory-based generator on GitOps repo)

### CICD — Jenkins + JNLP

- [ ] **CICD-01**: Jenkins controller deployed on system node pool (not spot — needs stability)
- [ ] **CICD-02**: Jenkins accessible via `HTTPRoute` at `jenkins.yourdomain.com` with TLS
- [ ] **CICD-03**: Jenkins Kubernetes plugin configured to spin ephemeral JNLP agent pods on spot node pool
- [ ] **CICD-04**: JNLP agent pods tolerate spot node taint and are scheduled exclusively on spot pool
- [ ] **CICD-05**: JNLP agent pods scale to zero when no builds running (no idle compute cost)
- [ ] **CICD-06**: Jenkins controller has Workload Identity annotation — can authenticate to ACR/AKV without secrets
- [ ] **CICD-07**: At least one sample Jenkinsfile pipeline builds a Docker image and pushes to ACR successfully

### OBS — Observability

- [ ] **OBS-01**: kube-prometheus-stack deployed with 30s scrape interval and 7-day retention to fit within memory budget
- [ ] **OBS-02**: Grafana accessible via `HTTPRoute` at `grafana.yourdomain.com` with TLS
- [ ] **OBS-03**: Dashboards for: node resource usage, ArgoCD app health, Jenkins build metrics, Traefik request rates
- [ ] **OBS-04**: Alertmanager configured with at least one alert rule (node memory > 85%)
- [ ] **OBS-05**: Grafana dashboards stored as ConfigMaps in GitOps repo (persist across pod restarts)

### SEC — Security Hardening

- [ ] **SEC-01**: `NetworkPolicy` default-deny applied to all platform namespaces; explicit allow rules for required traffic only
- [ ] **SEC-02**: Pod security: no privileged containers, no `hostNetwork`, no `hostPID` on any platform component
- [ ] **SEC-03**: `LimitRange` applied to all namespaces — no unbounded pods can consume entire node
- [ ] **SEC-04**: `PodDisruptionBudget` on ArgoCD and Prometheus (prevents eviction from taking both replicas if scaled)

---

## v2 Requirements (Deferred — Not in v1 Scope)

- External DNS operator (auto-create Cloudflare DNS records from HTTPRoute hostnames) — useful but manual DNS setup is fine for homelab
- ArgoCD Image Updater (auto-PR image tag updates from ACR) — Jenkins git push pattern is simpler for v1
- Vertical Pod Autoscaler in recommendation mode — useful after a week of runtime data
- Multi-cluster ArgoCD (Hub + spoke) — single cluster is right for budget
- ACR geo-replication — not needed for single-region homelab
- Kyverno policy engine — NetworkPolicies + LimitRanges sufficient for v1

---

## Out of Scope

- Service Mesh (Istio/Linkerd) — memory overhead exceeds 8GB node budget; NetworkPolicies achieve isolation
- Azure Firewall egress — $300+/month, completely out of budget
- True 3-zone HA — needs 3 system nodes minimum, triples base cost
- VPN Gateway / Private Cluster — $27/month, replaced by authorized IP ranges
- HashiCorp Vault — AKV CSI + Workload Identity is zero-ops and functionally equivalent
- B-series nodes for system pool — Microsoft explicitly warns against for AKS system pools
- Full LGTM observability stack (Loki + Tempo + Mimir) — too heavy for 8GB node

---

## Traceability

*(Filled by roadmap — maps REQ-IDs to phases)*

| Requirement | Phase |
|---|---|
| INFRA-01 to INFRA-05 | Phase 1 |
| CLUSTER-01 to CLUSTER-07 | Phase 2 |
| IDENTITY-01 to IDENTITY-06 | Phase 3 + 4 |
| GATEWAY-01 to GATEWAY-06 | Phase 5 |
| TLS-01 to TLS-05 | Phase 5 |
| GITOPS-01 to GITOPS-05 | Phase 6 |
| CICD-01 to CICD-07 | Phase 7 |
| OBS-01 to OBS-05 | Phase 8 |
| SEC-01 to SEC-04 | Phase 9 |
