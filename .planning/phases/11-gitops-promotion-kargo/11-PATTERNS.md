# Phase 11: GitOps Promotion (Kargo) — Pattern Mapping

**Generated:** 2026-07-09  
**Phase:** 11 — gitops-promotion-kargo  
**Status:** Complete

---

## Files to Create / Modify

| File | Role | Category |
|------|------|----------|
| `manifests/apps/kargo.yaml` | ArgoCD Application (Helm install of Kargo) | App-of-Apps entry |
| `manifests/apps/kargo-extras.yaml` | ArgoCD Application (HTTPRoute + ReferenceGrant) | App-of-Apps extras |
| `manifests/bootstrap/kargo/values.yaml` | Helm values for Kargo chart | Bootstrap config |
| `manifests/bootstrap/kargo/kargo-httproute.yaml` | Traefik HTTPRoute for Kargo UI | Networking |
| `manifests/bootstrap/kargo/kargo-referencegrant.yaml` | ReferenceGrant for cross-ns routing | Networking |
| `manifests/apps/homelab-demo-uat.yaml` | ArgoCD Application for UAT demo app | App-of-Apps entry |
| `manifests/apps/homelab-demo-prod.yaml` | ArgoCD Application for PROD demo app | App-of-Apps entry |
| `manifests/apps/homelab-demo/uat/deployment.yaml` | Kubernetes Deployment — demo app UAT | App overlay |
| `manifests/apps/homelab-demo/uat/service.yaml` | Kubernetes Service — demo app UAT | App overlay |
| `manifests/apps/homelab-demo/uat/kustomization.yaml` | Kustomize root for UAT overlay | App overlay |
| `manifests/apps/homelab-demo/prod/deployment.yaml` | Kubernetes Deployment — demo app PROD | App overlay |
| `manifests/apps/homelab-demo/prod/service.yaml` | Kubernetes Service — demo app PROD | App overlay |
| `manifests/apps/homelab-demo/prod/kustomization.yaml` | Kustomize root for PROD overlay | App overlay |
| `kargo/project.yaml` | Kargo Project CR | Kargo CRs |
| `kargo/warehouse.yaml` | Kargo Warehouse CR | Kargo CRs |
| `kargo/stage-uat.yaml` | Kargo Stage CR (UAT) | Kargo CRs |
| `kargo/stage-prod.yaml` | Kargo Stage CR (PROD) | Kargo CRs |
| `apps/homelab-demo/Dockerfile` | Container image build | CI artifact |
| `apps/homelab-demo/index.html` | Static page served by nginx | CI artifact |
| `apps/homelab-demo/Jenkinsfile.demo` | Jenkins pipeline for demo CI | CI pipeline |
| `docs/ADR-001-platform-auth-rbac.md` | **MODIFY**: add Kargo client to OIDC client table | Documentation |

---

## Pattern 1 — ArgoCD Application (Helm, OCI Chart)

### Analog: `manifests/apps/argocd.yaml` and `manifests/apps/jenkins.yaml`

All platform tools follow the same two-source App-of-Apps structure: one source for the Helm
chart, one source (`ref: values`) pointing to the GitOps repo for value files.

**Concrete excerpt from `manifests/apps/argocd.yaml`:**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: argocd
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  sources:
    - repoURL: 'https://argoproj.github.io/argo-helm'
      chart: argo-cd
      targetRevision: 10.1.2
      helm:
        valueFiles:
          - $values/manifests/bootstrap/argocd/values.yaml
    - repoURL: 'https://github.com/iemafzalhassan/homelab.git'
      targetRevision: HEAD
      ref: values
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: argocd
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

**Apply to `manifests/apps/kargo.yaml`:**
- Replace `repoURL` with `oci://ghcr.io/akuity/kargo-charts` and `chart: kargo`
- `targetRevision: 1.10.8`
- `namespace: kargo`
- Value file path: `$values/manifests/bootstrap/kargo/values.yaml`
- Keep `finalizers`, `automated.prune`, `selfHeal`, `CreateNamespace=true` identical

> NOTE: OCI chart sources use the `oci://` scheme, not HTTPS.
> ArgoCD v3.4.4 supports OCI natively — no extra configuration needed.

---

## Pattern 2 — ArgoCD "Extras" Application (HTTPRoute sidecar)

### Analog: `manifests/apps/argocd-extras.yaml` and `manifests/apps/jenkins-extras.yaml`

Each platform tool has a companion `-extras` Application that deploys namespace-scoped resources
(HTTPRoute, ServiceMonitors) which cannot be bundled with the Helm chart.

**Concrete excerpt from `manifests/apps/jenkins-extras.yaml`:**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: jenkins-extras
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: 'https://github.com/iemafzalhassan/homelab.git'
    targetRevision: HEAD
    path: manifests/bootstrap/jenkins
    directory:
      include: 'jenkins-httproute.yaml'
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: jenkins
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

**Apply to `manifests/apps/kargo-extras.yaml`:**
- `name: kargo-extras`
- `path: manifests/bootstrap/kargo`
- `include: '{kargo-httproute.yaml,kargo-referencegrant.yaml}'`
- `namespace: kargo`

---

## Pattern 3 — Helm `values.yaml` (OIDC + spot tolerations)

### Analog: `manifests/bootstrap/jenkins/values.yaml`

The Jenkins values demonstrate the canonical pattern for all OIDC-integrated tools:
1. OIDC client secret loaded from a K8s Secret via `secretKeyRef`
2. Spot node placement via `nodeSelector` + `tolerations`
3. OIDC issuerURL pointing to `https://sso.smapatticare.com/realms/Platform`

**Concrete spot toleration block (jenkins/values.yaml L18-24):**
```yaml
nodeSelector:
  kubernetes.azure.com/scalesetpriority: "spot"
tolerations:
  - key: "kubernetes.azure.com/scalesetpriority"
    operator: "Equal"
    value: "spot"
    effect: "NoSchedule"
```

**Concrete OIDC secret injection (jenkins/values.yaml L8-13):**
```yaml
containerEnv:
  - name: OIDC_CLIENT_SECRET
    valueFrom:
      secretKeyRef:
        name: jenkins-oidc-secret
        key: client-secret
```

**Concrete OIDC issuerURL pattern (argocd/values.yaml L14-19):**
```yaml
oidc.config: |
  name: Keycloak
  issuer: https://sso.smapatticare.com/realms/Platform
  clientID: argocd
  clientSecret: $argocd-oidc-secret:client-secret
  requestedScopes: ["openid", "profile", "email", "groups"]
```

**Apply to `manifests/bootstrap/kargo/values.yaml`:**
- Copy spot toleration block verbatim (same key/operator/value/effect across ALL platform tools)
- Secret name: `kargo-oidc-secret`, key: `client-secret` (env var: `OIDC_CLIENT_SECRET`)
- OIDC issuerURL: `https://sso.smapatticare.com/realms/Platform`
- OIDC clientID: `kargo`
- `api.host: kargo.smapatticare.com`
- RBAC groups: `admins → ["homelab-admins", "devops"]`, `viewers → ["developers", "viewers"]`

---

## Pattern 4 — Traefik HTTPRoute

### Analog: `manifests/bootstrap/jenkins/jenkins-httproute.yaml`

Jenkins is the closest analog: its service also exposes port `8080` and the HTTPRoute lives in
the service's own namespace.

**Concrete excerpt from `manifests/bootstrap/jenkins/jenkins-httproute.yaml`:**
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: jenkins
  namespace: jenkins
spec:
  parentRefs:
    - name: traefik-gateway
      namespace: traefik
      sectionName: web
    - name: traefik-gateway
      namespace: traefik
      sectionName: jenkins
  hostnames:
    - "jenkins.smapatticare.com"
  rules:
    - backendRefs:
        - name: jenkins
          port: 8080
```

**ArgoCD analog for comparison (`argocd-httproute.yaml`):**
```yaml
parentRefs:
  - name: traefik-gateway
    namespace: traefik
    sectionName: web
  - name: traefik-gateway
    namespace: traefik
    sectionName: websecure
rules:
  - matches:
      - path:
          type: PathPrefix
          value: /
    backendRefs:
      - name: argocd-server
        port: 80
```

**Apply to `manifests/bootstrap/kargo/kargo-httproute.yaml`:**
- `name: kargo`, `namespace: kargo`
- Single `parentRef` with `sectionName: websecure` only (Kargo UI is HTTPS-only)
- `hostname: kargo.smapatticare.com`
- `backendRefs[0].name: kargo-api` (Kargo Helm chart names the service `kargo-api`)
- `port: 8080`
- No `matches` block needed — bare backendRefs like Jenkins

---

## Pattern 5 — ReferenceGrant (cross-namespace Gateway→Service)

### Analog: Architecture pattern from `.planning/research/STACK.md` L218-234

No deployed ReferenceGrant YAML exists yet in `manifests/` — they are referenced in planning
docs as a pending item. The canonical pattern from research:

**Concrete excerpt from `.planning/research/STACK.md` L218-234:**
```yaml
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: allow-traefik-gateway
  namespace: argocd   # repeat per namespace
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: Gateway
    namespace: traefik
  to:
  - group: ""
    kind: Service
```

**Apply to `manifests/bootstrap/kargo/kargo-referencegrant.yaml`:**
- `name: kargo-traefik-grant`
- `namespace: kargo`

> DISCREPANCY ALERT: RESEARCH.md (11-RESEARCH.md L292-305) uses `from[0].kind: HTTPRoute`
> while STACK.md uses `from[0].kind: Gateway`. Per Gateway API spec, the correct form for
> allowing an HTTPRoute in namespace X to reference a Service in namespace Y is `kind: HTTPRoute`.
> Use the RESEARCH.md form — `kind: HTTPRoute`, `namespace: kargo` (where the HTTPRoute lives).

```yaml
# Correct form for kargo-referencegrant.yaml
spec:
  from:
    - group: gateway.networking.k8s.io
      kind: HTTPRoute
      namespace: kargo
  to:
    - group: ""
      kind: Service
```

---

## Pattern 6 — ArgoCD Application for a kustomize-managed app

### Analog: `manifests/apps/gpt-researcher/` subdirectory layout

The gpt-researcher app uses a flat kustomize layout within `manifests/apps/<appname>/`. Kargo-
managed apps follow the same structure but split into `uat/` and `prod/` subdirectories.

**Existing kustomize structure (`gpt-researcher/kustomization.yaml`):**
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - namespace.yaml
  - backend.yaml
  - frontend.yaml
  - gateway.yaml
```

**Apply to `manifests/apps/homelab-demo/uat/kustomization.yaml`:**
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - deployment.yaml
  - service.yaml
```

**ArgoCD Application for UAT (`manifests/apps/homelab-demo-uat.yaml`):**
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: homelab-demo-uat
  namespace: argocd
  annotations:
    kargo.akuity.io/authorized-stage: "homelab-demo:uat"   # CRITICAL for Kargo
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: 'https://github.com/iemafzalhassan/homelab.git'
    targetRevision: HEAD
    path: manifests/apps/homelab-demo/uat
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

**PROD Application differs in exactly these fields:**
1. `name: homelab-demo-prod`
2. `annotations: kargo.akuity.io/authorized-stage: "homelab-demo:prod"`
3. `path: manifests/apps/homelab-demo/prod`
4. `destination.namespace: demo-prod`
5. `syncPolicy: {}` — NO automated sync; Kargo controls PROD syncs

> CRITICAL: The `kargo.akuity.io/authorized-stage` annotation is the security boundary between
> Kargo and ArgoCD. Without it, Kargo's controller cannot patch the Application.

> CRITICAL: PROD ArgoCD Application must have NO `syncPolicy.automated` block.
> Without this constraint, ArgoCD and Kargo will race to control PROD.

---

## Pattern 7 — Kargo CRs (Project, Warehouse, Stage)

### Analog: None in codebase — new pattern

Conventions to follow from the rest of the platform:

**Namespace/label conventions:**
- Namespace = Project name: `homelab-demo`
- Kargo Project CR auto-creates its namespace
- `metadata.namespace: homelab-demo` on all Warehouse and Stage CRs
- No custom labels required by Kargo — existing platform CRs are label-free

**Directory placement:**
```
kargo/
  project.yaml
  warehouse.yaml
  stage-uat.yaml
  stage-prod.yaml
  kustomization.yaml
```

Deploy via a dedicated ArgoCD Application pointing to `path: kargo/` — follows the same
pattern as `manifests/apps/gpt-researcher/` but at the `kargo/` top-level path.

**`kargo/project.yaml` structure (from RESEARCH.md §2.1):**
```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Project
metadata:
  name: homelab-demo
spec:
  promotionPolicies:
    - stage: uat
      autoPromotionEnabled: true
    - stage: prod
      autoPromotionEnabled: false
```

**`kargo/warehouse.yaml` structure (from RESEARCH.md §2.2):**
```yaml
apiVersion: kargo.akuity.io/v1alpha1
kind: Warehouse
metadata:
  name: homelab-demo
  namespace: homelab-demo
spec:
  subscriptions:
    - git:
        repoURL: https://github.com/iemafzalhassan/homelab.git
        branch: main
        includePaths:
          - manifests/apps/homelab-demo/
```

---

## Pattern 8 — Keycloak OIDC Client (Documentation/Imperative)

### Analog: `docs/ADR-001-platform-auth-rbac.md` — OIDC Clients table (L68-74)

Existing client table (to be extended):
```markdown
| Client ID      | Redirect URI                                                  | Groups Claim Mapper |
|----------------|---------------------------------------------------------------|---------------------|
| `argocd`       | `https://argocd.smapatticare.com/auth/callback`              | YES - groups claim  |
| `jenkins`      | `https://jenkins.smapatticare.com/securityRealm/finishLogin` | YES - groups claim  |
| `oauth2-proxy` | `https://auth.smapatticare.com/oauth2/callback`              | YES - groups claim  |
```

**Row to add for Kargo:**
```markdown
| `kargo` | `https://kargo.smapatticare.com/auth/callback` | YES - groups claim (PKCE required) |
```

**Key differences from other clients:**
- Auth flow: **Authorization Code + PKCE** (mandatory for Kargo v1)
- The `groups` scope mapper already exists in Platform realm — reuse it
- Client secret stored as K8s Secret `kargo-oidc-secret` in `kargo` namespace

**RBAC mapping to add to ADR-001:**
```markdown
### Kargo (api.oidc claims mapping)
| Group            | Kargo Role  | Capability                              |
|------------------|-------------|-----------------------------------------|
| `homelab-admins` | `admin`     | Approve PROD, manage Projects           |
| `devops`         | `admin`     | Approve PROD on `homelab-demo` Project  |
| `developers`     | `viewer`    | View UAT status + Freight, no promote   |
| `viewers`        | `viewer`    | Read-only everywhere                    |
```

---

## Pattern 9 — Dockerfile (new file)

### Analog: `Utils/Jenkinsfile` ephemeral Dockerfile pattern

The existing Jenkinsfile (L10-11) dynamically creates a throwaway Alpine Dockerfile. The real
demo app needs a committed static Dockerfile.

**Existing pattern:**
```groovy
echo "FROM alpine:latest" > Dockerfile
echo "CMD echo 'Hello from Jenkins Spot Kaniko Build!'" >> Dockerfile
/kaniko/executor --context `pwd` --dockerfile `pwd`/Dockerfile --no-push
```

**Apply to `apps/homelab-demo/Dockerfile`:**
```dockerfile
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/index.html
EXPOSE 80
```

Minimal `index.html` shows version string for visible demo:
```html
<!DOCTYPE html>
<html><body><h1>homelab-demo: VERSION_PLACEHOLDER</h1></body></html>
```

The Jenkins pipeline replaces `VERSION_PLACEHOLDER` with `$IMAGE_TAG` before building.

---

## Pattern 10 — Jenkinsfile.demo (CI pipeline)

### Analog: `Utils/Jenkinsfile`

Existing structure:
```groovy
pipeline {
    agent { label 'kaniko' }
    stages {
        stage('Build with Kaniko') {
            steps {
                container('kaniko') {
                    sh '''
                        /kaniko/executor --context `pwd` --dockerfile `pwd`/Dockerfile --no-push
                    '''
                }
            }
        }
    }
}
```

**Apply to `apps/homelab-demo/Jenkinsfile.demo`:**
- Same `agent { label 'kaniko' }` and `container('kaniko')` wrapper
- Add `IMAGE_TAG = "git-${env.GIT_COMMIT[0..6]}"` variable
- Stage 1: `sed` to replace version placeholder in `index.html`
- Stage 2: Kaniko build + push to ACR (`--destination ${ACR_URL}/homelab-demo:${IMAGE_TAG}`)
- Stage 3: `git commit` of updated image tag to `manifests/apps/homelab-demo/uat/deployment.yaml`

**Jenkins job DSL to add to `jenkins/values.yaml` JCasC `jobs` block:**
```yaml
- script: >
    pipelineJob('homelab-demo') {
        definition {
            cpsScm {
                scm {
                    git {
                        remote {
                            url('https://github.com/iemafzalhassan/homelab.git')
                        }
                        branches('*/main')
                    }
                }
                scriptPath('apps/homelab-demo/Jenkinsfile.demo')
            }
        }
    }
```

---

## Data Flow Summary

```
Developer push to GitHub
    │
    ▼
Jenkins (kaniko agent — spot node)
    ├── builds nginx:alpine image with version tag
    ├── pushes to ACR as homelab-demo:git-<sha>
    └── commits tag update → manifests/apps/homelab-demo/uat/deployment.yaml
                                        │
                                        ▼
                          Kargo Warehouse (watches manifests/apps/homelab-demo/)
                                        │ creates Freight{commit: <sha>}
                                        ▼
                          Kargo Stage: uat (autoPromotionEnabled=true)
                                        │ argocd-update step
                                        ▼
                          ArgoCD Application: homelab-demo-uat
                                        │ kustomize sync
                                        ▼
                          Pod in demo-uat namespace ✓

DevOps clicks "Promote" in Kargo UI (kargo.smapatticare.com)
    │ (Keycloak OIDC login → groups claim → admin role)
    ▼
Kargo Stage: prod (autoPromotionEnabled=false → human gate)
    │ argocd-update step
    ▼
ArgoCD Application: homelab-demo-prod (manual sync — Kargo triggers it)
    │ kustomize sync
    ▼
Pod in demo-prod namespace ✓
```

---

## Critical Constraints

| Constraint | Source | Impact if violated |
|------------|--------|--------------------|
| `kargo.akuity.io/authorized-stage` annotation on ArgoCD Applications | RESEARCH §3.1 | Kargo cannot patch Applications |
| PROD ArgoCD Application MUST have NO `syncPolicy.automated` | RESEARCH §3.2 | Kargo and ArgoCD race; unpredictable PROD state |
| Kargo OIDC requires PKCE (Authorization Code + PKCE) | RESEARCH §4.1 | Login will fail |
| Kargo chart uses `oci://` scheme in ArgoCD source | RESEARCH §8 | ArgoCD will not find the chart |
| `root.yaml` uses `directory.recurse: false` | `manifests/apps/root.yaml` L15 | `homelab-demo/uat/` and `prod/` are NOT auto-picked up; explicit Application entries required |
| Kargo CRs MUST be in namespace matching Project name (`homelab-demo`) | RESEARCH §2.1 | Kargo controller ignores out-of-namespace CRs |

---

## PATTERN MAPPING COMPLETE
