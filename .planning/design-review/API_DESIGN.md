# Kyros API Design

## API Overview
Kyros provides a comprehensive RESTful API for platform management, operations, and integration. The API follows REST principles with JSON payloads, predictable resource URLs, and standard HTTP response codes.

## API Principles
1. **RESTful**: Resource-based URLs, standard HTTP methods, stateless interactions
2. **Versioned**: Explicit versioning in URL path (/api/v1/, /api/v2/)
3. **Consistent**: Predictable naming conventions, consistent error formats
4. **Documented**: OpenAPI 3.0 specification available
5. **Secure**: Authentication required for all endpoints (except health checks)
6. **Efficient**: Pagination, filtering, and field selection for large datasets
7. **Extensible**: Additive changes only, clear deprecation policy

## Authentication and Authorization
All API endpoints require authentication unless explicitly documented otherwise.

### Authentication Methods
1. **Bearer Token**: JWT access token in Authorization header
   ```
   Authorization: Bearer <jwt-token>
   ```
2. **Basic Auth**: For legacy clients (not recommended for new implementations)
3. **API Keys**: For service-to-service authentication (internal use)

### Authorization Model
- **RBAC**: Role-Based Access Control with predefined roles
- **Resource-Based**: Permissions tied to specific resources (repositories, namespaces)
- **Policy Evaluation**: OPA integration for complex authorization decisions
- **Tenant Context**: All operations scoped to authenticated user's tenant

### Error Responses
Standard error format for all API errors:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": [
      {
        "field": "field_name",
        "issue": "specific validation issue"
      }
    ],
    "request_id": "uuid-for-tracing",
    "timestamp": "ISO-8601 timestamp"
  }
}
```

HTTP Status Codes:
- 200: Success (GET, PUT, PATCH)
- 201: Created (POST)
- 202: Accepted (asynchronous operations)
- 204: No Content (DELETE)
- 400: Bad Request (validation errors)
- 401: Unauthorized (missing/invalid authentication)
- 403: Forbidden (authenticated but insufficient permissions)
- 404: Not Found (resource doesn't exist)
- 409: Conflict (resource already exists or state conflict)
- 422: Unprocessable Entity (semantic errors)
- 429: Too Many Requests (rate limiting)
- 500: Internal Server Error
- 503: Service Unavailable

## API Versioning Strategy
- **URI Versioning**: `/api/v1/`, `/api/v2/` in URL path
- **Version Lifecycle**: 
  - Stable versions maintained for minimum 12 months
  - Deprecation notice 6 months before removal
  - Clear migration guides provided
- **Backward Compatibility**: 
  - Additive changes only within minor versions
  - No breaking changes in patch versions
  - Deprecated fields marked with deprecation warnings

## Common Patterns

### Pagination
List endpoints support cursor-based pagination:
```http
GET /api/v1/repositories?limit=50&cursor=abc123
```
Response:
```json
{
  "items": [...],
  "pagination": {
    "limit": 50,
    "cursor": "abc123",
    "next_cursor": "def456",
    "has_more": true,
    "total_count": 1234
  }
}
```

### Filtering
Filter query parameters for list endpoints:
```http
GET /api/v1/repositories?visibility=public&namespace_id=uuid123
```

### Field Selection
Clients can request specific fields to reduce payload size:
```http
GET /api/v1/repositories?fields=id,name,description,created_at
```

### Sorting
Sort query parameters:
```http
GET /api/v1/repositories?sort=-created_at,name
```
Prefix `-` for descending order.

### Request IDs
All responses include a `X-Request-ID` header for tracing:
```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

## API Endpoints

### 1. Health and Metadata
```http
GET /healthz
```
Liveness probe - returns 200 if service is running

```http
GET /readyz
```
Readiness probe - returns 200 if service is ready to serve traffic

```http
GET /api/v1/
```
API root - returns available versions and links

```http
GET /api/v1/version
```
Returns version information and build metadata

### 2. Authentication
```http
POST /api/v1/auth/login
```
Authenticate user and return access/refresh tokens
```json
{
  "username": "string",
  "password": "string",
  "totp_code": "string (optional for MFA)"
}
```
Response:
```json
{
  "access_token": "jwt-token",
  "refresh_token": "refresh-token",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

```http
POST /api/v1/auth/refresh
```
Refresh access token using refresh token
```json
{
  "refresh_token": "refresh-token"
}
```

```http
POST /api/v1/auth/logout
```
Revoke refresh token and end session
```json
{
  "refresh_token": "refresh-token"
}
```

```http
POST /api/v1/auth/mfa/enable
```
Enable MFA for user account
```json
{
  "password": "string"
}
```

```http
POST /api/v1/auth/mfa/verify
```
Verify MFA setup
```json
{
  "code": "string",
  "secret": "string"
}
```

### 3. Users and Groups
```http
GET /api/v1/users
```
List users with pagination and filtering
Query params: `limit`, `cursor`, `search`, `enabled`, `sort`

```http
POST /api/v1/users
```
Create new user
```json
{
  "username": "string",
  "email": "string",
  "display_name": "string",
  "password": "string",
  "send_welcome_email": "boolean (optional)"
}
```

```http
GET /api/v1/users/{user_id}
```
Get user details

```http
PATCH /api/v1/users/{user_id}
```
Update user
```json
{
  "display_name": "string (optional)",
  "enabled": "boolean (optional)"
}
```

```http
DELETE /api/v1/users/{user_id}
```
Delete user (soft delete)

```http
POST /api/v1/users/{user_id}/disable
```
Disable user account

```http
POST /api/v1/users/{user_id}/enable
```
Enable user account

```http
GET /api/v1/groups
```
List groups

```http
POST /api/v1/groups
```
Create new group
```json
{
  "name": "string",
  "description": "string (optional)"
}
```

```http
GET /api/v1/groups/{group_id}
```
Get group details

```http
POST /api/v1/groups/{group_id}/users/{user_id}
```
Add user to group

```http
DELETE /api/v1/groups/{group_id}/users/{user_id}
```
Remove user from group

### 4. Roles and Permissions
```http
GET /api/v1/roles
```
List roles with filtering
Query params: `scope`, `client_id`, `name`

```http
POST /api/v1/roles
```
Create new role
```json
{
  "name": "string",
  "description": "string (optional)",
  "scope": "realm|client",
  "client_id": "uuid (required if scope=client)",
  "permissions": ["string"] // Permission names
}
```

```http
GET /api/v1/roles/{role_id}
```
Get role details

```http
PATCH /api/v1/roles/{role_id}
```
Update role
```json
{
  "description": "string (optional)",
  "permissions": ["string"] (optional)
}
```

```http
DELETE /api/v1/roles/{role_id}
```
Delete role

```http
POST /api/v1/roles/{role_id}/assign/users/{user_id}
```
Assign role to user

```http
DELETE /api/v1/roles/{role_id}/assign/users/{user_id}
```
Revoke role from user

```http
GET /api/v1/permissions
```
List all available permissions

### 5. Clients (Applications/Services)
```http
GET /api/v1/clients
```
List clients

```http
POST /api/v1/clients
```
Create new client application
```json
{
  "client_id": "string",
  "name": "string",
  "redirect_uris": ["string"],
  "enabled": "boolean (optional)",
  "protocol": "openid-connect|saml (optional, default=openid-connect)"
}
```
Response includes generated secret (only shown once)

```http
GET /api/v1/clients/{client_id}
```
Get client details (secret not returned)

```http
PATCH /api/v1/clients/{client_id}
```
Update client
```json
{
  "name": "string (optional)",
  "redirect_uris": ["string"] (optional),
  "enabled": "boolean (optional)"
}
```

```http
DELETE /api/v1/clients/{client_id}
```
Delete client

```http
POST /api/v1/clients/{client_id}/secret/rotate
```
Rotate client secret
Response includes new secret (only shown once)

### 6. Namespaces
```http
GET /api/v1/namespaces
```
List namespaces with pagination and filtering
Query params: `limit`, `cursor`, `visibility`, `search`, `sort`

```http
POST /api/v1/namespaces
```
Create new namespace
```json
{
  "name": "string",
  "description": "string (optional)",
  "visibility": "public|private|protected (optional, default=private)"
}
```

```http
GET /api/v1/namespaces/{namespace_id}
```
Get namespace details

```http
PATCH /api/v1/namespaces/{namespace_id}
```
Update namespace
```json
{
  "description": "string (optional)",
  "visibility": "public|private|protected (optional)"
}
```

```http
DELETE /api/v1/namespaces/{namespace_id}
```
Delete namespace (only if empty)

```http
GET /api/v1/namespaces/{namespace_id}/repositories
```
List repositories in namespace
Query params: `limit`, `cursor`, `visibility`, `search`, `sort`

### 7. Repositories
```http
GET /api/v1/repositories
```
List repositories with pagination and filtering
Query params: `limit`, `cursor`, `namespace_id`, `visibility`, `search`, `sort`

```http
POST /api/v1/repositories
```
Create new repository
```json
{
  "name": "string",
  "namespace_id": "uuid",
  "description": "string (optional)",
  "visibility": "public|private|protected|inherit (optional, default=inherit)"
}
```

```http
GET /api/v1/repositories/{repository_id}
```
Get repository details

```http
PATCH /api/v1/repositories/{repository_id}
```
Update repository
```json
{
  "description": "string (optional)",
  "visibility": "public|private|protected|inherit (optional)"
}
```

```http
DELETE /api/v1/repositories/{repository_id}
```
Delete repository (only if empty)

```http
GET /api/v1/repositories/{repository_id}/artifacts
```
List artifacts in repository
Query params: `limit`, `cursor`, `sort`

```http
GET /api/v1/repositories/{repository_id}/tags
```
List tags in repository
Query params: `limit`, `cursor`, `sort`

### 8. Artifacts and Tags (Registry API)
Note: The OCI Distribution API is available at `/v2/` for standard registry operations
These endpoints are for management and metadata operations

```http
GET /api/v1/artifacts/{artifact_id}
```
Get artifact details (manifest, layers, size, etc.)

```http
GET /api/v1/artifacts/{artifact_id}/blobs
```
List blobs that make up this artifact

```http
GET /api/v1/artifacts/{artifact_id}/tags
```
List tags pointing to this artifact

```http
POST /api/v1/artifacts/{artifact_id}/tags
```
Create new tag pointing to artifact
```json
{
  "name": "string"
}
```

```http
DELETE /api/v1/artifacts/{artifact_id}/tags/{tag_name}
```
Delete tag

```http
GET /api/v1/blobs/{digest}
```
Get blob metadata and download URL (time-limited)
```

### 9. Trust and Security
```http
GET /api/v1/trust/scores
```
List trust scores with filtering
Query params: `limit`, `cursor`, `min_score`, `max_score`, `level`, `repository_id`, `namespace_id`, `sort`

```http
GET /api/v1/trust/scores/{artifact_id}
```
Get trust score for specific artifact

```http
POST /api/v1/trust/scores/{artifact_id}/recalculate
```
Trigger recalculation of trust score for artifact

```http
GET /api/v1/trust/vulnerabilities
```
List vulnerabilities with filtering
Query params: `limit`, `cursor`, `artifact_id`, `severity`, `scanner`, `sort`

```http
GET /api/v1/trust/vulnerabilities/{vulnerability_id}
```
Get vulnerability details

```http
GET /api/v1/trust/sboms
```
List SBOMs with filtering
Query params: `limit`, `cursor`, `artifact_id`, `format`, `sort`

```http
GET /api/v1/trust/sboms/{sbom_id}
```
Get SBOM details
```json
{
  "artifact_id": "uuid",
  "format": "SPDX|CycloneDX",
  "content": {...}, // SBOM document
  "generated_at": "ISO-8601 timestamp",
  "generator": "string"
}
```

```http
GET /api/v1/trust/signatures
```
List signatures with filtering
Query params: `limit`, `cursor`, `artifact_id`, `type`, `verified`, `sort`

```http
GET /api/v1/trust/signatures/{signature_id}
```
Get signature details

```http
POST /api/v1/trust/signatures/{signature_id}/verify
```
Verify signature
```json
{
  "signature": "string" // if not already stored
}
```

```http
GET /api/v1/trust/policies
```
List policies with filtering
Query params: `limit`, `cursor`, `scope`, `namespace_id`, `repository_id`, `enabled`, `sort`

```http
POST /api/v1/trust/policies
```
Create new policy
```json
{
  "name": "string",
  "description": "string (optional)",
  "rules": {...}, // Policy rules in Rego or similar
  "scope": "global|namespace|repository",
  "namespace_id": "uuid (required if scope=namespace)",
  "repository_id": "uuid (required if scope=repository)",
  "enabled": "boolean (optional, default=true)"
}
```

```http
GET /api/v1/trust/policies/{policy_id}
```
Get policy details

```http
PATCH /api/v1/trust/policies/{policy_id}
```
Update policy
```json
{
  "description": "string (optional)",
  "rules": {...} (optional),
  "enabled": "boolean (optional)"
}
```

```http
DELETE /api/v1/trust/policies/{policy_id}
```
Delete policy

```http
POST /api/v1/trust/policies/{policy_id}/evaluate
```
Evaluate policy against artifact
```json
{
  "artifact_id": "uuid"
}
```
Response:
```json
{
  "policy_id": "uuid",
  "artifact_id": "uuid",
  "result": "pass|fail|warn|error",
  "details": {...}, // Detailed evaluation results
  "evaluated_at": "ISO-8601 timestamp"
}
```

```http
GET /api/v1/trust/evaluations
```
List policy evaluations with filtering
Query params: `limit`, `cursor`, `artifact_id`, `policy_id`, `result`, `sort`

### 10. Observability
```http
GET /api/v1/metrics
```
Get metrics (Prometheus format)
```http
GET /api/v1/metrics?format=json
```
Get metrics in JSON format

```http
GET /api/v1/logs
```
Query logs with filtering
Query params: `limit`, `cursor`, `service`, `level`, `start_time`, `end_time`, `trace_id`, `search`, `sort`

```http
GET /api/v1/traces
```
List traces with filtering
Query params: `limit`, `cursor`, `service`, `operation`, `start_time`, `end_time`, `sort`

```http
GET /api/v1/traces/{trace_id}
```
Get trace details including spans

```http
GET /api/v1/alerts
```
List alerts with filtering
Query params: `limit`, `cursor`, `status`, `severity`, `sort`

```http
POST /api/v1/alerts/{alert_id}/resolve
```
Resolve alert manually
```json
{
  "resolution_note": "string (optional)"
}
```

```http
GET /api/v1/alert-rules
```
List alert rules

```http
POST /api/v1/alert-rules
```
Create new alert rule
```json
{
  "name": "string",
  "description": "string (optional)",
  "condition": "string", // PromQL or similar expression
  "severity": "info|warning|error|critical",
  "notification_channels": [...], // Configured notification targets
  "enabled": "boolean (optional, default=true)"
}
```

```http
GET /api/v1/alert-rules/{rule_id}
```
Get alert rule details

```http
PATCH /api/v1/alert-rules/{rule_id}
```
Update alert rule
```json
{
  "description": "string (optional)",
  "condition": "string (optional)",
  "severity": "info|warning|error|critical (optional)",
  "notification_channels": [...] (optional),
  "enabled": "boolean (optional)"
}
```

```http
DELETE /api/v1/alert-rules/{rule_id}
```
Delete alert rule

### 11. GitOps and Automation (Webhooks)
```http
GET /api/v1/webhooks
```
List webhooks with filtering
Query params: `limit`, `cursor`, `enabled`, `sort`

```http
POST /api/v1/webhooks
```
Create new webhook subscription
```json
{
  "name": "string",
  "url": "string",
  "events": ["string"], // Event types to subscribe to
  "secret": "string (optional, for HMAC verification)",
  "format": "JSON|form-urlencoded (optional, default=JSON)",
  "headers": {...}, // HTTP headers to include (optional)
  "enabled": "boolean (optional, default=true)"
}
```

```http
GET /api/v1/webhooks/{webhook_id}
```
Get webhook details (secret not returned)

```http
PATCH /api/v1/webhooks/{webhook_id}
```
Update webhook
```json
{
  "name": "string (optional)",
  "url": "string (optional)",
  "events": ["string"] (optional),
  "secret": "string (optional)",
  "format": "JSON|form-urlencoded (optional)",
  "headers": {...} (optional),
  "enabled": "boolean (optional)"
}
```

```http
DELETE /api/v1/webhooks/{webhook_id}
```
Delete webhook subscription

```http
GET /api/v1/webhooks/{webhook_id}/deliveries
```
List delivery attempts for webhook
Query params: `limit`, `cursor`, `status`, `sort`

```http
GET /api/v1/webhooks/deliveries/{delivery_id}
```
Get delivery attempt details

```http
POST /api/v1/webhooks/deliveries/{delivery_id}/retry
```
Retry failed delivery
```

### 12. Multi-tenancy
```http
GET /api/v1/tenants
```
List tenants with filtering
Query params: `limit`, `cursor`, `search`, `sort`

```http
POST /api/v1/tenants
```
Create new tenant
```json
{
  "name": "string",
  "display_name": "string (optional)",
  "description": "string (optional)"
}
```

```http
GET /api/v1/tenants/{tenant_id}
```
Get tenant details

```http
PATCH /api/v1/tenants/{tenant_id}
```
Update tenant
```json
{
  "display_name": "string (optional)",
  "description": "string (optional)"
}
```

```http
DELETE /api/v1/tenants/{tenant_id}
```
Delete tenant (only if no users or namespaces assigned)

```http
GET /api/v1/tenants/{tenant_id}/users
```
List users in tenant
Query params: `limit`, `cursor`, `role`, `sort`

```http
POST /api/v1/tenants/{tenant_id}/users/{user_id}
```
Add user to tenant
```json
{
  "role": "owner|admin|member|viewer (optional, default=member)"
}
```

```http
DELETE /api/v1/tenants/{tenant_id}/users/{user_id}
```
Remove user from tenant

```http
GET /api/v1/tenants/{tenant_id}/groups
```
List groups in tenant
Query params: `limit`, `cursor`, `role`, `sort`

```http
POST /api/v1/tenants/{tenant_id}/groups/{group_id}
```
Add group to tenant
```json
{
  "role": "owner|admin|member|viewer (optional, default=member)"
}
```

```http
DELETE /api/v1/tenants/{tenant_id}/groups/{group_id}
```
Remove group from tenant

```http
GET /api/v1/namespace-quotas
```
List namespace quotas with filtering
Query params: `limit`, `cursor`, `tenant_id`, `namespace_name`, `sort`

```http
POST /api/v1/namespace-quotas
```
Create new namespace quota
```json
{
  "tenant_id": "uuid",
  "namespace_name": "string",
  "hard_limits": {...} // Resource quotas (CPU, memory, storage, etc.)
}
```

```http
GET /api/v1/namespace-quotas/{quota_id}
```
Get namespace quota details

```http
PATCH /api/v1/namespace-quotas/{quota_id}
```
Update namespace quota
```json
{
  "hard_limits": {...} (optional),
  "used_resources": {...} (optional)
}
```

```http
DELETE /api/v1/namespace-quotas/{quota_id}
```
Delete namespace quota

### 13. Audit and Compliance
```http
GET /api/v1/audit/events
```
List audit events with filtering
Query params: `limit`, `cursor`, `start_time`, `end_time`, `actor_id`, `resource_type`, `outcome`, `sort`

```http
GET /api/v1/audit/events/{event_id}
```
Get audit event details

```http
GET /api/v1/audit/reports/compliance
```
Generate compliance report
Query params: `standard` (soc2, iso27001, gdpr, hipaa), `format` (json, pdf, csv), `start_time`, `end_time`

```http
GET /api/v1/audit/reports/activity
```
Generate user activity report
Query params: `user_id`, `tenant_id`, `start_time`, `end_time`, `format` (json, pdf, csv)

## Rate Limiting
API implements rate limiting to prevent abuse:
- **Anonymous**: 100 requests/hour per IP
- **Authenticated**: 1000 requests/hour per user
- **Burst**: Allow short bursts up to 5x the rate limit
- **Headers**: 
  - `X-RateLimit-Limit`: Request limit per window
  - `X-RateLimit-Remaining`: Requests remaining in current window
  - `X-RateLimit-Reset`: Seconds until limit resets
  - `Retry-After`: Seconds to wait before retrying (when limit exceeded)

When rate limit exceeded, returns 429 Too Many Requests with:
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded",
    "details": null,
    "request_id": "uuid",
    "timestamp": "ISO-8601 timestamp"
  }
}
```

## OpenAPI Specification
The complete API specification is available in OpenAPI 3.0 format:
- JSON: `/api/v1/openapi.json`
- YAML: `/api/v1/openapi.yaml`
- Swagger UI: `/api/v1/docs/` (interactive documentation)

## SDKs and Client Libraries
Kyros provides official client libraries for popular languages:
- **Go**: `github.com/kyros-project/kyros-go-client`
- **Python**: `kyros-py-client` (PyPI)
- **JavaScript/TypeScript**: `@kyros/client` (npm)
- **Java**: `com.kyros:kyros-java-client` (Maven)
- **.NET**: `Kyros.Net.Client` (NuGet)

## Deprecation Policy
- **Deprecation Notice**: Deprecated endpoints return `Warning` header with sunset date
- **Migration Guide**: Available at `/api/v1/deprecations/{endpoint}`
- **Sunset Period**: Minimum 6 months between deprecation and removal
- **Removal**: Deprecated endpoints return 410 Gone after sunset date

## Examples

### Authentication Flow
```bash
# Login
curl -X POST https://kyros.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "secure-password"}'

# Response includes access token
# {
#   "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
#   "expires_in": 900,
#   "token_type": "Bearer"
# }

# Use access token for subsequent requests
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  https://kyros.example.com/api/v1/repositories
```

### Creating a Repository
```bash
curl -X POST https://kyros.example.com/api/v1/repositories \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-app",
    "namespace_id": "123e4567-e89b-12d3-a456-426614174000",
    "description": "My application container images",
    "visibility": "private"
  }'
```

### Pushing an Image (via Registry API)
```bash
# Standard Docker Registry API v2 at /v2/
docker tag my-app:latest kyros.example.com/my-app:latest
docker push kyros.example.com/my-app:latest
```

### Checking Trust Score
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  https://kyros.example.com/api/v1/trust/scores?repository_id=123e4567-e89b-12d3-a456-426614174000
```

### Setting Up a Webhook
```bash
curl -X POST https://kyros.example.com/api/v1/webhooks \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "CI Pipeline Trigger",
    "url": "https://ci.example.com/webhook/kyros",
    "events": ["artifact.pushed", "trust.score.updated"],
    "secret": "my-webhook-secret",
    "enabled": true
  }'
```

## Future API Enhancements
1. **GraphQL API**: Alternative to REST for flexible data fetching
2. **Real-time Updates**: WebSocket subscriptions for live updates
3. **Batch Operations**: Bulk create/update/delete operations
4. **Async Job Endpoints**: For long-running operations with status polling
5. **Webhook Templates**: Pre-configured webhook configurations for common integrations
6. **API Keys**: Programmatic access without JWT tokens for service accounts
7. **Field Expansion**: Ability to expand referenced objects in a single request
8. **Geolocation Awareness**: API endpoint selection based on client location
9. **API Usage Analytics**: Built-in tracking of API usage patterns
10. **Mock API**: Sandbox environment for testing and development

## Diagrams Reference
See [MERMAID.md](MERMAID.md) for detailed Mermaid diagrams including:
- API Endpoint Categories
- Authentication Flow
- Request/Response Patterns
- Rate Limiting Architecture
- Versioning Strategy
- Error Handling Flow