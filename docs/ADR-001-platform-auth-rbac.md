# ADR-001: Platform Authentication & RBAC Architecture

**Status:** Active  
**Date:** 2026-07-09  
**Decision by:** Afzal Hassan (@iemafzalhassan)

---

## Context

The homelab AKS platform uses Keycloak as a centralised Identity Provider (IdP) for Single Sign-On (SSO) across all platform tools. Multiple integration issues were discovered during initial setup that prevented users from logging in. This ADR documents the complete authentication architecture, known pitfalls, and their resolutions.

---

## Decision: Centralised OIDC via Keycloak Platform Realm

All platform applications authenticate via a single Keycloak realm (`Platform`) using the OIDC Authorization Code Flow. Keycloak issues JWT tokens carrying a `groups` claim that downstream applications use for RBAC — no application manages its own user database.

---

## Authentication Flow

```
User Browser
    │
    ▼
[Traefik Gateway] ─── TLS Terminate ───► [Application (ArgoCD/Jenkins)]
    │
    │  Redirect to Keycloak
    ▼
[Keycloak /realms/Platform]
    │  issues JWT with `groups` claim
    ▼
[Application verifies token]
    │  maps groups → local roles
    ▼
[User granted access per RBAC matrix]
```

For applications protected by OAuth2-Proxy (e.g. future admin-only services):
```
User → Traefik → OAuth2-Proxy → Keycloak (OIDC) → Cookie set → Upstream
```

---

## Keycloak Configuration

### Realm Structure

| Realm | Purpose |
|-------|---------|
| `master` | Keycloak system administration only |
| `Platform` | All homelab admin tools: ArgoCD, Jenkins, OAuth2-Proxy |
| `UAT` | Applications deployed in the UAT environment |
| `PROD` | Applications deployed in the PROD environment |

### Platform Realm: Groups

| Group | Members | Purpose |
|-------|---------|---------|
| `homelab-admins` | `admin@smapatticare.com`, `iemafzalhassan@gmail.com` | Super administrators |
| `devops` | `devops@smapatticare.com` | DevOps operators |
| `developers` | `dev@smapatticare.com` | Software developers |
| `viewers` | `viewer@smapatticare.com` | Read-only auditors |

### Platform Realm: OIDC Clients

| Client ID | Redirect URI | Groups Claim Mapper |
|-----------|-------------|-------------------|
| `argocd` | `https://argocd.smapatticare.com/auth/callback` | ✅ `oidc-group-membership-mapper` → `groups` claim |
| `jenkins` | `https://jenkins.smapatticare.com/securityRealm/finishLogin` | ✅ `oidc-group-membership-mapper` → `groups` claim |
| `oauth2-proxy` | `https://auth.smapatticare.com/oauth2/callback` | ✅ `oidc-group-membership-mapper` → `groups` claim |
| `kargo` | `https://kargo.smapatticare.com/auth/callback` | ✅ `oidc-group-membership-mapper` → `groups` claim (PKCE required) |

---

## RBAC Mapping Matrix

### ArgoCD (`configs.rbac.policy.csv`)

```csv
g, homelab-admins, role:admin
g, devops, role:admin
g, developers, role:readonly
g, viewers, role:readonly
```

`policy.default: role:''` — unauthenticated / unrecognised groups are denied.

### Jenkins (JCasC `authorizationStrategy.globalMatrix.permissions`)

| Permission | `homelab-admins` | `devops` | `developers` | `viewers` |
|-----------|:---:|:---:|:---:|:---:|
| `Overall/Administer` | ✅ | ✅ | ❌ | ❌ |
| `Overall/Read` | (inherited) | (inherited) | ✅ (authenticated) | ✅ (authenticated) |
| `Job/Build` | ✅ | ✅ | ✅ | ❌ |
| `Job/Cancel` | ✅ | ✅ | ✅ | ❌ |
| `Job/Read` | ✅ | ✅ | ✅ | ✅ |
| `Job/Workspace` | ✅ | ✅ | ✅ | ❌ |
| `Run/Replay` | ✅ | ✅ | ✅ | ❌ |
| `Run/Update` | ✅ | ✅ | ✅ | ❌ |

### Keycloak Admin Console

| Group | Access |
|-------|--------|
| `homelab-admins` | Full realm admin via `realm-management:realm-admin` client role |
| `devops` | _(Future: scope client roles for user/client mgmt only)_ |
| `developers` | ❌ No login access to Keycloak console |
| `viewers` | ❌ No login access to Keycloak console |

### Kargo (api.oidc claims mapping)

| Group | Kargo Role | Capability |
|-------|------------|------------|
| `homelab-admins` | `admin` | Approve PROD, manage Projects |
| `devops` | `admin` | Approve PROD on `homelab-demo` Project |
| `developers` | `viewer` | View UAT status + Freight, no promote |
| `viewers` | `viewer` | Read-only everywhere |

---

## Known Issues & Resolutions (Post-Mortem)

### Issue 1: OAuth2-Proxy `500` — `email isn't verified`

**Symptom:** After Keycloak login, OAuth2-Proxy redirects to `/oauth2/callback` and returns `HTTP 500` with log:
```
Error redeeming code during OAuth2 callback: email in id_token (dev@smapatticare.com) isn't verified
```

**Root Cause:** All Keycloak users were created with `emailVerified: false`. By default, OAuth2-Proxy (and many OIDC clients) enforce that the email address in the JWT must be verified (`email_verified: true` claim). Keycloak does NOT auto-verify emails unless configured to do so or email verification is explicitly triggered.

**Resolution:** Set `emailVerified=true` on all Platform realm users via `kcadm.sh`:
```bash
kubectl exec keycloak-0 -n keycloak -- /opt/bitnami/keycloak/bin/kcadm.sh \
  update users/{USER_ID} -s emailVerified=true -r Platform --config /tmp/kcadm.config
```

**Permanent Fix:** When creating new Keycloak users, always set `emailVerified: true` immediately. For future automation, add this to the Keycloak realm export/import bootstrap script.

---

### Issue 2: All users land in no group — RBAC has no effect

**Symptom:** Users can authenticate but ArgoCD and Jenkins give them no permissions or wrong permissions. The `groups` claim in the JWT is an empty array `[]`.

**Root Cause:** The Keycloak groups `devops`, `developers`, `viewers` did not exist. Only `homelab-admins` was created. No users were assigned to any group (including `homelab-admins`). The `groups` claim mapper was present on all clients but returned `[]` because users had no group memberships.

**Resolution:**
1. Created missing groups: `devops`, `developers`, `viewers`
2. Assigned users to groups per the RBAC matrix above
3. Updated ArgoCD `policy.csv` and Jenkins JCasC permissions to reference all 4 groups

---

### Issue 3: Keycloak `sso-tls-secret` covers wrong SANs initially

**Symptom:** The original `sso-tls-secret` certificate only covered `auth.smapatticare.com`, not `sso.smapatticare.com`.

**Root Cause:** The Traefik Gateway `keycloak` listener was added after the initial certificate was issued for the `authsecure` listener. cert-manager's Gateway shim combines DNS names from all listeners referencing the same secret name — but this only triggers a certificate refresh when the `Certificate` resource's `spec.dnsNames` diverges from the existing secret.

**Resolution:** Traefik Helm chart was updated to define separate listeners (`authsecure` for `auth.smapatticare.com`, `keycloak` for `sso.smapatticare.com`) both referencing `sso-tls-secret`. cert-manager automatically detected the SAN change and reissued the certificate with both domains.

---

## Observability Integration Notes

- Keycloak logs are collected by `alloy-logs` DaemonSet in the `monitoring` namespace via pod log scraping
- Authentication failures (`Error redeeming code`, `email isn't verified`) are visible in Loki with label `{namespace="oauth2-proxy"}`
- Future: add a Grafana alert on `count_over_time({namespace="oauth2-proxy"} |= "Error redeeming code" [5m]) > 3`

---

## Files Modified by This Decision

| File | Change |
|------|--------|
| `manifests/bootstrap/argocd/values.yaml` | Added `devops`, `developers`, `viewers` to `policy.csv`; added `policy.default: role:''` |
| `manifests/bootstrap/jenkins/values.yaml` | Added group-scoped permissions to `globalMatrix` for `devops`, `developers`, `viewers` |
| `manifests/bootstrap/kargo/values.yaml` | Added Kargo OIDC config and mapped Keycloak groups |
| `manifests/bootstrap/traefik/values.yaml` | Added HTTPS listener for Kargo and certificate generation |

All Keycloak changes (group creation, user assignments, `emailVerified`) were applied imperatively via `kcadm.sh`. These should be captured in a future Keycloak realm export and committed as `manifests/bootstrap/keycloak/realm-export.json` for GitOps reproducibility.
