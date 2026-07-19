# Kyros Domain Model

## Business Domains

Kyros encompasses several interconnected business domains that collectively provide a comprehensive software supply chain platform:

### 1. Registry Domain
**Responsibility**: Storage, distribution, and management of OCI-compliant artifacts (container images, Helm charts, SBOMs, etc.)

**Core Concepts**:
- **Repository**: Collection of related artifacts (typically images of the same application)
- **Namespace**: Logical grouping of repositories (often corresponding to teams or projects)
- **Artifact**: Stored OCI object (image manifest, Helm chart, SBOM, etc.)
- **Tag**: Mutable pointer to an artifact version
- **Digest**: Immutable content-addressed identifier for an artifact
- **Blob**: Content-addressed storage unit (layers, config, etc.)

### 2. Identity and Access Management (IAM) Domain
**Responsibility**: Authentication, authorization, and user/identity management

**Core Concepts**:
- **User**: Individual human actor in the system
- **Group**: Collection of users for simplified permission management
- **Service Account**: Non-human identity for automated processes
- **Role**: Collection of permissions that can be assigned to users/groups
- **Permission**: Specific action that can be performed on a resource
- **Realm**: Security boundary containing users, groups, roles, and clients
- **Client**: Application or service that requests authentication

### 3. Trust and Security Domain
**Responsibility**: Security scanning, vulnerability management, signature verification, and trust scoring

**Core Concepts**:
- **Vulnerability**: Security weakness in an artifact or its dependencies
- **SBOM (Software Bill of Materials)**: Inventory of components in an artifact
- **Signature**: Cryptographic proof of artifact origin and integrity
- **Trust Score**: Quantitative measure of artifact security and quality
- **Policy**: Rule that defines acceptable artifact characteristics
- **Quarantine**: Isolation of artifacts that fail policy checks
- **Allowlist/Denylist**: Lists of trusted/untrusted entities

### 4. Observability Domain
**Responsibility**: Collection, storage, and visualization of metrics, logs, and traces

**Core Concepts**:
- **Metric**: Quantitative measurement of system behavior
- **Log**: Immutable record of discrete events over time
- **Trace**: Record of a request's journey through distributed systems
- **Dashboard**: Visualization of metrics and logs
- **Alert**: Notification triggered by metric thresholds
- **Trace ID**: Unique identifier for correlating events across services

### 5. GitOps and Automation Domain
**Responsibility**: Integration with GitOps workflows and automation capabilities

**Core Concepts**:
- **Webhook**: HTTP callback triggered by system events
- **Pipeline**: Automated sequence of operations (build, test, deploy)
- **Promotion**: Movement of artifacts between environments
- **Manifest**: Declarative description of desired system state
- **Sync**: Process of aligning actual state with desired state
- **Hook**: Extension point for custom automation logic

### 6. Multi-tenancy Domain
**Responsibility**: Isolation and resource management for multiple tenants

**Core Concepts**:
- **Tenant**: Isolated environment for a team, project, or organization
- **Namespace**: Kubernetes namespace providing isolation boundary
- **Quota**: Limits on resource consumption (storage, bandwidth, etc.)
- **Resource**: Consumable entity (CPU, memory, storage, etc.)
- **LimitRange**: Default resource limits for containers in a namespace
- **PriorityClass**: Priority assignment for pod scheduling

## Bounded Contexts

Following Domain-Driven Design principles, Kyros implements these bounded contexts:

### 1. Registry Context
**Boundary**: OCI artifact storage and distribution
**Shared Kernel**: None (independent context)
**Relationships**:
- Customer-Supplier with Trust Context (provides artifacts for scanning)
- Conformist with IAM Context (consumes authentication/authorization)
- Anticorruption Layer with GitOps Context (translates webhook events)

### 2. Identity and Access Management Context
**Boundary**: Authentication, authorization, and identity management
**Shared Kernel**: None (independent context)
**Relationships**:
- Shared Kernel with all other contexts (provides authz/authentication)
- Conformist relationship with external Keycloak (adapts to external IdP)

### 3. Trust and Security Context
**Boundary**: Security analysis, vulnerability management, and trust scoring
**Shared Kernel**: Vulnerability data formats with external scanners
**Relationships**:
- Customer-Supplier with Registry Context (consumes artifacts for analysis)
- Conformist with Policy Context (implements external policy standards)
- Anticorruption Layer with Notification Context (formats alerts for delivery)

### 4. Observability Context
**Boundary**: Metrics, logging, tracing, and alerting
**Shared Kernel**: OpenTelemetry instrumentation library
**Relationships**:
- Shared Kernel with all contexts (provides observability data)
- Publisher-Subscriber with Notification Context (publishes alerts)

### 5. GitOps and Automation Context
**Boundary**: Webhooks, automation pipelines, and GitOps integration
**Shared Kernel**: Webhook payload formats
**Relationships**:
- Customer-Supplier with Registry Context (consumes registry events)
- Conformist with external GitOps tools (Adapts to ArgoCD/Kargo formats)
- Anticorruption Layer with Notification Context (formats webhook payloads)

### 6. Multi-tenancy Context
**Boundary**: Tenant isolation, quotas, and resource management
**Shared Kernel**: Resource quota definitions with Kubernetes
**Relationships**:
- Shared Kernel with Registry Context (enforces namespace quotas)
- Customer-Supplier with IAM Context (consumes tenant/user mappings)
- Conformist with Kubernetes (implements Kubernetes namespace concepts)

## Entities and Value Objects

### Registry Context Entities
- **Repository** (Entity)
  - ID: UUID
  - Name: String (unique within namespace)
  - NamespaceID: UUID (foreign key)
  - Description: String (optional)
  - Visibility: Enum (public, private, protected)
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp
  - CreatedBy: UUID (foreign key to User)
  - UpdatedBy: UUID (foreign key to User)

- **Artifact** (Entity)
  - ID: UUID
  - RepositoryID: UUID (foreign key)
  - Digest: String (SHA256, unique)
  - MediaType: String (OCI media type)
  - Size: Integer (bytes)
  - CreatedAt: Timestamp
  - UploadedBy: UUID (foreign key to User or ServiceAccount)

- **Tag** (Entity)
  - ID: UUID
  - ArtifactID: UUID (foreign key)
  - Name: String (unique within repository)
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp

- **Blob** (Value Object)
  - Digest: String (SHA256)
  - Size: Integer (bytes)
  - MediaType: String (OCI media type)
  - UploadedAt: Timestamp

### IAM Context Entities
- **User** (Entity)
  - ID: UUID
  - Username: String (unique)
  - Email: String (unique)
  - DisplayName: String
  - Enabled: Boolean
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp
  - LastLoginAt: Timestamp (nullable)
  - EmailVerified: Boolean

- **Group** (Entity)
  - ID: UUID
  - Name: String (unique within realm)
  - Description: String (optional)
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp

- **Role** (Entity)
  - ID: UUID
  - Name: String (unique within realm)
  - Description: String (optional)
  - Scope: Enum (realm, client-specific)
  - ClientID: UUID (foreign key to Client, nullable for realm roles)
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp

- **Client** (Entity)
  - ID: UUID
  - ClientID: String (unique identifier)
  - Name: String
  - Secret: String (hashed)
  - RedirectURIs: String array
  - Enabled: Boolean
  - Protocol: Enum (openid-connect, saml)
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp

### Trust and Security Context Entities
- **Vulnerability** (Entity)
  - ID: UUID
  - ArtifactID: UUID (foreign key)
  - Scanner: String (Trivy, Grype, etc.)
  - VulnerabilityID: String (CVE identifier or scanner-specific)
  - Severity: Enum (unknown, low, medium, high, critical)
  - Title: String
  - Description: String
  - References: String array
  - FixedVersion: String (nullable)
  - DiscoveredAt: Timestamp

- **SBOM** (Entity)
  - ID: UUID
  - ArtifactID: UUID (foreign key)
  - Format: Enum (SPDX, CycloneDX)
  - Content: JSON (SBOM document)
  - GeneratedAt: Timestamp
  - Generator: String (Syft, etc.)

- **Signature** (Entity)
  - ID: UUID
  - ArtifactID: UUID (foreign key)
  - Type: Enum (cosign, notary, pgp)
  - KeyID: String (identifier of signing key)
  - Payload: String (signed data)
  - Signature: String (cryptographic signature)
  - VerifiedAt: Timestamp (nullable)
  - VerificationStatus: Enum (pending, verified, failed)

- **TrustScore** (Entity)
  - ID: UUID
  - ArtifactID: UUID (foreign key)
  - Score: Float (0.0-1.0)
  - Level: Enum (unknown, untrusted, low, medium, high, trusted)
  - Factors: JSON (breakdown of scoring factors)
  - CalculatedAt: Timestamp
  - ExpiresAt: Timestamp (nullable)

- **Policy** (Entity)
  - ID: UUID
  - Name: String (unique)
  - Description: String (optional)
  - Rules: JSON (policy rules in Rego or similar)
  - Scope: Enum (global, namespace, repository)
  - NamespaceID: UUID (foreign key, nullable for global policies)
  - RepositoryID: UUID (foreign key, nullable for non-repository policies)
  - Enabled: Boolean
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp

### Observability Context Entities
- **Metric** (Value Object)
  - Name: String
  - Value: Float or Integer
  - Tags: Map[String, String]
  - Timestamp: Timestamp
  - Unit: String (optional)

- **LogEntry** (Entity)
  - ID: UUID
  - Service: String
  - Level: Enum (debug, info, warn, error, fatal)
  - Message: String
  - TraceID: String (nullable)
  - SpanID: String (nullable)
  - Attributes: Map[String, String]
  - Timestamp: Timestamp

- **Trace** (Entity)
  - ID: UUID (TraceID)
  - ServiceName: String
  - OperationName: String
  - StartTime: Timestamp
  - Duration: Integer (nanoseconds)
  - Tags: Map[String, String]
  - Logs: LogEntry array
  - Status: Enum (ok, error)

- **Alert** (Entity)
  - ID: UUID
  - RuleID: UUID (foreign key to AlertRule)
  - Status: Enum (pending, firing, resolved)
  - Severity: Enum (info, warning, error, critical)
  - Summary: String
  - Description: String
  - StartedAt: Timestamp
  - EndedAt: Timestamp (nullable)
  - Value: String (current value that triggered alert)

### GitOps and Automation Context Entities
- **Webhook** (Entity)
  - ID: UUID
  - Name: String (unique)
  - URL: String
  - Events: String array (subscribed event types)
  - Secret: String (hashed, for HMAC verification)
  - Format: Enum (JSON, form-urlencoded)
  - Headers: Map[String, String]
  - Enabled: Boolean
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp
  - LastTriggeredAt: Timestamp (nullable)
  - FailureCount: Integer
  - NextRetryAt: Timestamp (nullable)

- **WebhookDelivery** (Entity)
  - ID: UUID
  - WebhookID: UUID (foreign key)
  - EventID: String (NATS JetStream message ID or similar)
  - Attempt: Integer
  - Status: Enum (pending, success, failed)
  - ResponseCode: Integer (nullable)
  - ResponseBody: String (truncated)
  - AttemptedAt: Timestamp
  - CompletedAt: Timestamp (nullable)

### Multi-tenancy Context Entities
- **Tenant** (Entity)
  - ID: UUID
  - Name: String (unique)
  - DisplayName: String
  - Description: String (optional)
  - CreatedAt: Timestamp
  - UpdatedAt: Timestamp
  - CreatedBy: UUID (foreign key to User)

- **NamespaceQuota** (Entity)
  - ID: UUID
  - TenantID: UUID (foreign key)
  - NamespaceID: String (Kubernetes namespace name)
  - HardLimits: JSON (resource quotas)
  - Used: JSON (current usage)
  - UpdatedAt: Timestamp

## Relationships

### Registry Context Relationships
- Repository **contains** zero or more Artifacts (one-to-many)
- Repository **has** zero or more Tags pointing to Artifacts (one-to-many)
- Artifact **is identified by** exactly one Digest (one-to-one)
- Artifact **consists of** one or more Blobs (one-to-many)
- Tag **points to** exactly one Artifact (many-to-one)
- Namespace **contains** zero or more Repositories (one-to-many)

### IAM Context Relationships
- User **belongs to** zero or more Groups (many-to-many)
- User **can be assigned** zero or more Roles (many-to-many)
- Group **can be assigned** zero or more Roles (many-to-many)
- Role **grants** one or more Permissions (one-to-many)
- Client **defines** zero or more Roles (one-to-many)
- Realm **contains** zero or more Users, Groups, Roles, Clients (one-to-many each)

### Trust and Security Context Relationships
- Artifact **may have** zero or more Vulnerabilities (one-to-many)
- Artifact **may have** exactly one SBOM (one-to-one, optional)
- Artifact **may have** zero or more Signatures (one-to-many)
- Artifact **has** exactly one TrustScore (one-to-one, optional)
- Policy **applies to** zero or more Artifacts via evaluation (one-to-many)
- Namespace **may have** zero or more Policies (one-to-many)
- Repository **may have** zero or more Policies (one-to-many)

### Observability Context Relationships
- Service **emits** zero or more Metrics (one-to-many)
- Service **emits** zero or more LogEntries (one-to-many)
- Trace **records** one service interaction (one-to-one, represents a request)
- Span **is part of** exactly one Trace (many-to-one)
- LogEntry **may be associated** with exactly one Trace via TraceID (many-to-one, optional)
- AlertRule **triggers** zero or more Alerts (one-to-many)
- Alert **is caused by** exactly one AlertRule (many-to-one)

### GitOps and Automation Context Relationships
- Webhook **delivers** zero or more WebhookDeliveries (one-to-many)
- WebhookDelivery **is attempt** to deliver for exactly one Webhook (many-to-one)
- Event **triggers** zero or more Webhooks (many-to-many via subscription)
- WebhookDelivery **represents** delivery attempt of exactly one Event (many-to-one)

### Multi-tenancy Context Relationships
- Tenant **owns** zero or more Namespaces (one-to-many)
- Namespace **belongs to** exactly one Tenant (many-to-one)
- Tenant **has** zero or more NamespaceQuotas (one-to-many)
- NamespaceQuota **applies to** exactly one Namespace (many-to-one)
- User **belongs to** exactly one Tenant (many-to-one, via tenant mapping)
- Group **may be associated** with zero or more Tenants (many-to-many)

## Aggregates

### Registry Aggregate
**Root**: Repository
**Boundaries**: 
- Includes: Repository, Artifact, Tag, Blob
- Excludes: Namespace (referenced by ID), User/ServiceAccount (referenced by ID)
**Invariants**:
- Tag name must be unique within a repository
- Artifact digest must be unique globally
- Blob digest must be unique globally
- Repository name must be unique within a namespace
**Operations**:
- CreateRepository(name, description, visibility)
- DeleteRepository(hard: boolean)
- PushArtifact(manifest, layers, config) -> Artifact
- PullArtifact(digest) -> Artifact layers and manifest
- TagArtifact(artifactDigest, tagName)
- DeleteTag(tagName)
- GarbageCollectUnreferencedBlobs()

### IAM Aggregate
**Root**: Realm
**Boundaries**:
- Includes: Realm, User, Group, Role, Client, Permission
- Excludes: None (self-contained)
**Invariants**:
- Username must be unique within realm
- Email must be unique within realm
- ClientID must be unique within realm
- Role name must be unique within scope (realm or client)
**Operations**:
- CreateUser(username, email, displayName, password)
- UpdateUser(userID, updates)
- DisableUser(userID)
- CreateGroup(name, description)
- AddUserToGroup(userID, groupID)
- RemoveUserFromGroup(userID, groupID)
- CreateRole(name, description, scope, clientID, permissions)
- AssignRoleToUser(roleID, userID)
- AssignRoleToGroup(roleID, groupID)
- CreateClient(clientID, name, secret, redirectURIs, protocol)

### Trust and Security Aggregate
**Root**: Artifact
**Boundaries**:
- Includes: Artifact, Vulnerability, SBOM, Signature, TrustScore
- Excludes: Policy (referenced by ID), Scanner services (external)
**Invariants**:
- An artifact can have at most one SBOM per format
- Trust score must be between 0.0 and 1.0
- Vulnerability severity must be valid enum value
**Operations**:
- ScanArtifactForVulnerabilities(artifactID) -> Vulnerability array
- GenerateSBOM(artifactID, format) -> SBOM
- VerifySignature(artifactID, signatureID) -> VerificationResult
- CalculateTrustScore(artifactID, policyID) -> TrustScore
- EvaluatePolicy(policyID, artifactID) -> PolicyResult

### Observability Aggregate
**Root**: ServiceInstance
**Boundaries**:
- Includes: ServiceInstance, Metric, LogEntry, Trace, Span
- Excludes: AlertRule (referenced by ID)
**Invariants**:
- Metric timestamp must be valid
- LogEntry timestamp must be valid
- Trace startTime must be before endTime (if ended)
**Operations**:
- RecordMetric(name, value, tags, timestamp)
- LogEntry(service, level, message, traceID, spanID, attributes, timestamp)
- StartTrace(traceID, serviceName, operationName, startTime, tags)
- EndTrace(traceID, endTime, status)
- AddTraceLog(traceID, logEntry)
- CreateAlertRule(name, condition, severity, notificationChannels)
- EvaluateAlertRule(ruleID, metricValues) -> Alert array

### GitOps and Automation Aggregate
**Root**: Webhook
**Boundaries**:
- Includes: Webhook, WebhookDelivery
- Excludes: Event definitions (shared kernel), Notification channels (external)
**Invariants**:
- Webhook URL must be valid
- Webhook secret must be properly formatted if provided
- Delivery attempt number must be positive
**Operations**:
- CreateWebhook(name, url, events, secret, format, headers)
- UpdateWebhook(webhookID, updates)
- DeleteWebhook(webhookID)
- TriggerWebhook(webhookID, eventType, payload) -> WebhookDelivery
- RetryFailedDelivery(deliveryID) -> WebhookDelivery
- GetWebhookDeliveries(webhookID, statusFilter, limit, offset)

### Multi-tenancy Aggregate
**Root**: Tenant
**Boundaries**:
- Includes: Tenant, NamespaceQuota
- Excludes: Namespace (Kubernetes concept), User/Group (IAM context)
**Invariants**:
- Tenant name must be unique
- Namespace quota must reference valid namespace
- Used resources cannot exceed hard limits (eventually consistent)
**Operations**:
- CreateTenant(name, displayName, description)
- UpdateTenant(tenantID, updates)
- DeleteTenant(tenantID, force: boolean)
- CreateNamespaceQuota(tenantID, namespaceName, hardLimits)
- UpdateNamespaceQuota(quotaID, usedResources)
- DeleteNamespaceQuota(quotaID)
- GetTenantUsage(tenantID) -> ResourceUsage breakdown

## Ownership and Responsibility Matrix

| Entity/Aggregate | Primary Owner | Secondary Owners | Responsibility |
|------------------|---------------|------------------|----------------|
| Repository | Registry Context | IAM Context (access control) | Storage and lifecycle management of artifacts |
| Artifact | Registry Context | Trust Context (scoring), Observability (metrics) | Immutable storage of OCI objects |
| Tag | Registry Context |  | Mutable pointers to artifact versions |
| Blob | Registry Context |  | Content-addressed storage of artifact layers |
| User | IAM Context | All Contexts (authentication) | Identity management and authentication |
| Group | IAM Context |  | Grouping of users for permission management |
| Role | IAM Context | All Contexts (authorization) | Definition of permissions and access rights |
| Client | IAM Context |  | Application/service identity for authentication |
| Vulnerability | Trust Context | Registry Context (artifact linkage) | Security vulnerability tracking |
| SBOM | Trust Context |  | Software bill of materials storage |
| Signature | Trust Context | Registry Context (artifact linkage) | Cryptographic signature verification |
| TrustScore | Trust Context | Registry Context (artifact linkage), API (exposure) | Trust score calculation and storage |
| Policy | Trust Context |  | Security and quality policy definition |
| Metric | Observability Context | All Contexts (instrumentation) | Quantitative system behavior measurement |
| LogEntry | Observability Context | All Contexts (instrumentation) | Immutable event recording |
| Trace | Observability Context | All Contexts (instrumentation) | Distributed request tracing |
| Alert | Observability Context | Notification Context (delivery) | System condition notification |
| Webhook | GitOps Context | Notification Context (delivery) | Event-driven HTTP callbacks |
| WebhookDelivery | GitOps Context |  | Webhook delivery attempt tracking |
| Tenant | Multi-tenancy Context | IAM Context (user mapping) | Tenant isolation and resource management |
| NamespaceQuota | Multi-tenancy Context | Registry Context (enforcement) | Resource quota management per namespace |

## Domain Events

### Registry Context Events
- **RepositoryCreated**: Repository created with metadata
- **RepositoryDeleted**: Repository deleted (soft or hard)
- **ArtifactPushed**: New artifact pushed to registry
- **ArtifactPulled**: Artifact pulled from registry
- **TagCreated**: New tag created pointing to artifact
- **TagDeleted**: Tag deleted
- **BlobStored**: Blob stored in blob storage
- **BlobDeleted**: Blob deleted during garbage collection
- **GarbageCollectionCompleted**: Garbage collection cycle finished

### IAM Context Events
- **UserCreated**: New user registered
- **UserUpdated**: User profile updated
- **UserDeleted**: User account deleted
- **UserDisabled**: User account disabled
- **UserEnabled**: User account re-enabled
- **UserLoggedIn**: User successfully authenticated
- **UserLoggedOut**: User session ended
- **GroupCreated**: New group created
- **GroupUpdated**: Group information updated
- **GroupDeleted**: Group deleted
- **UserAddedToGroup**: User added to group
- **UserRemovedFromGroup**: User removed from group
- **RoleCreated**: New role defined
- **RoleUpdated**: Role definition updated
- **RoleDeleted**: Role deleted
- **RoleAssignedToUser**: Role assigned to user
- **RoleAssignedToGroup**: Role assigned to group
- **RoleRevokedFromUser**: Role revoked from user
- **RoleRevokedFromGroup**: Role revoked from group
- **ClientCreated**: New client application registered
- **ClientUpdated**: Client information updated
- **ClientDeleted**: Client application deleted
- **ClientSecretRotated**: Client secret rotated

### Trust and Security Context Events
- **VulnerabilityDiscovered**: New vulnerability found in artifact
- **VulnerabilityUpdated**: Vulnerability information updated
- **VulnerabilityRemediated**: Vulnerability fixed in artifact
- **SBOMGenerated**: SBOM created for artifact
- **SBOMUpdated**: SBOM regenerated for artifact
- **SignatureVerified**: Signature verified as valid
- **SignatureFailed**: Signature verification failed
- **TrustScoreCalculated**: Trust score calculated for artifact
- **TrustScoreUpdated**: Trust score updated for artifact
- **PolicyEvaluated**: Policy evaluated against artifact
- **PolicyViolation**: Artifact violates policy
- **PolicyCompliance**: Artifact complies with policy
- **QuarantineApplied**: Artifact placed in quarantine
- **QuarantineRemoved**: Artifact removed from quarantine

### Observability Context Events
- **MetricThresholdExceeded**: Metric value crosses configured threshold
- **MetricThresholdCleared**: Metric value returns below threshold
- **LogRetentionExpired**: Log entries deleted due to retention policy
- **TraceCompleted**: Distributed trace completed
- **AlertFired**: Alert rule condition met, alert activated
- **AlertResolved**: Alert condition cleared, alert deactivated
- **AlertNotificationSent**: Alert notification sent to channels
- **AlertNotificationFailed**: Failed to send alert notification

### GitOps and Automation Context Events
- **WebhookCreated**: New webhook subscription created
- **WebhookUpdated**: Webhook subscription updated
- **WebhookDeleted**: Webhook subscription deleted
- **WebhookTriggered**: Webhook triggered by event
- **WebhookDeliverySuccess**: Webhook delivery successful
- **WebhookDeliveryFailed**: Webhook delivery failed after retries
- **WebhookDeliveryRetried**: Webhook delivery attempt retried
- **EventPublished**: System event published to event bus
- **EventConsumed**: System event consumed by subscriber

### Multi-tenancy Context Events
- **TenantCreated**: New tenant created
- **TenantUpdated**: Tenant information updated
- **TenantDeleted**: Tenant deleted
- **NamespaceQuotaCreated**: New namespace quota defined
- **NamespaceQuotaUpdated**: Namespace quota updated
- **NamespaceQuotaDeleted**: Namespace quota deleted
- **QuotaExceeded**: Namespace resource usage exceeds quota
- **QuotaCleared**: Namespace resource usage returns below quota
- **UserAssignedToTenant**: User assigned to tenant
- **UserRemovedFromTenant**: User removed from tenant

## Cross-Cutting Concerns

### Auditing
All privileged operations across domains generate audit events:
- **AuditEvent** Entity:
  - ID: UUID
  - Timestamp: Timestamp
  - Actor: UUID (User or ServiceAccount)
  - Action: String (operation performed)
  - ResourceType: String (type of resource affected)
  - ResourceID: UUID (identifier of resource affected)
  - ResourceName: String (name of resource affected)
  - Outcome: Enum (success, failure)
  - FailureReason: String (nullable)
  - IPAddress: String (actor IP address)
  - UserAgent: String (actor user agent)
  - RequestID: String (correlation ID)
  - Changes: JSON (before/after values for updates)

### Multi-tenancy Isolation
- **Data Isolation**: Each tenant's data is isolated via tenant_id foreign key
- **Namespace Isolation**: Kubernetes namespaces provide runtime isolation
- **Resource Isolation**: Quotas prevent resource starvation between tenants
- **Network Isolation**: NetworkPolicies restrict cross-namespace communication
- **Authentication Isolation**: Users and groups are tenant-scoped

### Security Propagation
- **Authentication Tokens**: JWT tokens carry tenant and user information
- **Authorization Context**: Policies evaluated with tenant and user context
- **Audit Context**: Audit events include tenant and user information
- **Observability Context**: Metrics, logs, and traces include tenant tags

### Consistency Patterns
- **Eventual Consistency**: Most cross-domain updates use event-driven eventual consistency
- **Sagas**: Long-running transactions use saga pattern with compensating actions
- **Read-Through Caching**: Frequently accessed data cached with read-through pattern
- **Write-Behind Caching**: Updates written to cache first, then persisted asynchronously

## Data Access Patterns

### Repository Pattern
Each domain implements repository interfaces for data access:
- **Read Operations**: GetByID, GetByCriteria, List, Count
- **Write Operations**: Create, Update, Delete, Upsert
- **Transactional Operations**: Unit of work with commit/rollback
- **Query Specification**: Criteria objects for complex queries
- **Pagination**: Offset/limit or cursor-based pagination

### Caching Strategy
- **L1 Cache**: Local in-memory cache (per service instance)
- **L2 Cache**: Redis shared cache (distributed)
- **L3 Cache**: HTTP/CDN cache (for static assets and public blobs)
- **Cache Invalidation**: Event-driven invalidation or TTL-based
- **Cache Warming**: Pre-loading of frequently accessed data

### Search and Filtering
- **Full-Text Search**: Elasticsearch for log and audit trail search
- **Filtered Queries**: Database indexes for common query patterns
- **Faceted Navigation**: Pre-computed counts for UI filtering
- **Search Relevance**: TF-IDF or BM25 scoring for text search

### Reporting and Analytics
- **Real-time Dashboards**: Prometheus + Grafana for operational metrics
- **Historical Analytics**: Data warehouse for trend analysis
- **Ad-hoc Queries**: SQL interface for exploratory analysis
- **Export Capabilities**: CSV, JSON, PDF export for reports

## Evolution and Versioning

### Schema Evolution
- **Backward Compatible Changes**: 
  - Add new tables/columns with nullable defaults
  - Add new enum values
  - Add new optional fields to JSON documents
- **Breaking Changes**:
  - Require versioned APIs
  - Provide migration scripts
  - Maintain dual-write during transition period
- **Schema Versioning**: 
  - Migration tracking table
  - Automated migration execution
  - Rollback capability

### API Versioning
- **URI Versioning**: /api/v1/, /api/v2/
- **Header Versioning**: Accept: application/vnd.kyros.v2+json
- **Deprecation Policy**: 
  - Deprecation notice in API responses
  - Sunset date provided
  - Migration guide available
- **Semantic Versioning**: MAJOR.MINOR.PATCH for breaking/feature/fix

### Data Migration
- **Migration Framework**: 
  - Versioned migration scripts
  - Up/down functions for each migration
  - Transactional execution where possible
  - Irreversible migrations marked explicitly
- **Data Backfilling**: 
  - Background jobs for data enrichment
  - Batch processing for large datasets
  - Progress reporting and resumability
- **Rollback Procedures**: 
  - Tested down migrations
  - Backup/restore procedures
  - Point-in-time recovery capabilities

## Diagrams Reference
See [MERMAID.md](MERMAID.md) for detailed Mermaid diagrams including:
- Domain Model Entity Relationship Diagrams
- Aggregate Boundaries
- Bounded Context Maps
- Domain Event Flow
- Cross-Cutting Concerns