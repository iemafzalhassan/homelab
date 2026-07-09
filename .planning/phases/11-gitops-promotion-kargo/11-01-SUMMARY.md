# 11-01 Summary: Kargo Installation & SSO Integration

## Execution Summary
We have successfully deployed Kargo `v1.10.8` onto our AKS spot node pool and integrated it with our centralized Keycloak SSO server and Traefik Gateway API.

### Key Deliverables Completed
1. **Kargo Namespace & Secret Provisioning**:
   - Created the `kargo` namespace.
   - Configured `kargo-oidc-secret` with OIDC client credentials.
2. **Kargo Helm Deployment**:
   - Deployed Kargo via ArgoCD using the two-source Helm application pattern targeting OCI registry `oci://ghcr.io/akuity/kargo-charts/kargo`.
   - Set limits and requests matching the budget and capacity constraints.
   - Enabled spot pool scheduling with tolerations and node selectors.
3. **SSO Authentication**:
   - Integrated Keycloak Platform realm OIDC client (`kargo`) using PKCE.
   - Configured group-to-role mappings (`homelab-admins`/`devops` as admin, `developers`/`viewers` as viewer).
   - Disabled API TLS (`api.tls.enabled: false`) to avoid self-signed handshake errors, relying on upstream TLS termination at the Traefik Gateway, and set `api.tls.terminatedUpstream: true`.
4. **Gateway API Routing**:
   - Deployed an `HTTPRoute` for `kargo.smapatticare.com` targeting port `80` (HTTP port exposed by the service when TLS is disabled).
   - Configured `ReferenceGrant` allowing cross-namespace route binding from the `traefik` namespace.

### Verification Status
- **Namespace & Pods**: All Kargo pods are running and healthy on the spot node pool.
- **Routing**: `kargo.smapatticare.com` resolves, terminates TLS successfully at Traefik, and routes to Kargo API.
- **SSO Challenge**: Unauthenticated requests to API endpoints (`/v1beta1/projects`) return `401 Unauthorized` as expected. Public configuration `/v1beta1/system/public-server-config` successfully exposes the OIDC configuration.
