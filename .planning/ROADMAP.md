# Roadmap: Homelab AKS Platform

## Overview

Build a production-grade, cost-optimized Kubernetes homelab on Azure (Central India) running 24x7 for ~$20-25/month. The journey starts with raw Azure infrastructure (VNet, AKS, identities) and progresses through each platform layer: ingress (Traefik v3 + Gateway API), certificates (cert-manager + Cloudflare), GitOps (ArgoCD), CI/CD (Jenkins + JNLP spot agents), and observability (kube-prometheus-stack). The final two phases harden security and validate the whole stack. Every phase builds on the last — nothing is deployed before its dependency is proven working.

## Phases

- [x] **Phase 1: Azure Foundation** — VNet, subnets, NSGs, resource groups, budget alert
- [x] **Phase 2: AKS Cluster** — Free tier AKS, system node pool, CNI Overlay, OIDC + Workload Identity
- [x] **Phase 3: Identities & Secrets** — Per-workload managed identities, federated credentials, AKV + Private Endpoint
- [x] **Phase 4: Secrets Distribution** — Secrets Store CSI Driver + AKV Provider; validate zero-credential secret mounting
- [x] **Phase 5: Gateway & TLS** — Gateway API CRDs → Traefik v3 → cert-manager → Cloudflare wildcard cert
- [x] **Phase 6: GitOps Bootstrap** — ArgoCD Helm install → App-of-Apps pattern → self-managed ArgoCD
- [x] **Phase 7: CI/CD** — Jenkins controller + spot node pool + JNLP ephemeral agents + sample pipeline
- [x] **Phase 8: Observability** — kube-prometheus-stack + Grafana dashboards + Alertmanager
- [ ] **Phase 9: Security Hardening** — NetworkPolicies, PodDisruptionBudgets, LimitRanges, pod security
- [ ] **Phase 10: Validation** — End-to-end smoke test, cost review, runbook documentation


## Phase Details

### Phase 1: Azure Foundation
**Goal**: All Azure baseline infrastructure exists and can be re-created from `terraform apply` alone
**Depends on**: Nothing
**Requirements**: INFRA-01, INFRA-02, INFRA-03, INFRA-04, INFRA-05
**Success Criteria** (what must be TRUE):
  1. `terraform apply` completes without errors from a clean state
  2. VNet `10.0.0.0/16` with 4 subnets exists and is visible in Azure Portal
  3. All subnets have NSGs attached with default-deny egress rules
  4. Azure Budget Alert at 80%+100% is configured and shows in Cost Management
  5. All resources have `environment=homelab` and `managed-by=terraform` tags
**Plans**: TBD

Plans:
- [ ] 01-01: Terraform root module, providers, backend config (`versions.tf`, `main.tf`, `variables.tf`, `outputs.tf`)
- [ ] 01-02: Networking module (VNet, subnets, NSGs, public IP placeholder)
- [ ] 01-03: Azure Budget Alert resource + `terraform.tfvars` template

---

### Phase 2: AKS Cluster
**Goal**: Running AKS cluster accessible via `kubectl` with Workload Identity ready for use
**Depends on**: Phase 1
**Requirements**: CLUSTER-01, CLUSTER-02, CLUSTER-03, CLUSTER-04, CLUSTER-05, CLUSTER-06, CLUSTER-07
**Success Criteria** (what must be TRUE):
  1. `kubectl get nodes` returns 1 Ready system node (Standard_D2as_v5, Zone 1)
  2. `kubectl version` shows K8s server 1.35.x
  3. AKS API server endpoint allows access only from authorized home IP — blocked from any other IP
  4. `kubectl get pods -A` shows all system pods Running (CoreDNS, kube-proxy, etc.)
  5. OIDC issuer URL is non-empty in `terraform output oidc_issuer_url`
**Plans**: TBD

Plans:
- [ ] 02-01: AKS Terraform module (cluster, system pool, CNI Overlay, OIDC, authorized IP ranges)
- [ ] 02-02: kubeconfig output + `make update-ip` Makefile target for dynamic IP refresh

---

### Phase 3: Identities & Secrets
**Goal**: Per-workload managed identities exist with federated credentials; AKV provisioned with Private Endpoint
**Depends on**: Phase 2 (needs OIDC issuer URL)
**Requirements**: IDENTITY-01, IDENTITY-02, IDENTITY-03, IDENTITY-04, IDENTITY-05
**Success Criteria** (what must be TRUE):
  1. Three User-Assigned Managed Identities exist in Azure (argocd, jenkins, traefik)
  2. Each has a federated credential with correct `subject = "system:serviceaccount:<ns>:<sa>"` 
  3. Key Vault exists with public network access disabled
  4. Key Vault has a Private Endpoint in the infra-subnet resolving via Private DNS Zone
  5. Each identity has `Key Vault Secrets User` role on AKV — no broader permissions
**Plans**: TBD

Plans:
- [ ] 03-01: Identity Terraform module (3× user-assigned managed identity + federated credentials)
- [ ] 03-02: Key Vault Terraform module (vault, private endpoint, private DNS zone, RBAC assignments)

---

### Phase 4: Secrets Distribution
**Goal**: Pods can mount secrets from AKV as files via Secrets Store CSI — no static K8s Secrets in etcd
**Depends on**: Phase 3
**Requirements**: IDENTITY-06
**Success Criteria** (what must be TRUE):
  1. Secrets Store CSI Driver and AKV provider DaemonSet pods are Running on system node
  2. A test pod mounting a `SecretProviderClass` that reads from AKV starts successfully
  3. `kubectl exec` into test pod shows secret file content is correct
  4. `kubectl get secrets -A` shows no application secrets in etcd — only CSI-projected tokens
**Plans**: TBD

Plans:
- [ ] 04-01: Secrets Store CSI + AKV provider Helm install (via ArgoCD later — manual bootstrap for now)
- [ ] 04-02: Test `SecretProviderClass` + test pod YAML; validate end-to-end secret mounting

---

### Phase 5: Gateway & TLS
**Goal**: External HTTPS traffic reaches the cluster and all platform services get valid wildcard TLS certs automatically
**Depends on**: Phase 4
**Requirements**: GATEWAY-01, GATEWAY-02, GATEWAY-03, GATEWAY-04, GATEWAY-05, GATEWAY-06, TLS-01, TLS-02, TLS-03, TLS-04, TLS-05
**Success Criteria** (what must be TRUE):
  1. `kubectl api-resources | grep gateway` shows `HTTPRoute`, `Gateway`, `GatewayClass` registered
  2. Traefik `LoadBalancer` service has an assigned Azure Public IP
  3. Cloudflare `*.yourdomain.com` A record points to Traefik public IP with proxy enabled
  4. `curl -I https://traefik.yourdomain.com` returns HTTP 200 with a valid Let's Encrypt cert (no browser warnings)
  5. HTTP → HTTPS redirect works: `curl -I http://traefik.yourdomain.com` returns 301/302
  6. `kubectl get certificaterequest -A` shows cert requests Approved and Issued
**Plans**: TBD

Plans:
- [ ] 05-01: Gateway API CRDs v1.5.1 install + Traefik v3 Helm values (Gateway API first-class)
- [ ] 05-02: cert-manager Helm install with `enableGatewayAPI=true` + Cloudflare DNS-01 ClusterIssuer
- [ ] 05-03: `Gateway` resource with cert-manager annotation + wildcard `HTTPRoute` for Traefik dashboard
- [ ] 05-04: `ReferenceGrant` per-namespace template + validate cross-namespace routing works

---

### Phase 6: GitOps Bootstrap
**Goal**: ArgoCD manages its own lifecycle and all platform Applications; App-of-Apps pattern operational
**Depends on**: Phase 5 (needs Traefik HTTPRoute for ArgoCD UI access)
**Requirements**: GITOPS-01, GITOPS-02, GITOPS-03, GITOPS-04, GITOPS-05
**Success Criteria** (what must be TRUE):
  1. `https://argocd.yourdomain.com` loads ArgoCD UI with valid TLS
  2. ArgoCD root Application shows all child Applications as Healthy + Synced
  3. Manually deleting a deployed ConfigMap causes ArgoCD to re-create it within 3 minutes (auto-heal works)
  4. GitOps repo has `apps/` directory with at least 3 platform Application definitions (traefik, cert-manager, monitoring)
  5. ArgoCD itself is managed by an ArgoCD Application (self-managed pattern)
**Plans**: TBD

Plans:
- [ ] 06-01: ArgoCD Helm bootstrap (values file, ServiceAccount with WI annotation, HTTPRoute for UI)
- [ ] 06-02: App-of-Apps root Application YAML + GitOps repo `apps/` directory structure
- [ ] 06-03: Self-managed ArgoCD Application + sync wave configuration (CRDs wave -1, apps wave 0)

---

### Phase 7: CI/CD
**Goal**: Jenkins can build Docker images using ephemeral JNLP pods on spot nodes, authenticated via Workload Identity
**Depends on**: Phase 6 (ArgoCD deploys Jenkins)
**Requirements**: CICD-01, CICD-02, CICD-03, CICD-04, CICD-05, CICD-06, CICD-07
**Success Criteria** (what must be TRUE):
  1. `https://jenkins.yourdomain.com` loads Jenkins UI with valid TLS
  2. Running a Jenkins pipeline creates an agent pod on the spot node pool and deletes it after the build completes
  3. `kubectl get pods -n jenkins` shows 0 agent pods when no builds running (scale-to-zero confirmed)
  4. Jenkins agent pod can authenticate to ACR without any stored credentials (Workload Identity)
  5. A sample pipeline that builds a Docker image and pushes to ACR completes successfully
  6. Spot node pool exists and shows 1-2 nodes during build, 0 nodes at idle
**Plans**: TBD

Plans:
- [ ] 07-01: Spot node pool Terraform module + validation (taint, zone spread, scale 0→4)
- [ ] 07-02: Jenkins ArgoCD Application (Helm values: JNLP config, spot toleration, WI ServiceAccount, HTTPRoute)
- [ ] 07-03: Sample Jenkinsfile (Docker build → ACR push) + end-to-end pipeline run

---

### Phase 8: Observability
**Goal**: Prometheus scrapes all platform components; Grafana shows actionable dashboards; Alertmanager fires on critical conditions
**Depends on**: Phase 6 (ArgoCD deploys monitoring stack)
**Requirements**: OBS-01, OBS-02, OBS-03, OBS-04, OBS-05
**Success Criteria** (what must be TRUE):
  1. `https://grafana.yourdomain.com` loads with valid TLS
  2. Grafana dashboard shows live node CPU/memory usage for the system node
  3. Grafana dashboard shows ArgoCD application health and sync status
  4. Grafana dashboard shows Traefik request rate and 5xx error rate
  5. Alertmanager has at least one firing-capable rule (e.g., node memory > 85%)
  6. Grafana dashboards survive a pod restart (persisted via ConfigMap)
**Plans**: TBD

Plans:
- [ ] 08-01: kube-prometheus-stack ArgoCD Application (Helm values: 30s scrape, 7d retention, HTTPRoute for Grafana)
- [ ] 08-02: Grafana dashboard ConfigMaps (node, ArgoCD, Traefik, Jenkins) stored in GitOps repo
- [ ] 08-03: Alertmanager rule for node memory + email/webhook notification config

---

### Phase 9: Security Hardening
**Goal**: All namespaces have default-deny NetworkPolicies; no unbounded resource consumption; platform components survive node disruption
**Depends on**: Phase 8 (all components deployed — can now apply policies without breaking things)
**Requirements**: SEC-01, SEC-02, SEC-03, SEC-04
**Success Criteria** (what must be TRUE):
  1. `kubectl get networkpolicies -A` shows a default-deny policy in every platform namespace
  2. Cross-namespace traffic that is expected (Prometheus scraping Jenkins) still works after NetworkPolicies applied
  3. `kubectl describe pod <any-platform-pod>` shows no privileged containers, no hostNetwork, no hostPID
  4. `kubectl get limitrange -A` shows LimitRange in every namespace with CPU/memory limits
  5. `kubectl get poddisruptionbudget -A` shows PDB on ArgoCD and Prometheus
**Plans**: TBD

Plans:
- [ ] 09-01: NetworkPolicy manifests (default-deny + explicit allow rules per namespace) in GitOps repo
- [ ] 09-02: LimitRange manifests per namespace + PodDisruptionBudgets for ArgoCD and Prometheus
- [ ] 09-03: Pod security audit (no privileged containers, no host-level access) — fix any violations found

---

### Phase 10: Validation
**Goal**: End-to-end smoke test passes; cost is within budget; runbook exists for common ops tasks
**Depends on**: Phase 9
**Requirements**: All v1 requirements implicitly validated
**Success Criteria** (what must be TRUE):
  1. Full end-to-end test: push code → Jenkins build → Docker push → ArgoCD sync → pod updated
  2. Azure Cost Management shows actual spend within $25/month budget
  3. Deleting and re-applying the root ArgoCD Application re-creates the entire platform within 15 minutes
  4. `make update-ip` successfully refreshes the authorized IP and `kubectl get nodes` works from new IP
  5. Runbook document exists covering: IP rotation, cert renewal verification, spot node pool scale-up, Jenkins agent debugging
**Plans**: TBD

Plans:
- [ ] 10-01: End-to-end smoke test script + manual test run with documented results
- [ ] 10-02: Cost review (Azure Cost Management) + budget alert verification
- [ ] 10-03: RUNBOOK.md with operational procedures

---

## Progress

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Azure Foundation | 0/3 | Not started | - |
| 2. AKS Cluster | 0/2 | Not started | - |
| 3. Identities & Secrets | 0/2 | Not started | - |
| 4. Secrets Distribution | 0/2 | Not started | - |
| 5. Gateway & TLS | 0/4 | Not started | - |
| 6. GitOps Bootstrap | 0/3 | Not started | - |
| 7. CI/CD | 0/3 | Not started | - |
| 8. Observability | 0/3 | Not started | - |
| 9. Security Hardening | 0/3 | Not started | - |
| 10. Validation | 0/3 | Not started | - |

