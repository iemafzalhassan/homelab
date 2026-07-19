# Kyros Registry Implementation

## Overview
Kyros implements a fully OCI-compliant container registry based on the cncf/distribution v3 registry engine, extended with additional features for security, observability, and enterprise capabilities. This document explains the registry implementation details, architecture, and extensions.

## Core Registry Engine

### Base Implementation
Kyros uses [cncf/distribution](https://github.com/distribution/distribution) v3 as the foundation for its registry functionality:
- **OCI Distribution Spec Compliance**: Full implementation of OCI Distribution Specification v1.0.1
- **OCI Image Spec Compliance**: Support for OCI Image Specification v1.0.1
- **Docker Registry HTTP API V2**: Backward compatibility with Docker Engine
- **Storage Drivers**: Pluggable backend for blob and metadata storage
- **Authentication Framework**: Extensible authentication and authorization system
- **Event Notifications**: Built-in webhook and notification system
- **Garbage Collection**: Built-in blob garbage collection capabilities

### Registry Service Architecture
The Kyros registry service extends the base distribution engine with:

#### 1. Storage Layer Abstraction
```go
// Registry storage interface extending distribution's storage.StorageDriver
type StorageDriver interface {
    distribution.StorageDriver
    
    // Extended methods for Kyros-specific features
    GetArtifactMetadata(ctx context.Context, digest string) (*ArtifactMetadata, error)
    StoreArtifactMetadata(ctx context.Context, metadata *ArtifactMetadata) error
    GetTrustScore(ctx context.Context, digest string) (*TrustScore, error)
    StoreTrustScore(ctx context.Context, digest string, score *TrustScore) error
    ListArtifactsByNamespace(ctx context.Context, namespaceID string) ([]*ArtifactListing, error)
    GetNamespaceUsage(ctx context.Context, namespaceID string) (*NamespaceUsage, error)
}
```

#### 2. Metadata Management
While distribution stores minimal metadata in its storage backend, Kyros extends this with:
- **PostgreSQL Backend**: Primary store for rich metadata (users, repositories, trust scores, etc.)
- **Hybrid Approach**: Blobs stored in object storage (MinIO/S3), metadata in PostgreSQL
- **Cached Layer**: Redis cache for frequently accessed metadata
- **Event Synchronization**: Ensures consistency between distribution metadata and PostgreSQL

#### 3. Authentication Extension
Kyros replaces distribution's simple token authentication with:
- **Keycloak Integration**: OpenID Connect bearer token validation
- **Service Account Support**: Dedicated service accounts for automated systems
- **Fine-Grained Permissions**: Repository-scoped and namespace-scoped permissions
- **Session Management**: Web session support for browser-based UI
- **MFA Support**: Multi-factor authentication for sensitive operations

#### 4. Authorization Extension
Kyros enhances distribution's basic authorization with:
- **RBAC Engine**: Role-based access control with custom roles
- **Policy Integration**: OPA integration for complex authorization decisions
- **Resource Scoping**: Permissions scoped to specific repositories/namespaces
- **Dynamic Policies**: Runtime-updatable authorization policies
- **Audit Logging**: Comprehensive logging of all authorization decisions

## OCI Registry Concepts in Kyros

### 1. Blobs and Layers
- **Content Addressable Storage**: Each blob identified by SHA256 digest
- **Immutable Storage**: Blobs never modified after upload
- **Deduplication**: Automatic deduplication of identical layers across images
- **Storage Backend**: Configurable (MinIO, S3, Azure Blob, GCS, etc.)
- **Chunked Upload**: Support for resumable chunked uploads for large layers
- **Garbage Collection**: Periodic cleanup of unreferenced blobs

### 2. Manifests and Artifacts
- **OCI Manifest Format**: Support for OCI Image Index and Image Manifest
- **Manifest Schema Versions**: Support for schemaVersion 2 and 2 (multi-arch)
- **Artifact References**: Support for OCI Artifacts beyond container images
- **Referrers API**: Implementation of OCI Referrers API for discovering artifacts that reference others
- **Subject References**: Tracking of what references each artifact (for garbage collection)

### 3. Tags and References
- **Mutable Tags**: Human-readable pointers to immutable manifests
- **Tag Immutability Option**: Repositories can enforce immutable tags
- **Tag Retention Policies**: Automatic cleanup of old tags based on policies
- **Referrers Tracking**: Tracking which artifacts reference a given manifest
- **Garbage Collection Safety**: Ensuring referenced manifests are not collected

### 4. Namespaces and Repositories
- **Hierarchical Organization**: Namespaces contain repositories
- **Access Control**: Permissions can be set at namespace or repository level
- **Visibility Settings**: Public, private, or protected visibility levels
- **Automatic Creation**: Repositories can be created on first push (configurable)
- **Namespace Quotas**: Storage and rate limits per namespace

## Registry API Extensions

Kyros extends the standard Docker Registry HTTP API V2 with additional endpoints for management and enterprise features:

### Standard Registry Endpoints (Unchanged)
```http
GET /v2/                          # API version check
GET /v2/_catalog                  # List repositories
GET /v2/<name>/tags/list          # List tags in repository
HEAD /v2/<name>/blobs/<digest>    # Check blob existence
GET /v2/<name>/blobs/<digest>     # Download blob
PUT /v2/<name>/blobs/uploads/     # Initiate blob upload
PATCH /v2/<name>/blobs/uploads/<uuid>  # Upload chunk
PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>  # Complete upload
POST /v2/<name>/manifests/<reference>  # Push manifest
GET /v2/<name>/manifests/<reference>   # Pull manifest
HEAD /v2/<name>/manifests/<reference>  # Check manifest existence
DELETE /v2/<name>/manifests/<reference> # Delete manifest
```

### Kyros Management Endpoints
```http
# Repository Management
POST /api/v1/repositories          # Create repository
GET /api/v1/repositories/{id}      # Get repository details
PATCH /api/v1/repositories/{id}    # Update repository
DELETE /api/v1/repositories/{id}   # Delete repository

# Trust and Security
GET /api/v1/trust/scores           # List trust scores
GET /api/v1/trust/scores/{digest}  # Get trust score for artifact
POST /api/v1/trust/scores/{digest}/recalculate  # Recalculate trust score
GET /api/v1/trust/vulnerabilities  # List vulnerabilities
GET /api/v1/trust/sboms            # List SBOMs
GET /api/v1/trust/signatures       # List signatures

# Policy Management
GET /api/v1/trust/policies         # List policies
POST /api/v1/trust/policies        # Create policy
GET /api/v1/trust/policies/{id}    # Get policy details
PATCH /api/v1/trust/policies/{id}  # Update policy
DELETE /api/v1/trust/policies/{id} # Delete policy
POST /api/v1/trust/policies/{id}/evaluate  # Evaluate policy

# Webhook Management
GET /api/v1/webhooks               # List webhooks
POST /api/v1/webhooks              # Create webhook
GET /api/v1/webhooks/{id}          # Get webhook details
PATCH /api/v1/webhooks/{id}        # Update webhook
DELETE /api/v1/webhooks/{id}       # Delete webhook
GET /api/v1/webhooks/{id}/deliveries  # List delivery attempts

# Namespace Management
GET /api/v1/namespaces             # List namespaces
POST /api/v1/namespaces            # Create namespace
GET /api/v1/namespaces/{id}        # Get namespace details
PATCH /api/v1/namespaces/{id}      # Update namespace
DELETE /api/v1/namespaces/{id}     # Delete namespace
GET /api/v1/namespaces/{id}/repositories  # List repositories in namespace

# Multi-tenancy
GET /api/v1/tenants                # List tenants
POST /api/v1/tenants               # Create tenant
GET /api/v1/tenants/{id}           # Get tenant details
PATCH /api/v1/tenants/{id}         # Update tenant
DELETE /api/v1/tenants/{id}        # Delete tenant
GET /api/v1/namespace-quotas       # List namespace quotas
POST /api/v1/namespace-quotas      # Create namespace quota
```

## Storage Architecture

### Blob Storage
Kyros stores blob data (layers, configs) in object storage:
- **Primary Backend**: MinIO (S3-compatible) for self-hosted deployments
- **Cloud Options**: AWS S3, Google Cloud Storage, Azure Blob Storage
- **Configuration**: Pluggable via storage driver interface
- **Encryption**: Server-side encryption (SSE-S3, SSE-KMS) or client-side encryption
- **Versioning**: Object versioning enabled for accidental deletion protection
- **Lifecycle Rules**: Automatic transition to cheaper storage classes
- **Replication**: Cross-region replication for disaster recovery
- **Performance**: CDN caching for frequently accessed blobs (with security considerations)

### Metadata Storage
Kyros stores metadata in PostgreSQL:
- **Primary Store**: Users, groups, roles, permissions, repositories, trusts
- **Relationships**: Foreign keys maintain data integrity
- **Indexes**: Optimized for common query patterns
- **Partitioning**: Time-based partitioning for high-volume tables (events, logs)
- **Connection Pooling**: PgBouncer for efficient connection management
- **Read Replicas**: For scaling read-heavy operations
- **Backup Strategy**: Regular logical and physical backups

### Caching Layer
Kyros uses Redis for distributed caching:
- **User Information**: Cached user profiles and group memberships
- **Token Validation**: Cached JWT validation results (short TTL)
- **Repository Metadata**: Frequently accessed repository information
- **Trust Scores**: Cached trust score calculations
- **Policy Decisions**: Cached OPA policy evaluations
- **Rate Limiting**: Distributed rate limiting counters
- **Session Storage**: Web session data for browser UI
- **Configuration**: Feature flags and dynamic configuration

## Security Features

### Transport Security
- **TLS Everywhere**: All service-to-service and client-to-service communication encrypted
- **Certificate Management**: Automated certificate rotation via cert-manager
- **Cipher Suites**: Modern, secure cipher suites only
- **Protocol Versions**: TLS 1.2 and 1.3 only
- **Perfect Forward Secrecy**: Ephemeral key exchanges

### Authentication Security
- **Short-Lived Tokens**: Access tokens valid for 15 minutes
- **Refresh Token Rotation**: Refresh tokens rotated on each use
- **Token Revocation**: Immediate revocation capability
- **MFA Support**: Optional multi-factor authentication
- **Password Policies**: Strong password requirements via Keycloak
- **Brute Force Protection**: Rate limiting on authentication endpoints
- **Credential Stuffing Defense**: CAPTCHA and login throttling

### Authorization Security
- **Principle of Least Privilege**: Minimal permissions granted by default
- **Resource Scoping**: Permissions limited to specific resources
- **Policy Enforcement**: Centralized policy decision points
- **Audit Logging**: All authorization decisions logged
- **Default Deny**: Default deny posture for new resources
- **Segregation of Duties**: Separation between administrative and operational roles

### Data Security
- **Encryption at Rest**:
  - Database: Transparent encryption or application-level encryption
  - Object Storage: Server-side encryption (SSE-S3, SSE-KMS)
  - Backups: Encrypted storage
  - Secrets: External secret manager (Vault, Azure Key Vault)
  
- **Secrets Management**:
  - Database credentials: External secret manager
  - Service credentials: Kubernetes secrets or external secret manager
  - Encryption keys: Hardware Security Module (HSM) or cloud KMS
  
- **Input Validation**:
  - All API inputs validated and sanitized
  - Manifest schema validation
  - Tag name validation (prevent injection)
  - Blob upload limits (prevent DoS)

### Network Security
- **Network Policies**: Kubernetes NetworkPolicies restricting service communication
- **Service Mesh**: Optional Istio/Linkerd for mTLS and traffic policies
- **Ingress Control**: Traefik with rate limiting, WAF, and access controls
- **Private Networks**: Services deployed in private subnets
- **Egress Control**: Restricted outbound communication where possible
- **DDoS Protection**: Rate limiting and connection limits at ingress

### Audit and Compliance
- **Comprehensive Logging**: All privileged operations logged to immutable storage
- **Log Integrity**: Cryptographic hashing and chaining of log entries
- **Log Retention**: Configurable retention periods with archival
- **Compliance Reporting**: Built-in reports for SOC 2, ISO 27001, GDPR, HIPAA
- **Vulnerability Scanning**: Regular scanning of Kyros container images
- **Penetration Testing**: Regular third-party security assessments
- **Security Headers**: CSP, HSTS, X-Frame-Options, etc.

## Observability Features

### Metrics
Kyros exports Prometheus metrics for monitoring:
- **Registry Metrics**:
  - `registry_http_requests_total`: HTTP requests by method, endpoint, status
  - `registry_http_request_duration_seconds`: Request latency distribution
  - `registry_uploaded_bytes_total`: Bytes uploaded by repository
  - `registry_downloaded_bytes_total`: Bytes downloaded by repository
  - `registry_blob_storage_bytes_total`: Total blob storage used
  - `registry_blob_count_total`: Total blob count
  - `registry_garbage_collected_blobs_total`: Blobs collected by GC
  - `registry_trust_score_calculations_total`: Trust score calculations
  - `registry_policy_evaluations_total`: Policy evaluations
  - `registry_webhook_deliveries_total`: Webhook delivery attempts
  
- **Business Metrics**:
  - `registry_repositories_total`: Total repositories
  - `registry_namespaces_total`: Total namespaces
  - `registry_users_total`: Total users
  - `registry_artifacts_pushed_total`: Artifacts pushed
  - `registry_artifacts_pulled_total`: Artifacts pulled
  
- **Service Metrics**:
  - `go_gc_duration_seconds`: Go GC pause distribution
  - `go_threads`: Number of OS threads created
  - `process_resident_memory_bytes`: Resident memory size
  - `process_cpu_seconds_total`: Total user and system CPU time

### Logging
Kyros implements structured logging:
- **Format**: JSON-formatted logs for easy parsing
- **Fields**: Timestamp, level, message, service, traceID, spanID, attributes
- **Levels**: debug, info, warn, error, fatal
- **Sampling**: Adaptive sampling for high-volume debug logs
- **Redaction**: Automatic redaction of sensitive information (tokens, passwords)
- **Routing**: Logs sent to Loki or similar log aggregation system
- **Retention**: Configurable retention with automatic cleanup

### Tracing
Kyros implements distributed tracing:
- **Instrumentation**: OpenTelemetry SDK in all services
- **Context Propagation**: W3C TraceContext headers
- **Span Creation**: Automatic span creation for incoming/outgoing requests
- **Custom Spans**: Manual spans for business logic operations
- **Attributes**: Rich attributes for filtering and analysis
- **Events**: Span events for significant occurrences
- **Links**: Span links for causal relationships
- **Status**: Span status (ok, error) with optional message
- **Collection**: OpenTelemetry Collector agents
- **Storage**: Tempo trace database
- **Visualization**: Grafana TraceQL integration
- **Sampling**: Adaptive sampling based on traffic volume and error rates

### Health Checks
Kyros provides comprehensive health endpoints:
- **Liveness Probe**: `/healthz` - basic service liveness
- **Readiness Probe**: `/readyz` - service readiness to serve traffic
- **Startup Probe**: `/startupz` - service startup completion
- **Detailed Health**: `/healthz/detailed` - component-level health status
- **Dependency Checks**: Verifies connectivity to downstream services
- **Resource Checks**: Memory, disk space, file descriptor usage
- **Dependency Graph**: Shows health of all dependencies

## Garbage Collection

Kyros implements robust garbage collection for storage efficiency:

### Blob Garbage Collection
- **Trigger**: Manual initiation or scheduled (cron-based)
- **Algorithm**: Mark-and-sweep algorithm
  1. **Mark Phase**: Identify all referenced blobs via manifest traversal
  2. **Sweep Phase**: Delete unreferenced blobs from storage
- **Concurrent Safe**: Safe to run while registry is serving traffic
- **Progress Reporting**: Progress metrics and logging
- **Interruptible**: Can be paused and resumed
- **Dry Run Mode**: Preview what would be deleted without actual deletion
- **Configuration**: Configurable thresholds (age, size, percentage)

### Manifest Garbage Collection
- **Tag Expiration**: Automatic deletion of old tags based on policies
- **Manifest Retention**: Retention policies for manifests themselves
- **Referrers Tracking**: Using OCI Referrers API to track what references manifests
- **Legal Hold**: Ability to place manifests under legal hold to prevent deletion
- **Artifact Retention**: Different policies for different artifact types (images vs. Helm charts)

### Storage Optimization
- **Layer Deduplication**: Automatic deduplication of identical layers
- **Content-Defined Chunking**: For future chunk-based storage optimization
- **Compression**: Optional compression of stored blobs
- **Storage Tiering**: Automatic movement of infrequently accessed blobs to cheaper storage
- **Quota Enforcement**: Per-namespace and per-repository storage quotas

## API Extensions and Custom Endpoints

### Extended Registry API
Kyros extends the standard registry API with additional capabilities:

#### Blob Upload Enhancements
```http
# Resumable chunked upload with progress tracking
PATCH /v2/<name>/blobs/uploads/<uuid>?range=<bytes>-<bytes>
Response: 202 Accepted with Location header and Range header showing uploaded range

# Upload verification
GET /v2/<name>/blobs/uploads/<uuid>?verify
Response: 200 OK with JSON containing uploaded ranges and total size
```

#### Manifest Enhancements
```http
# Manifest validation without storage
POST /v2/<name>/manifests/<reference>?_validate
Response: 201 Created if manifest is valid, 400 Bad Request if invalid

# Manifest blob sum verification
HEAD /v2/<name>/manifests/<reference>?_checksum
Response: 200 OK with Docker-Content-Digest header if layers present
```

#### Referrers API (OCI Referrers)
```http
# Get artifacts that reference a given manifest
GET /v2/<name>/referrers/<digest>
Response: 200 OK with JSON array of referencing manifests

# Get referrers with specific properties
GET /v2/<name>/referrers/<digest>?kind=<artifact-type>&size=<n>
```

### Management API Endpoints
As detailed in the API Design document, Kyros provides comprehensive management APIs for:
- Repository and namespace management
- Trust score and security features
- Policy management
- Webhook configuration
- Multi-tenancy controls
- Audit and compliance reporting

## Integration Points

### 1. Authentication Integration
Kyros registry integrates with the authentication service:
- **Token Validation**: Delegates JWT validation to auth service
- **Permission Checks**: Calls auth service for authorization decisions
- **User Information**: Retrieves user info from auth service for audit logs
- **Session Validation**: Validates web sessions for browser UI
- **Service Account Handling**: Special handling for service-to-service requests

### 2. Trust Score Integration
Kyros registry interacts with the trust score service:
- **Pre-Push Scanning**: Optional scanning before allowing push (configurable)
- **Post-Push Scanning**: Automatic scanning after push completion
- **Trust Score Storage**: Stores trust scores in PostgreSQL via auth service
- **Trust Score Retrieval**: Provides trust scores via management API
- **Blocking Pushes**: Can block pushes based on trust score policies
- **Event Publishing**: Publishes events when trust scores are updated

### 3. Webhook Integration
Kyros registry extends distribution's webhook capabilities:
- **Extended Events**: Additional event types (trust.score.updated, policy.violated, etc.)
- **Enhanced Payloads**: Richer payloads with trust score, vulnerability data, etc.
- **Retry Logic**: Improved retry with exponential backoff and jitter
- **Dead Letter Queues**: Failed deliveries sent to DLQ for later processing
- **Event Filtering**: Server-side event filtering to reduce unnecessary deliveries
- **Payload Transformation**: Optional transformation of webhook payloads
- **Signature Verification**: HMAC-SHA256 verification of webhook payloads
- **Delivery Guarantees**: At-least-once delivery guarantee
- **Circuit Breaker**: Circuit breaker pattern for unhealthy endpoints

### 4. Policy Integration
Kyros registry integrates with the policy engine:
- **Admission Control**: Policies evaluated during manifest upload
- **Pre-Receive Hooks**: Policies evaluated before storing blobs
- **Post-Receive Hooks**: Policies evaluated after manifest upload
- **Sync Policies**: Policies evaluated during mirroring/sync operations
- **Garbage Collection Policies**: Policies influencing GC decisions
- **Quota Policies**: Policies affecting quota calculations
- **Audit Integration**: Policy decisions logged to audit trail

### 5. Observability Integration
Kyros registry exports comprehensive observability data:
- **Metrics**: Prometheus metrics for all registry operations
- **Logging**: Structured logs with trace IDs for correlation
- **Tracing**: OpenTelemetry instrumentation for distributed tracing
- **Health Checks**: Liveness, readiness, and startup probes
- **Profiling**: CPU, heap, block, and mutex profiling endpoints
- **Debug Endpoints**: Conditional debug endpoints for troubleshooting

## Deployment and Configuration

### Deployment Options
Kyros registry can be deployed in various configurations:

#### 1. Standalone Deployment
- **Single Binary**: All registry functionality in one process
- **Embedded Dependencies**: Uses embedded SQLite for simple deployments
- **Object Storage**: Configured external blob storage
- **Use Cases**: Development, testing, small edge deployments

#### 2. Microservices Deployment
- **Separate Services**: Registry service separated from auth, trustscore, etc.
- **External Dependencies**: Uses external PostgreSQL, Redis, MinIO
- **Kubernetes Native**: Designed for Kubernetes deployment
- **Helm Chart**: Available for easy installation
- **Operator Pattern**: Kyros Operator for lifecycle management
- **Use Cases**: Production deployments, scalable installations

#### 3. High Availability Deployment
- **Multiple Replicas**: Multiple registry service replicas behind load balancer
- **Shared Storage**: All replicas share same object storage and database
- **Database Clustering**: PostgreSQL with replication and failover
- **Redis Clustering**: Redis cluster for shared caching
- **Object Storage**: Multi-region or replicated blob storage
- **Use Cases**: Enterprise production, mission-critical deployments

### Configuration Management
Kyros registry configuration via:
- **Environment Variables**: Primary configuration method
- **Configuration Files**: YAML file for complex configurations
- **Command-Line Flags**: Override specific settings
- **Secrets Management**: External secret store for sensitive values
- **Dynamic Configuration**: Runtime configuration changes without restart
- **Feature Flags**: Gradual rollout of new features
- **Validation**: Configuration validation on startup

### Key Configuration Areas
#### Storage Configuration
```env
REGISTRY_STORAGE=blobs
REGISTRY_STORAGE_BLOBS=blobs
REGISTRY_STORAGE_BLOBS_ROOTDIRECTORY=/var/lib/registry
# For object storage
REGISTRY_STORAGE=s3
REGISTRY_STORAGE_S3_ACCESSKEY=...
REGISTRY_STORAGE_S3_SECRETKEY=...
REGISTRY_STORAGE_S3_REGION=us-east-1
REGISTRY_STORAGE_S3_BUCKET=my-registry-bucket
REGISTRY_STORAGE_S3_ENCRYPT=true
REGISTRY_STORAGE_S3_SECURE=true
REGISTRY_STORAGE_S3_V4AUTH=true
REGISTRY_STORAGE_S3_CHUNKSIZE=5242880  # 5MiB
```

#### Authentication Configuration
```env
AUTH_TYPE=keycloak
AUTH_KEYCLOAK_URL=https://keycloak.example.com
AUTH_KEYCLOAK_REALM=kyros
AUTH_KEYCLOAK_CLIENT_ID=kyros-registry
AUTH_KEYCLOAK_CLIENT_SECRET=${KEYCLOAK_CLIENT_SECRET}
AUTH_TOKEN_VALIDATION_LOCAL=true
AUTH_TOKEN_INTROSPECTION_CACHE_TTL=30
AUTH_USER_INFO_CACHE_TTL=300
```

#### Authorization Configuration
```env
AUTHORIZATION_TYPE=rbac+opa
AUTHORIZATION_POLICY_PATH=/policies
AUTHORIZATION_SUPER_USER=admin
AUTHORIZATION_ALWAYS_ALLOW=false
```

#### Trust Score Configuration
```env
TRUST_SCORE_ENABLED=true
TRUST_SCORE_SERVICE_URL=https://trustscore.kyros.example.com
TRUST_SCORE_AUTO_CALCULATE=true
TRUST_SCORE_MINIMUM_SCORE=0.6
TRUST_SCORE_BLOCK_BELOW_SCORE=0.3
TRUST_SCORE_POLICY_ID=default-trust-policy
```

#### Webhook Configuration
```env
WEBHOOK_ENABLED=true
WEBHOOK_ENDPOINT=webhook.kyros.example.com
WEBHOOK_MAXRETRIES=10
WEBHOOK_RETRYINTERVAL=60
WEBHOOK_MAXPAYLOAD=4194304  # 4MiB
WEBHOOK_HEADERSIZE=8192
WEBHOOK_TIMEOUT=30
WEBHOOK_LEASEDURATION=10m
```

#### Logging Configuration
```env
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
LOG_FIELDS=service,hostname,traceID,spanID
LOG_ACCESSLOG=true
LOG_ACCESSLOG_IGNOREDCODES=
LOG_ACCESSLOG_RESPONSEBODYLIMIT=0
```

#### HTTP Configuration
```env
HTTP_ADDR=:5000
HTTP_NET=tcp
HTTP_HOST=https://registry.kyros.example.com
HTTP_TLS_CERTIFICATE=/certs/domain.crt
HTTP_TLS_KEY=/certs/domain.key
HTTP_TLS_LETSENCRYPT_CACHE=/letsencrypt
HTTP_TLS_LETSENCRYPT_RESOLVERS=myresolver
HTTP_HEADERSERVER="Kyros Registry"
HTTP_HEADERX-CONTENT-TYPE-OPTIONS=nosniff
HTTP_HEADERX-FRAME-OPTIONS=DENY
HTTP_HEADERXXSSPROTECTION=1; mode=block
```

## Scaling and Performance

### Horizontal Scaling
Kyros registry scales horizontally:
- **Stateless Design**: Registry service instances are stateless
- **Shared Storage**: All instances share same blob storage and database
- **Load Balancing**: Traffic distributed via TCP/HTTP load balancer
- **Session Affinity**: Not required for API traffic (stateless JWT)
- **WebSocket Connections**: Require session affinity for real-time UI updates
- **Database Connections**: Connection pooling prevents connection exhaustion
- **Redis Connections**: Shared Redis cluster for distributed caching
- **Object Storage**: Designed for high concurrent access

### Performance Optimization
#### Upload Optimization
- **Parallel Chunk Upload**: Multiple chunks can be uploaded in parallel
- **Chunk Size Configuration**: Configurable chunk size for optimal throughput
- **Resume Support**: Interrupted uploads can be resumed from last checkpoint
- **Expect Header**: Support for `Expect: 100-continue` to reduce bandwidth
- **Content-Length Validation**: Early rejection of oversized uploads

#### Download Optimization
- **Range Requests**: Support for HTTP range requests for partial downloads
- **ETag Support**: ETag headers for cache validation
- **Cache-Control Headers**: Configurable cache control for proxies
- **Content-Disposition**: Proper content disposition headers for downloads
- **Blast Radius Limiting**: Rate limiting per IP to prevent bandwidth exhaustion

#### Metadata Optimization
- **Database Indexing**: Optimized indexes for common query patterns
- **Query Caching**: Frequently accessed metadata cached in Redis
- **Connection Pooling**: PgBouncer for efficient database connections
- **Read Replicas**: Read-only queries directed to database replicas
- **Batch Operations**: Bulk operations for efficiency
- **Asynchronous Processing**: Non-critical operations processed asynchronously

#### Garbage Collection Optimization
- **Incremental GC**: Garbage collection can run incrementally
- **Priority Scheduling**: Lower priority during peak hours
- **Resource Throttling**: CPU and I/O throttling during GC
- **Progress Tracking**: Detailed progress reporting
- **Concurrent Safe**: Safe to run while serving traffic
- **Storage Driver Optimization**: Optimized for specific storage backends

### Resource Requirements
#### Minimum Requirements
- **CPU**: 2 cores
- **Memory**: 4 GB RAM
- **Storage**: 
  - Blob Storage: Depends on image count and size
  - Metadata: ~100 MB for basic deployment
  - Redis: ~100 MB for caching
- **Network**: 1 Gbps recommended

#### Recommended Production
- **CPU**: 8-32 cores (depending on load)
- **Memory**: 16-64 GB RAM
- **Storage**:
  - Blob Storage: Scalable object storage (S3/minIO)
  - Metadata: PostgreSQL with adequate IOPS
  - Redis: Redis cluster for HA
- **Network**: 10 Gbps recommended for high-throughput scenarios

#### Scaling Guidelines
- **Requests Per Second**: 
  - Read-heavy (pulls): 1000+ RPM per core
  - Write-heavy (pushes): 100-500 RPM per core (depending on layer size)
  - Mixed workload: 200-800 RPM per core
- **Concurrent Connections**: 1000+ per instance with proper tuning
- **Storage Throughput**: 
  - Blob Upload: Limited by storage backend and network
  - Blob Download: Limited by storage backend and network
  - Metadata: Limited by database connection pool and query complexity

## Data Durability and Consistency

### Blob Storage Durability
- **Object Storage Durability**: 11 nines (99.999999999%) for S3/MinIO with versioning
- **Geo-Replication**: Cross-Region Replication**: For disaster recovery
- **Versioning**: Object versioning enabled to protect against accidental deletion
- **Lifecycle Policies**: Automatic transition to glacier/deep archive for old data
- **Checksum Validation**: Automatic checksum validation on read/write
- **Self-Healing**: Automatic repair of corrupted objects (where supported)

### Metadata Durability
- **Database Replication**: PostgreSQL streaming replication for HA
- **Backup Strategy**: 
  - Logical backups: Daily pg_dump
  - Physical backups: Continuous archiving with WAL-E/WAL-G
  - Point-in-Time Recovery: Capable of restoring to any point
- **Consistency Model**: Strong consistency for metadata operations
- **Transaction Isolation**: Serializable isolation for critical operations
- **Deadlock Detection**: Automatic deadlock detection and resolution
- **Connection Limits**: Configurable max connections to prevent overload

### Consistency Guarantees
#### Read-After-Write Consistency
- **Blob Uploads**: Immediate consistency for blob storage (object storage guarantee)
- **Manifest Uploads**: Immediate consistency for manifest storage
- **Metadata Updates**: Eventually consistent with short convergence time (<1 second)
- **Trust Score Updates**: Eventually consistent (typically <5 seconds)
- **Policy Updates**: Eventually consistent (typically <1 second)

#### Transaction Boundaries
- **Atomic Blob Upload**: Blob upload is atomic (either complete or not stored)
- **Atomic Manifest Upload**: Manifest upload is atomic
- **Multi-Blob Upload**: Each blob uploaded atomically, but manifest references all
- **Tag Creation**: Tag creation is atomic
- **Repository Creation**: Repository creation is atomic

#### Conflict Resolution
- **Last Write Wins**: For most metadata updates (timestamps, counters)
- **Merge Strategies**: For complex data structures (application-level merging)
- **Optimistic Concurrency**: Version vectors or timestamps for conflict detection
- **Manual Resolution**: Administrative interface for conflict resolution
- **Audit Trail**: Complete audit trail for all changes

## Failure Scenarios and Recovery

### Storage Failures
#### Blob Storage Unavailable
- **Detection**: Failed blob upload/download operations
- **Impact**: New pushes/pulls fail; existing cached layers may still work
- **Recovery**: 
  - Failover to secondary storage (if configured)
  - Queue uploads for retry when storage returns
  - Serve errors to clients with appropriate retry-after headers
  - Alert operators immediately

#### Metadata Database Unavailable
- **Detection**: Failed database queries/connection attempts
- **Impact**: 
  - New repository creation fails
  - Permission checks may fail (fallback to cached permissions)
  - Trust score storage fails
  - Webhook configuration changes fail
  - Existing operations continue with cached data where possible
- **Recovery**:
  - Failover to database replica (if configured)
  - Queue metadata updates for retry
  - Serve cached data where possible with degraded functionality
  - Alert operators immediately

#### Cache Unavailable
- **Detection**: Increased latency or fallback to slower data sources
- **Impact**: 
  - Increased latency for frequent operations
  - Higher database load
  - Reduced performance but continued operation
- **Recovery**:
  - Automatic recovery when cache returns
  - Temporary increase in database load acceptable
  - No data loss, only performance impact

### Network Failures
#### Service-to-Service Communication Failure
- **Detection**: Failed RPC calls or message publishing
- **Impact**: 
  - Degraded functionality depending on service
  - Potential fallback to cached data or default behavior
  - Increased latency due to retries and timeouts
- **Recovery**:
  - Circuit breaker pattern prevents cascading failures
  - Exponential backoff and jitter for retries
  - Fallback to cached data or default behavior where appropriate
  - Alert operators for persistent issues

#### Client-to-Service Communication Failure
- **Detection**: Connection timeouts, refused connections
- **Impact**: Clients unable to reach service
- **Recovery**:
  - Load balancer health checks remove unhealthy instances
  - DNS TTL controls client reconnection behavior
  - Clients implement retry logic with exponential backoff
  - Service scales to handle increased load from retries

### Software Failures
#### Service Crash
- **Detection**: Process exit, restart loops
- **Impact**: Temporary unavailability until restart
- **Recovery**:
  - Kubernetes restart policy automatically restarts failed containers
  - Health checks prevent traffic to unhealthy instances
  - Stateful data preserved in external storage
  - Crash logs collected for analysis
  - Alert operators for frequent crashes

#### Deadlock
- **Detection**: Increased latency, timeout errors, thread dumps showing blocked threads
- **Impact**: Gradual degradation leading to complete unavailability
- **Recovery**:
  - Database deadlock detection automatically resolves deadlocks
  - Application-level timeouts prevent indefinite blocking
  - Thread dump analysis for root cause identification
  - Code fixes for recurring deadlock patterns
  - Alert operators for persistent issues

#### Memory Leak
- **Detection**: Gradually increasing memory usage over time
- **Impact**: Eventually leads to OOM kills and service crashes
- **Recovery**:
  - Memory profiling to identify leak sources
  - Automatic restart based on memory usage thresholds
  - Resource limits in Kubernetes to prevent node exhaustion
  - Alert operators for memory usage trends
  - Code fixes for memory leak sources

### Data Corruption
#### Blob Corruption
- **Detection**: Checksum mismatches during download
- **Impact**: Corrupted layers leading to image pull failures
- **Recovery**:
  - Automatic retry from alternative source (if mirrored)
  - Mark blob as corrupted and quarantine
  - Trigger garbage collection to remove corrupted blob
  - Alert operators for investigation
  - Potential need to re-push affected images

#### Metadata Corruption
- **Detection**: Database constraint violations, failed queries
- **Impact**: Inconsistent state, failed operations
- **Recovery**:
  - Database backup and point-in-time recovery
  - Manual repair of corrupted records
  - Application-level validation and correction
  - Alert operators immediately
  - Potential need for manual reconciliation

### Disaster Recovery
#### Regional Outage
- **Detection**: Loss of connectivity to primary region
- **Impact**: Complete service outage in affected region
- **Recovery**:
  - Automated failover to secondary region (if active-active)
  - DNS failover with low TTL
  - Manual initiation of failover procedures
  - Data synchronization from backups
  - Validation of failed-over services
  - Notification of users and stakeholders

#### Data Loss Incident
- **Detection**: Missing or corrupted data beyond recovery capabilities
- **Impact**: Permanent data loss requiring restoration from backup
- **Recovery**:
  - Identify point-in-time for recovery
  - Restore from backup to temporary environment
  - Validate restored data integrity
  - Switch production to restored environment
  - Perform post-incident analysis
  - Implement preventive measures

## Compliance and Certifications

### Regulatory Compliance
Kyros registry is designed to support compliance with:
- **SOC 2 Type II**: Security, availability, processing integrity, confidentiality, privacy
- **ISO 27001**: Information security management systems
- **GDPR**: General Data Protection Regulation (data protection and privacy)
- **HIPAA**: Health Insurance Portability and Accountability Act (PHI protection)
- **PCI DSS**: Payment Card Industry Data Security Standard (cardholder data)
- **FedRAMP**: Federal Risk and Authorization Management Program (US government)
- **CCPA**: California Consumer Privacy Act (consumer privacy rights)

### Security Certifications
Kyros aims to achieve or support:
- **Common Criteria**: ISO/IEC 15408 certification for security evaluation
- **FIPS 140-2/3**: Federal Information Processing Standards for cryptographic modules
- **SOC 1**: Financial reporting controls
- **ISO 22301**: Business continuity management systems
- **ISO 20000-1**: IT service management
- **ISO 9001**: Quality management systems

### Audit Features
Kyros provides comprehensive audit capabilities:
- **Immutable Audit Log**: Cryptographically secured audit trail
- **Log Integrity Verification**: Ability to verify log integrity
- **Log Search and Retrieval**: Efficient search and retrieval of audit events
- **Report Generation**: Pre-built and customizable audit reports
- **Retention Policies**: Configurable retention with archival
- **Access Controls**: Role-based access to audit logs
- **Real-Time Alerting**: Alerting on suspicious audit events
- **Export Capabilities**: Multiple formats (JSON, CSV, XML, PDF)

## Performance Benchmarks

### Baseline Performance
Single instance performance on moderate hardware (8 cores, 32GB RAM, SSD):
- **Blob Upload**: 150-300 MB/s (network and storage dependent)
- **Blob Download**: 300-600 MB/s (network and storage dependent)
- **Manifest Operations**: 500-1000 operations/second
- **API Requests**: 1000-2000 requests/second (simple operations)
- **Authenticated Requests**: 800-1500 requests/second (with JWT validation)
- **Concurrent Connections**: 2000-5000 simultaneous connections

### Scaling Characteristics
- **Linear Scaling**: Near-linear scaling with additional instances for stateless operations
- **Database Bound**: Metadata operations eventually limited by database capacity
- **Storage Bound**: Blob operations limited by storage backend throughput
- **Network Bound**: High-volume scenarios limited by network bandwidth
- **Cache Effective**: High cache hit ratios (>95%) for repeated access patterns

### Optimization Targets
- **Target Latency**: <50ms for 95% of API requests
- **Target Throughput**: 10,000+ RPM for mixed workload
- **Target Availability**: 99.9% uptime SLA
- **Target Recovery Time**: <15 minutes for partial failures, <4 hours for complete site failure
- **Target Data Durability**: 99.999999999% (11 nines) annual durability
- **Target RPO**: <15 minutes for user data
- **Target RTO**: <4 hours for full service restoration

## Future Enhancements

### Planned Features
1. **Manifest Lists (Multi-Arch Images)**: Enhanced support for manifest lists and platform selection
2. **Referrers API Enhancements**: Advanced referrers querying and filtering capabilities
3. **Garbage Collection Policies**: Policy-driven garbage collection decisions
4. **Storage Plugins**: Additional storage backends (Ceph, HDFS, on-premises NAS)
5. **Image Signing Integration**: Native Cosign integration for push-time signing
6. **SBOM Storage**: Native SBOM storage and retrieval capabilities
7. **Scan Integration**: Built-in vulnerability scanning during push
8. **Policy as Code**: Enhanced OPA integration with built-in policy library
9. **Replication Enhancements**: Active-active geo-replication with conflict resolution
10. **Mirror Improvements**: Smart mirroring with intelligent source selection

### Performance Improvements
1. **Read-Ahead Prefetching**: Predictive blob prefetching for sequential access
2. **Adaptive Chunking**: Dynamic chunk size based on network conditions
3. **Connection Multiplexing**: HTTP/2 connection multiplexing for reduced overhead
4. **Compression**: Optional blob compression for storage and bandwidth efficiency
5. **Edge Caching**: CDN integration for geographically distributed caching
6. **Database Sharding**: Horizontal scaling of metadata storage
7. **Blob Storage Tiering**: Automatic movement to optimal storage tiers
8. **Query Optimization**: Advanced query planning and execution optimization
9. **Async Processing**: Increased use of asynchronous processing for non-critical paths
10. **Resource Pooling**: Enhanced connection and resource pooling mechanisms

### Security Improvements
1. **Hardware Security Modules**: HSM integration for key protection
2. **Confidential Computing**: TEEs for sensitive operations
3. **Zero Trust Networking**: Microsegmentation and encrypted service mesh
4. **Advanced Threat Detection**: Behavioral analysis and anomaly detection
5. **Automated Response**: Automated containment and remediation of threats
6. **Supply Chain Security**: Enhanced SLSA framework integration
7. **Code Signing**: Registry-level code signing for provenance
8. **Provenance Tracking**: Detailed provenance tracking for all artifacts
9. **Immutable Logs**: Write-once-read-many (WORM) storage for audit logs
10. **Quantum-Resistant Cryptography**: Preparation for post-quantum cryptography

### User Experience Improvements
1. **Enhanced Web UI**: Improved repository browsing, search, and discovery
2. **CLI Improvements**: Enhanced CLI with autocomplete and better error messages
3. **Mobile Support**: Optimized mobile experience for registry management
4. **Self-Service Portal**: Tenant self-service for common operations
5. **API Discoverability**: Improved API documentation and discovery
6. **SDK Improvements**: Enhanced client libraries with better ergonomics
7. **Template Library**: Pre-built configurations for common use cases
8. **Troubleshooting Tools**: Built-in diagnostics and troubleshooting utilities
9. **Training Materials**: Comprehensive documentation and learning resources
10. **Community Features**: Forums, marketplace, and community contributions

## Diagrams Reference
See [MERMAID.md](MERMAID.md) for detailed Mermaid diagrams including:
- Registry Architecture Overview
- Storage Layer Implementation
- Authentication and Authorization Flow
- API Request Processing Pipeline
- Garbage Collection Algorithm
- High Availability Deployment Patterns
- Disaster Recovery Architecture
- Performance Optimization Techniques
- Security Implementation Details
- Observability Integration Points