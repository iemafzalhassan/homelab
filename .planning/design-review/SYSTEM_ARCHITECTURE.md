# Kyros System Architecture

## High-Level Architecture

Kyros follows a cloud-native, microservices architecture designed for scalability, resilience, and extensibility. The system is composed of loosely coupled services that communicate through well-defined APIs and event-driven mechanisms.

### Core Architectural Principles
1. **Microservices**: Independently deployable, scalable services
2. **API-First**: All service interactions through well-defined contracts
3. **Event-Driven**: Asynchronous communication via NATS JetStream
4. **Observability**: Built-in metrics, logging, and tracing
5. **Security**: Zero-trust model with mutual TLS and JWT authentication
6. **Extensibility**: Plugin architecture for custom functionality
7. **GitOps Native**: Designed for declarative infrastructure management

## Component Boundaries

### Core Services
1. **Registry Service** (`kyros/cmd/registry/`)
   - OCI-compliant container image storage and distribution
   - Based on cncf/distribution v3 registry engine
   - Handles image push/pull, blob storage, manifest management

2. **API Service** (`kyros/cmd/api/`)
   - RESTful API for platform management and operations
   - User authentication, authorization, and session management
   - Repository and namespace management
   - Webhook configuration and delivery

3. **Trust Score Service** (`kyros/cmd/trustscore/`)
   - Calculates trust scores for images based on multiple factors
   - Integrates with SBOM, vulnerability scanners, and signature verification
   - Provides policy evaluation engine

4. **Webhook Service** (`kyros/cmd/webhook/`)
   - Manages webhook subscriptions and delivery
   - Handles event filtering and transformation
   - Implements retry mechanisms and dead-letter queues

5. **Authentication Service** (`kyros/cmd/auth/`)
   - Integrates with Keycloak for OIDC authentication
   - Manages tokens, sessions, and user information
   - Provides authentication middleware for other services

6. **Operator Service** (`kyros/cmd/operator/`)
   - Kubernetes Operator for Kyros platform management
   - Handles installation, upgrades, and configuration
   - Manages custom resources (TrustPolicies, Webhooks, etc.)

### Supporting Services
1. **Keycloak** - Identity and access management (external dependency)
2. **PostgreSQL** - Primary database for metadata and user information
3. **Elasticsearch** - Search and analytics for logs and audit trails
4. **Redis** - Caching and session storage
5. **NATS JetStream** - Event streaming and messaging backbone
6. **MinIO/S3** - Blob storage for image layers and artifacts
7. **Prometheus** - Metrics collection and storage
8. **Grafana** - Visualization and dashboarding
9. **Tempo** - Distributed tracing backend
10. **Loki** - Log aggregation system

## Microservices Details

### Registry Service
- **Responsibilities**: OCI registry functionality, image storage, distribution
- **Interfaces**: 
  - Docker Registry HTTP API v2 (REST)
  - Internal gRPC for service-to-service communication
  - NATS JetStream for event publishing
- **Dependencies**: 
  - PostgreSQL (metadata)
  - MinIO/S3 (blob storage)
  - Redis (caching)
  - NATS JetStream (events)
- **Scalability**: Horizontally scalable behind load balancer
- **Failure Handling**: Stateless design allows for quick recovery

### API Service
- **Responsibilities**: Platform management API, UI backend, business logic
- **Interfaces**:
  - REST API (external and internal)
  - GraphQL API (planned)
  - Internal gRPC
  - NATS JetStream (events)
- **Dependencies**:
  - PostgreSQL (primary data store)
  - Redis (caching and sessions)
  - Keycloak (authentication)
  - NATS JetStream (events)
- **Scalability**: Horizontally scalable with session affinity for WebSocket connections
- **Failure Handling**: Stateless except for WebSocket connections

### Trust Score Service
- **Responsibilities**: Trust score calculation, policy evaluation, security analysis
- **Interfaces**:
  - gRPC API (internal)
  - NATS JetStream (event consumption and publishing)
  - REST API (administrative)
- **Dependencies**:
  - PostgreSQL (trust score storage)
  - Redis (caching of scan results)
  - External scanners (Trivy, Grype, Syft)
  - NATS JetStream (events)
- **Scalability**: Horizontally scalable; CPU-intensive tasks can be distributed
- **Failure Handling**: Idempotent processing allows for safe retries

### Webhook Service
- **Responsibilities**: Webhook management, event delivery, retry logic
- **Interfaces**:
  - REST API (configuration)
  - NATS JetStream (event consumption)
  - HTTP/Webhook (event delivery)
- **Dependencies**:
  - PostgreSQL (webhook configuration storage)
  - Redis (delivery attempt tracking)
  - NATS JetStream (event source)
- **Scalability**: Horizontally scalable; delivery workers can be scaled independently
- **Failure Handling**: Built-in retry with exponential backoff and dead-letter queue

### Authentication Service
- **Responsibilities**: Authentication middleware, token validation, user info
- **Interfaces**:
  - Internal gRPC (service-to-service auth)
  - HTTP (token introspection endpoints)
  - NATS JetStream (auth events)
- **Dependencies**:
  - Keycloak (primary identity provider)
  - Redis (token caching)
  - NATS JetStream (events)
- **Scalability**: Horizontally scalable; primarily I/O bound
- **Failure Handling**: Stateless; failures degrade to unauthenticated requests

### Operator Service
- **Responsibilities**: Kubernetes-native platform management
- **Interfaces**:
  - Kubernetes API (custom resources)
  - Internal gRPC (service communication)
  - NATS JetStream (operational events)
- **Dependencies**:
  - Kubernetes cluster
  - All Kyros services (for health checks and management)
  - Helm (for chart management)
- **Scalability**: Single instance (leader-elected) for coordination
- **Failure Handling**: Leader election ensures high availability

## Dependencies

### External Dependencies
1. **Keycloak** - Identity and Access Management
2. **PostgreSQL** - Relational Database
3. **MinIO/S3** - Object Storage (blob storage)
4. **Redis** - In-Memory Data Store
5. **NATS JetStream** - Event Streaming Platform
6. **Trivy/Grype** - Vulnerability Scanners
7. **Syft** - SBOM Generator
8. **Cosign** - Image Signing Utility
9. **Prometheus** - Metrics Collection
10. **Grafana** - Visualization
11. **Tempo** - Distributed Tracing
12. **Loki** - Log Aggregation

### Internal Dependencies
- Services communicate via gRPC for low-latency, type-safe communication
- Event-driven communication via NATS JetStream for loose coupling
- Shared libraries for common functionality (logging, configuration, tracing)

## Data Flow

### Image Push Flow
1. Client authenticates with API Service (via Keycloak)
2. Client initiates image push to Registry Service
3. Registry Service validates authentication and permissions
4. Registry Service receives and stores image layers (blobs) in MinIO/S3
5. Registry Service stores manifest and metadata in PostgreSQL
6. Registry Service publishes `image.pushed` event to NATS JetStream
7. Trust Score Service consumes event and initiates analysis
8. Webhook Service consumes event and delivers configured webhooks
9. API Service updates repository metadata and search indexes

### Image Pull Flow
1. Client authenticates with API Service (via Keycloak)
2. Client initiates image pull from Registry Service
3. Registry Service validates authentication and permissions
4. Registry Service retrieves image layers from MinIO/S3
5. Registry Service assembles manifest and delivers to client
6. Registry Service publishes `image.pulled` event to NATS JetStream
7. API Service updates pull statistics and last accessed timestamps

### Trust Score Calculation Flow
1. Trust Score Service consumes `image.pushed` event from NATS JetStream
2. Service retrieves image manifest and layers from Registry Service
3. Service runs SBOM generation (Syft) on image layers
4. Service runs vulnerability scanning (Trivy/Grype) on SBOM
5. Service verifies image signatures (Cosign) if present
6. Service calculates trust score based on configurable policy
7. Service stores trust score and analysis results in PostgreSQL
8. Service publishes `trust.score.updated` event to NATS JetStream
9. Webhook Service consumes event and delivers trust score webhooks
10. API Service makes trust score available via API

### Webhook Delivery Flow
1. Event published to NATS JetStream (image.pushed, trust.score.updated, etc.)
2. Webhook Service consumes event and matches against subscriptions
3. For each matching subscription:
   - Service prepares payload according to subscription configuration
   - Service attempts HTTP delivery to webhook URL
   - On failure, service implements exponential backoff retry
   - After max retries, event sent to dead-letter queue
4. Delivery attempts and results stored in PostgreSQL for auditing
5. Service publishes `webhook.delivered` or `webhook.failed` events

## Failure Scenarios

### Service Unavailability
- **Registry Service Down**: Image push/pull operations fail; API returns 503
- **API Service Down**: Platform management unavailable; registry operations continue via direct registry access
- **Trust Score Service Down**: Trust score calculation delayed; images marked as "unscored"
- **Webhook Service Down**: Webhook delivery delayed; events queued in NATS until service recovers
- **Authentication Service Down**: All services fall back to unauthenticated mode or cached tokens
- **Operator Service Down**: Platform management operations paused; existing services continue running

### Data Store Failures
- **PostgreSQL Unavailable**: 
  - Registry Service: Can't persist new images or metadata (failures after timeout)
  - API Service: Can't manage users/repositories; read-only operations possible
  - Trust Score Service: Can't store analysis results
  - Mitigation: Read replicas, connection pooling, circuit breakers
  
- **Redis Unavailable**:
  - API Service: Session storage falls back to JWT tokens; performance degradation
  - Trust Score Service: Loss of caching for scan results; increased scanner load
  - Webhook Service: Loss of delivery attempt tracking; potential duplicate deliveries
  - Mitigation: Redis clustering, fallback to database storage

- **NATS JetStream Unavailable**:
  - All services: Event-driven communication breaks down
  - Mitigation: NATS clustering; services implement local queuing for critical events
  
- **MinIO/S3 Unavailable**:
  - Registry Service: Cannot store or retrieve image layers; operations fail
  - Mitigation: Multi-region replication, gateway pattern

### Network Partitions
- **Inter-service Communication**: 
  - Services implement circuit breaker pattern (e.g., using resilience4j)
  - Fallback to cached data where possible
  - Degraded functionality rather than complete failure
  
- **External Dependencies**:
  - Keycloak: Services cache tokens and validate signatures locally
  - Scanners: Trust Score Service queues scans and processes when connectivity restored

### Resource Exhaustion
- **CPU/Memory**: 
  - Kubernetes HPA scales services based on resource utilization
  - Trust Score Service: Scan workers can be scaled independently
  - Webhook Service: Delivery workers can be scaled based on volume
  
- **Storage**:
  - Blob Storage: Quotas enforced per namespace/repository
  - Database: Connection pooling and query optimization
  - Logs: Retention policies and rotation

## Scalability Considerations

### Horizontal Scaling
- **Stateless Services**: API, Auth, Webhook services scale horizontally behind load balancer
- **Stateful Services**: 
  - Registry Service: Scalable with shared blob storage (MinIO/S3) and database
  - Trust Score Service: Scalable with work distribution via NATS work queues
  - Operator Service: Single active instance with leader election

### Vertical Scaling
- **Resource Limits**: 
  - CPU/Memory requests and limits defined per service
  - Trust Score Service allocated more CPU for scanning workloads
  - Registry Service allocated more memory for blob cache
  
### Database Scaling
- **PostgreSQL**:
  - Read replicas for query-heavy operations (API service reads)
  - Connection pooling (PgBouncer) to manage connection counts
  - Partitioning for large tables (events, audit logs)
  
### Caching Strategy
- **Redis Layers**:
  - L1: Local in-memory caches (where appropriate)
  - L2: Redis shared cache for distributed services
  - L3: HTTP/CDN caching for static assets and frequently accessed blobs
  
### CDN and Edge Caching
- **Static Assets**: Served via CDN for global distribution
- **Image Layers**: Optional CDN caching for public repositories (with security considerations)
- **API Responses**: Cache-control headers for cacheable responses

### Performance Optimization
- **Blob Storage**: 
  - Chunked upload/download for large images
  - Parallel layer transfer
  - Blob deduplication
  
- **Database**:
  - Index optimization for frequent queries
  - Read-only replicas for reporting and analytics
  - Connection pooling
  
- **Network**:
  - gRPC for efficient service-to-service communication
  - HTTP/2 for client-service communication
  - Compression for large payloads

## Technology Stack

### Languages
- **Go**: Primary language for backend services (performance, concurrency, cloud-native)
- **TypeScript**: Frontend web UI (React/Next.js ecosystem)
- **Python**: Some utility scripts and ML components (Trust Score AI enhancements)

### Frameworks & Libraries
- **Go**:
  - Gin/Gorilla Mux for HTTP routing
  - gRPC for service-to-service communication
  - Viper for configuration management
  - Zap for structured logging
  - Prometheus client for metrics
  - OpenTelemetry for tracing
  
- **TypeScript/JS**:
  - Next.js 15 for React framework
  - Tailwind CSS for styling
  - React Query for data fetching
  - WebSocket for real-time updates
  
### Infrastructure
- **Containerization**: Docker images for all services
- **Orchestration**: Kubernetes (primary target)
- **Service Mesh**: Optional (Istio/Linkerd for advanced traffic management)
- **CI/CD**: GitHub Actions with automated testing and deployment
- **GitOps**: ArgoCD for application deployment, Kargo for promotion

### Protocols & Standards
- **OCI Distribution Spec**: Container registry protocol
- **OCI Image Spec**: Image format standard
- **REST/JSON**: Primary API format
- **gRPC/Protobuf**: Internal service communication
- **NATS**: Event streaming protocol
- **OpenAPI 3.0**: API specification standard
- **OIDC/OAuth 2.0**: Authentication and authorization
- **JWT**: Token format
- **Webhook**: HTTP-based event delivery

## Security Architecture

### Zero Trust Principles
1. **Verify Explicitly**: All requests authenticated and authorized
2. **Least Privilege**: Minimal permissions granted
3. **Assume Breach**: Encryption everywhere, microsegmentation
4. **Secure by Default**: Secure configurations as default

### Authentication Flow
1. User authenticates with Keycloak (username/password, SSO, MFA)
2. Keycloak issues JWT access token and refresh token
3. Client presents JWT to Kyros services in Authorization header
4. Services validate JWT signature and claims
5. Services check permissions against policy engine
6. Short-lived tokens (15min) with refresh capability

### Authorization Model
- **RBAC**: Role-Based Access Control with predefined roles
- **ABAC**: Attribute-Based Access Control for fine-grained policies
- **Resource-Based**: Permissions tied to specific resources (repositories, namespaces)
- **Policy Engine**: OPA integration for complex policy decisions

### Data Protection
- **Encryption at Rest**:
  - Database: Transparent encryption or application-level encryption
  - Blob Storage: Server-side encryption (SSE-S3, SSE-KMS)
  - Backups: Encrypted storage
  
- **Encryption in Transit**:
  - mTLS between all services
  - HTTPS for all external communications
  - Encrypted NATS connections
  
- **Secrets Management**:
  - External secrets manager (HashiCorp Vault, Azure Key Vault)
  - Kubernetes secrets for cluster-only secrets
  - Automatic rotation of service credentials

### Network Security
- **Network Policies**: Kubernetes NetworkPolicies restricting service communication
- **Service Mesh**: Optional mTLS and traffic policies
- **Ingress Control**: API Gateway (Traefik) with rate limiting and WAF
- **Private Networks**: Services deployed in private subnets; no direct internet access

### Audit and Compliance
- **Audit Logging**: All privileged operations logged to immutable storage
- **Log Retention**: Configurable retention periods with archival
- **Compliance Reporting**: Built-in reports for SOC 2, ISO 27001, etc.
- **Vulnerability Management**: Regular scanning of Kyros container images

## Observability Architecture

### Metrics
- **Service-Level Metrics**: 
  - Request rates, error rates, latency (RED metrics)
  - Resource utilization (CPU, memory, disk, network)
  - Business metrics (image pushes/pulls, trust scores, webhook deliveries)
  
- **Collection**: Prometheus client libraries in each service
- **Storage**: Prometheus TSDB with remote write to long-term storage
- **Visualization**: Grafana dashboards with pre-built panels
- **Alerting**: Prometheus Alertmanager with routing to Slack, email, PagerDuty

### Logging
- **Structured Logging**: JSON-formatted logs with trace IDs
- **Collection**: Agent daemons (Fluent Bit/Filebeat) on each node
- **Storage**: Loki log aggregation system
- **Retention**: Index-based retention with automated cleanup
- **Search**: Grafana Explore or Loki CLI for log investigation

### Tracing
- **Instrumentation**: OpenTelemetry SDK in all services
- **Context Propagation**: W3C TraceContext headers
- **Collection**: OpenTelemetry Collector agents
- **Storage**: Tempo trace database
- **Visualization**: Grafana TraceQL integration
- **Sampling**: Adaptive sampling based on traffic volume

### Profiling
- **Continuous Profiling**: eBPF-based profiling for CPU/memory analysis
- **Heap Profiling**: Go pprof endpoints exposed via debug endpoints
- **Block Profiling**: Goroutine blocking analysis
- **Mutex Profiling**: Contention analysis

### Health Checks
- **Liveness Probes**: Kubernetes liveness probes for container restart
- **Readiness Probes**: Kubernetes readiness probes for service endpoint removal
- **Startup Probes**: Kubernetes startup probes for slow-starting applications
- **Health Endpoints**: `/healthz` endpoint returning detailed service status

### DORA Metrics
- **Deployment Frequency**: Tracked via GitOps deployment events
- **Lead Time for Changes**: Measured from commit to production
- **Change Failure Rate**: Percentage of deployments causing incidents
- **Time to Restore Service**: MTTR from incident detection to resolution
- **Collection**: Automated from GitOps and incident management systems

## Deployment Architecture

### Kubernetes-Native Deployment
- **Namespace Isolation**: 
  - `kyros-platform`: Core platform services
  - `kyros-monitoring`: Observability stack (Prometheus, Grafana, etc.)
  - `kyros-dependencies`: External dependencies (PostgreSQL, Redis, etc.)
  - `kyros-workloads`: User workloads (if running user containers in same cluster)
  
- **Helm Charts**:
  - Umbrella chart for complete platform deployment
  - Subcharts for each service and dependency
  - Configuration values for environment-specific settings
  
- **Operator Pattern**:
  - Kyros Operator manages application lifecycle
  - Custom Resources: KyrosCluster, TrustPolicy, Webhook, Repository
  - Reconciliation loops for drift correction
  
### Deployment Strategies
- **Blue/Green**: For major version upgrades
- **Canary**: For gradual rollout of new features
- **Rolling Update**: Default for patch releases
- **Recreate**: For stateful services requiring downtime

### Configuration Management
- **External Configuration**: 
  - ConfigMaps for non-sensitive configuration
  - Secrets for sensitive data (passwords, keys)
  - External secret stores (Vault, Azure Key Vault) for enterprise secrets
  
- **Dynamic Configuration**:
  - Feature flags via Unleash or similar
  - Runtime configuration updates without restart
  - Configuration validation on startup

### Backup and Disaster Recovery
- **Backup Strategy**:
  - Database: Regular logical backups (pg_dump) + snapshots
  - Blob Storage: Versioning + cross-region replication
  - Configuration: GitOps repository as source of truth
  
- **Recovery Procedures**:
  - Step-by-step recovery guides for different failure scenarios
  - Automated failover for multi-cluster deployments
  - RTO/RPO targets defined per service type
  
### Multi-Cluster and Geographic Distribution
- **Active-Active**: 
  - Multi-cluster deployment with traffic routing
  - Conflict-free replicated data types (CRDTs) for eventual consistency
  - Registry mirroring for image availability
  
- **Active-Passive**:
  - Warm standby cluster with automated failover
  - Regular synchronization of critical data
  
- **Geo-Distribution**:
  - Edge caching for frequently accessed images
  - Regional replicas for compliance requirements
  - Latency-based routing via Global Server Load Balancing (GSLB)

## Evolution and Extensibility

### Plugin Architecture
- **Extension Points**:
  - Authentication providers (LDAP, SAML, custom OIDC)
  - Storage backends (Azure Blob, GCS, on-premises NAS)
  - Scanning engines (different vulnerability scanners)
  - Notification channels (Slack, Teams, custom webhooks)
  - Policy engines (OPA, custom policy languages)
  
- **Plugin Framework**:
  - Go plugins or gRPC-based external services
  - Versioned plugin APIs
  - Plugin marketplace and registry
  
### API Extensibility
- **Versioning**: 
  - URI versioning (/api/v1/, /api/v2/)
  - Header-based versioning as alternative
  - Deprecation policy with sunset dates
  
- **Expansion**:
  - Additive changes only (new fields, optional fields)
  - Clear deprecation warnings
  - Semantic versioning for API contracts
  
### UI Extensibility
- **Plugin Slots**: 
  - Dashboard widgets
  - Navigation extensions
  - Context menu items
  - Custom views and pages
  
- **Themeability**:
  - CSS variables for easy theming
  - Dark/light mode support
  - Custom branding and white-labeling
  
### Data Model Evolution
- **Backward Compatibility**:
  - Database migrations that are backward-compatible
  - Schema versioning
  - Rolling upgrade support
  
- **Forward Compatibility**:
  - Ignore unknown fields in API requests/responses
  - Default values for new fields
  - Feature flags for gradual rollout

## Diagrams Reference
See [MERMAID.md](MERMAID.md) for detailed Mermaid diagrams including:
- System Architecture
- Service Communication Patterns
- Data Flow Diagrams
- Deployment Architecture
- Failure Scenario Handling
- Scaling Patterns