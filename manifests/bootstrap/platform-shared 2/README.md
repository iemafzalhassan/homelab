# Platform shared-services design

## Intent

This platform should use a single shared identity and platform-services layer for the core tools:
- Keycloak for SSO/OIDC
- Argo CD for GitOps orchestration
- Jenkins for CI/CD
- Kargo for promotion
- Traefik for ingress
- Monitoring for observability

Application-specific workloads should live in their own app namespaces, but shared platform services should be consolidated where possible.

## Recommended model

### 1. Keep platform services in shared namespaces

Keep these in dedicated shared namespaces:
- argocd
- jenkins
- kargo
- monitoring
- traefik
- keycloak
- oauth2-proxy
- cnpg-system

This avoids mixing control-plane services with application workloads and keeps upgrades and RBAC boundaries predictable.

### 2. Prefer one shared PostgreSQL cluster for platform apps

Instead of one PostgreSQL cluster per app, use a shared CNPG cluster in the platform namespace and create multiple databases per consumer.

Example:
- `platform` database for shared platform metadata
- `keycloak` database for Keycloak
- `jenkins` database if Jenkins needs an internal DB
- `kargo` database if Kargo requires one

### 3. Use one shared Keycloak realm for all platform apps

The current setup already uses the `Platform` realm for ArgoCD, Jenkins, Kargo, and OAuth2-Proxy. That is the right model.

Use the same realm and same user/group model for:
- ArgoCD
- Jenkins
- Kargo
- OAuth2-Proxy
- future admin portals

### 4. Keep app namespaces isolated, but reuse the shared platform services

For customer or demo apps, use separate namespaces such as:
- demo-prod
- demo-uat
- homelab-demo

These namespaces should consume the shared platform services rather than each creating its own copy of the same base services.

## Why not put everything in one platform namespace?

Putting every application into a single namespace is usually not ideal because:
- it mixes production, staging, and experimental workloads together
- it makes RBAC and resource isolation harder
- it creates larger blast radius during upgrades or incidents
- it makes it harder to delete or isolate a single app without affecting others

A good compromise is:
- one shared platform namespace for control-plane services
- one namespace per application or environment

## Cost optimization approach

### Shared resources to prefer

- CNPG cluster shared across services instead of per-app Postgres deployments
- one Keycloak realm and one Keycloak instance instead of app-specific identity providers
- one Traefik ingress controller shared across all apps
- one monitoring stack shared across namespaces
- one shared ingress certificate strategy via cert-manager

### Avoided duplication

Avoid creating separate PostgreSQL, Keycloak, ingress, or observability stacks for each application unless the app has a hard compliance or isolation requirement.

## Suggested next steps

1. Add a shared CNPG cluster for platform databases.
2. Create separate databases inside the cluster for Keycloak, Jenkins, Kargo, and future apps.
3. Keep app-specific namespaces for workloads such as demo-prod and demo-uat.
4. Continue using the existing shared Keycloak realm for SSO.
5. Reuse the monitoring, ingress, and certificate stack across all applications.
