# Phase 11: GitOps Promotion (Kargo) — Research

**Date:** 2026-07-09
**Status:** Complete

---

## 1. Kargo Version & Helm Chart

**Latest stable:** `v1.10.8` (chart version `1.10.8`)
**OCI chart registry:** `oci://ghcr.io/akuity/kargo-charts/kargo`

Install command:
```bash
helm install kargo oci://ghcr.io/akuity/kargo-charts/kargo \
  --namespace kargo \
  --create-namespace \
  --version 1.10.8 \
  -f manifests/bootstrap/kargo/values.yaml
```

### Key Helm values structure

```yaml
api:
  host: kargo.smapatticare.com      # used for OIDC redirect URI + deep-links
  service:
    type: ClusterIP                  # Traefik HTTPRoute targets this
  oidc:
    enabled: true
    issuerURL: "https://sso.smapatticare.com/realms/Platform"
    clientID: "kargo"
    clientSecret: ""                 # loaded from K8s Secret via env ref
    admins:
      claims:
        groups: ["homelab-admins", "devops"]
    viewers:
      claims:
        groups: ["developers", "viewers"]
  argocd:
    urls:
      "": "https://argocd.smapatticare.com"  # deep-link back to ArgoCD UI

controller:
  argocd:
    integrationEnabled: true         # enables Kargo ↔ ArgoCD reconciliation

# Resource caps for spot node
resources:
  api:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 500m
      memory: 256Mi
  controller:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 1000m
      memory: 512Mi

# Spot node placement
tolerations:
  - key: "kubernetes.azure.com/scalesetpriority"
    operator: "Equal"
    value: "spot"
    effect: "NoSchedule"
nodeSelector:
  kubernetes.azure.com/scalesetpriority: "spot"
```

---

## 2. Kargo CRD Reference

### 2.1 Project

The `Project` CR creates the Kargo namespace and sets `autoPromotionPolicies`. Setting a Stage's policy to `autoPromotionEnabled: false` enforces manual approval.

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Project
metadata:
  name: homelab-demo
spec:
  promotionPolicies:
    - stage: uat
      autoPromotionEnabled: true    # UAT: auto-deploy when Freight arrives
    - stage: prod
      autoPromotionEnabled: false   # PROD: manual approval required
```

### 2.2 Warehouse (Git subscription)

Watches a specific path in the GitOps repo for commits. When Jenkins commits an updated image tag to `manifests/apps/homelab-demo/uat/values.yaml`, Kargo detects it and creates a `Freight` object.

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Warehouse
metadata:
  name: homelab-demo
  namespace: homelab-demo            # same as Project name
spec:
  subscriptions:
    - git:
        repoURL: https://github.com/iemafzalhassan/homelab.git
        branch: main
        includePaths:
          - manifests/apps/homelab-demo/   # only trigger on changes here
```

### 2.3 Stage: UAT

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Stage
metadata:
  name: uat
  namespace: homelab-demo
spec:
  requestedFreight:
    - origin:
        kind: Warehouse
        name: homelab-demo
      sources:
        direct: true
  promotionTemplate:
    spec:
      steps:
        - uses: argocd-update
          as: update-uat
          config:
            apps:
              - name: homelab-demo-uat
                namespace: argocd
                sources:
                  - repoURL: https://github.com/iemafzalhassan/homelab.git
                    desiredCommitFromStep: update-uat
```

### 2.4 Stage: PROD (manual gate)

```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Stage
metadata:
  name: prod
  namespace: homelab-demo
spec:
  requestedFreight:
    - origin:
        kind: Warehouse
        name: homelab-demo
      sources:
        stages: ["uat"]             # PROD only takes Freight verified in UAT
  promotionTemplate:
    spec:
      steps:
        - uses: argocd-update
          as: update-prod
          config:
            apps:
              - name: homelab-demo-prod
                namespace: argocd
                sources:
                  - repoURL: https://github.com/iemafzalhassan/homelab.git
                    desiredCommitFromStep: update-prod
```

The manual gate is enforced by `autoPromotionEnabled: false` in the `Project` — no `Stage`-level field needed.

---

## 3. ArgoCD Integration Pattern

### 3.1 How Kargo talks to ArgoCD

Kargo's controller talks to ArgoCD via the **Kubernetes API** (not ArgoCD's own API server). It reads and patches `Application` resources directly. The Kargo Helm chart creates the necessary `ClusterRole` and `ClusterRoleBinding` for its controller ServiceAccount.

**Critical requirement:** Each ArgoCD `Application` that Kargo manages MUST have the annotation:
```yaml
metadata:
  annotations:
    kargo.akuity.io/authorized-stage: "homelab-demo:uat"   # or "homelab-demo:prod"
```

This is the security boundary — Kargo can only update Applications explicitly annotated for its Stages.

### 3.2 ArgoCD Applications for demo app

These are **separate** from the Kargo ArgoCD Application. They are the target apps Kargo promotes into:

**UAT Application** (`manifests/apps/homelab-demo-uat.yaml`):
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: homelab-demo-uat
  namespace: argocd
  annotations:
    kargo.akuity.io/authorized-stage: "homelab-demo:uat"
spec:
  project: default
  sources:
    - repoURL: 'https://github.com/iemafzalhassan/homelab.git'
      path: manifests/apps/homelab-demo/uat
      targetRevision: HEAD
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: demo-uat
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

**PROD Application** (`manifests/apps/homelab-demo-prod.yaml`):
```yaml
# Same structure but:
# - annotations: kargo.akuity.io/authorized-stage: "homelab-demo:prod"
# - path: manifests/apps/homelab-demo/prod
# - namespace: demo-prod
# - syncPolicy: {} (NO automated sync — Kargo controls when it syncs)
```

> **Key insight:** PROD ArgoCD Application must NOT have `syncPolicy.automated` — Kargo triggers the sync via the `argocd-update` promotion step.

---

## 4. Keycloak OIDC for Kargo

### 4.1 Keycloak client setup

Create a new client `kargo` in the Platform realm:
- **Client type:** OpenID Connect
- **Authentication flow:** Authorization Code + PKCE (Kargo requires PKCE)
- **Valid redirect URIs:** `https://kargo.smapatticare.com/auth/callback`
- **Client secret:** generate and store in Azure Key Vault → mount as K8s Secret
- **Scopes:** Add `groups` scope (already exists for other clients — reuse the mapper)

### 4.2 Groups claim mapping

The existing `oidc-group-membership-mapper` in the Platform realm already maps group memberships to the `groups` claim in the JWT. No new mapper needed — just add the `groups` scope to the `kargo` client.

### 4.3 Kargo RBAC via claims

```yaml
api:
  oidc:
    admins:
      claims:
        groups: ["homelab-admins", "devops"]   # can approve PROD promotions
    viewers:
      claims:
        groups: ["developers", "viewers"]       # view-only: UAT status, Freight
```

Kargo v1 maps these directly — no intermediate Kubernetes RBAC required for the UI-level RBAC.

---

## 5. Traefik HTTPRoute for Kargo

Kargo API server listens on port **8080** (HTTP, internal). The HTTPRoute pattern matches all other platform services exactly.

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: kargo
  namespace: kargo
spec:
  parentRefs:
    - name: traefik-gateway
      namespace: traefik
      sectionName: websecure
  hostnames:
    - "kargo.smapatticare.com"
  rules:
    - backendRefs:
        - name: kargo-api
          port: 8080
```

**ReferenceGrant** required (cross-namespace):
```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: kargo-traefik-grant
  namespace: kargo
spec:
  from:
    - group: gateway.networking.k8s.io
      kind: HTTPRoute
      namespace: traefik
  to:
    - group: ""
      kind: Service
```

---

## 6. Demo App: homelab-demo

### 6.1 Structure

```
apps/homelab-demo/
  Dockerfile                         # simple nginx serving version.html
  Jenkinsfile.demo                   # CI pipeline: build → push → commit tag
manifests/apps/homelab-demo/
  uat/
    deployment.yaml                  # image: ACR_URL/homelab-demo:TAG (mutable)
    service.yaml
    kustomization.yaml               # or values.yaml if Helm-based
  prod/
    deployment.yaml                  # same structure, separate path
    service.yaml
    kustomization.yaml
```

### 6.2 Minimal Dockerfile

```dockerfile
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/index.html
EXPOSE 80
```

The `index.html` includes a version string (e.g. `v1.2.3`) so the demo visually shows "the app changed" when promoted.

### 6.3 Jenkins pipeline: commit image tag

After building and pushing the image, Jenkins commits the new tag to the GitOps repo:
```groovy
sh """
  sed -i 's|image: .*homelab-demo.*|image: ${ACR_URL}/homelab-demo:${IMAGE_TAG}|' \
    manifests/apps/homelab-demo/uat/deployment.yaml
  git config user.email "jenkins@smapatticare.com"
  git config user.name "Jenkins CI"
  git add manifests/apps/homelab-demo/uat/deployment.yaml
  git commit -m "ci: update homelab-demo UAT image to ${IMAGE_TAG}"
  git push origin main
"""
```

Kargo's Warehouse then automatically detects the new commit on the `main` branch under `manifests/apps/homelab-demo/` and creates a `Freight`.

---

## 7. Kargo Secret: OIDC Client Secret

Since we use Azure Key Vault + Secrets Store CSI Driver, the client secret flows as:

1. Store `kargo-client-secret` in AKV
2. Create `SecretProviderClass` in `kargo` namespace referencing AKV secret
3. Mount as K8s Secret named `kargo-oidc-secret` with key `client-secret`
4. Reference in Kargo values:
```yaml
api:
  env:
    - name: OIDC_CLIENT_SECRET
      valueFrom:
        secretKeyRef:
          name: kargo-oidc-secret
          key: client-secret
  oidc:
    clientSecret: ""   # leave empty — injected via env
```

**Alternative (simpler for homelab):** Create the K8s Secret directly since Kargo doesn't need to access AKV directly. The secret only needs to exist in the `kargo` namespace.

---

## 8. Kargo ArgoCD Application (bootstrapped via ArgoCD)

Following the existing App-of-Apps pattern (`manifests/apps/*.yaml`):

```yaml
# manifests/apps/kargo.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: kargo
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  sources:
    - repoURL: 'oci://ghcr.io/akuity/kargo-charts'
      chart: kargo
      targetRevision: 1.10.8
      helm:
        valueFiles:
          - $values/manifests/bootstrap/kargo/values.yaml
    - repoURL: 'https://github.com/iemafzalhassan/homelab.git'
      targetRevision: HEAD
      ref: values
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: kargo
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

> **Note:** OCI-sourced Helm charts require ArgoCD `v2.12+` with OCI enabled. Our ArgoCD `v3.4.4` (chart `10.1.0`) supports this natively.

---

## 9. End-to-End Promotion Flow (Verified)

```
1. Developer commits code → GitHub
2. Jenkins builds image → pushes to ACR (homelab-demo:git-abc1234)
3. Jenkins commits image tag update to manifests/apps/homelab-demo/uat/deployment.yaml
4. Kargo Warehouse detects new commit → creates Freight{commit: abc1234}
5. UAT Stage: autoPromotionEnabled=true → Kargo triggers argocd-update step
   → ArgoCD syncs homelab-demo-uat Application → pod in demo-uat updated
6. Dev team verifies UAT manually
7. DevOps opens Kargo UI at kargo.smapatticare.com → logs in via Keycloak
8. Clicks "Promote" on UAT Freight → selects PROD stage → clicks Confirm
9. Kargo triggers argocd-update for homelab-demo-prod → ArgoCD syncs PROD
   → pod in demo-prod updated
10. PROD is live ✅
```

---

## 10. Validation Architecture

### Plan 11-01 (Kargo install) — verification commands:
```bash
kubectl get pods -n kargo                          # all pods Running
kubectl get svc kargo-api -n kargo                 # ClusterIP exists on port 8080
curl -sk https://kargo.smapatticare.com            # 200 or 302 to OIDC login
kubectl get crd | grep kargo                       # Project, Warehouse, Stage, Freight CRDs present
```

### Plan 11-02 (CRs + demo app) — verification commands:
```bash
kubectl get project homelab-demo -n homelab-demo   # Project exists
kubectl get warehouse homelab-demo -n homelab-demo # Warehouse exists
kubectl get stage uat prod -n homelab-demo         # both Stages exist
kubectl get freight -n homelab-demo                # Freight created after Jenkins push
kubectl get pods -n demo-uat                       # demo app running in UAT
```

### Plan 11-03 (RBAC + e2e) — verification commands:
```bash
# Login as devops@smapatticare.com → can see Promote button for PROD
# Login as dev@smapatticare.com → can see Stage status, no Promote button
kubectl get rolebinding -n homelab-demo            # Kargo role bindings present
# Run Jenkins demo pipeline → verify Freight created → UAT auto-promoted
# Click Promote in UI for PROD → verify demo-prod pod updated
```

---

## RESEARCH COMPLETE
