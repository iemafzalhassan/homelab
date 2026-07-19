# Kyros Database Design

## Entity Relationship Diagram Overview

Kyros uses a relational database (PostgreSQL) as the primary store for metadata, user information, configuration, and operational data. The database design follows normalization principles while considering performance requirements for common access patterns.

## Core Tables

### 1. Users and Authentication
```sql
-- Users table stores user profiles
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    enabled BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Groups for organizing users
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Many-to-many relationship between users and groups
CREATE TABLE user_groups (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, group_id)
);

-- Roles define permissions that can be assigned
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    scope VARCHAR(50) DEFAULT 'realm', -- realm, client
    client_id UUID NULL REFERENCES clients(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(name, scope, client_id)
);

-- Many-to-many relationship between users and roles
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    PRIMARY KEY (user_id, role_id)
);

-- Many-to-many relationship between groups and roles
CREATE TABLE group_roles (
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    PRIMARY KEY (group_id, role_id)
);

-- Clients represent applications/services that authenticate
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    secret_hash TEXT NOT NULL, -- bcrypt hash of client secret
    redirect_uris TEXT ARRAY DEFAULT '{}',
    enabled BOOLEAN DEFAULT TRUE,
    protocol VARCHAR(50) DEFAULT 'openid-connect',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Access tokens for API authentication
CREATE TABLE access_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    client_id UUID REFERENCES clients(id) ON DELETE SET NULL,
    token_hash TEXT NOT NULL, -- bcrypt hash of token
    scopes TEXT ARRAY DEFAULT '{}',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(token_hash)
);

-- Refresh tokens for obtaining new access tokens
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    client_id UUID REFERENCES clients(id) ON DELETE SET NULL,
    token_hash TEXT NOT NULL, -- bcrypt hash of token
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(token_hash)
);
```

### 2. Registry Core
```sql
-- Namespaces provide logical grouping of repositories
CREATE TABLE namespaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    visibility VARCHAR(50) DEFAULT 'private', -- public, private, protected
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- Repositories store collections of artifacts
CREATE TABLE repositories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    namespace_id UUID NOT NULL REFERENCES namespaces(id) ON DELETE CASCADE,
    description TEXT,
    visibility VARCHAR(50) DEFAULT 'inherit', -- public, private, protected, inherit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(name, namespace_id)
);

-- Artifacts represent stored OCI objects (images, charts, etc.)
CREATE TABLE artifacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    digest VARCHAR(255) UNIQUE NOT NULL, -- SHA256 digest
    media_type VARCHAR(255) NOT NULL,
    size_bytes BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(repository_id, digest)
);

-- Tags provide mutable pointers to artifacts
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(artifact_id, name)
);

-- Blobs store content-addressed layers and configs
CREATE TABLE blobs (
    digest VARCHAR(255) PRIMARY KEY, -- SHA256 digest
    media_type VARCHAR(255) NOT NULL,
    size_bytes BIGINT NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- Many-to-many relationship between artifacts and blobs (manifest -> layers/config)
CREATE TABLE artifact_blobs (
    artifact_id UUID REFERENCES artifacts(id) ON DELETE CASCADE,
    blob_digest VARCHAR(255) REFERENCES blobs(digest) ON DELETE CASCADE,
    path VARCHAR(255), -- JSON path in manifest referencing this blob
    PRIMARY KEY (artifact_id, blob_digest)
);
```

### 3. Trust and Security
```sql
-- Vulnerabilities discovered in artifacts
CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    scanner VARCHAR(100) NOT NULL, -- Trivy, Grype, etc.
    vulnerability_id VARCHAR(255) NOT NULL, -- CVE or scanner-specific ID
    severity VARCHAR(50) NOT NULL, -- unknown, low, medium, high, critical
    title VARCHAR(255) NOT NULL,
    description TEXT,
    references TEXT ARRAY DEFAULT '{}',
    fixed_version VARCHAR(255),
    discovered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for vulnerability queries
CREATE INDEX idx_vulnerabilities_artifact_id ON vulnerabilities(artifact_id);
CREATE INDEX idx_vulnerabilities_severity ON vulnerabilities(severity);
CREATE INDEX idx_vulnerabilities_scanner ON vulnerabilities(scanner);

-- SBOMs (Software Bills of Materials)
CREATE TABLE sboms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL UNIQUE REFERENCES artifacts(id) ON DELETE CASCADE,
    format VARCHAR(50) NOT NULL, -- SPDX, CycloneDX
    content JSONB NOT NULL,
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    generator VARCHAR(255) NOT NULL -- Syft, etc.
);

-- Image signatures (Cosign, Notary, PGP)
CREATE TABLE signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- cosign, notary, pgp
    key_id VARCHAR(255) NOT NULL,
    payload TEXT NOT NULL,
    signature_text TEXT NOT NULL,
    verified_at TIMESTAMP WITH TIME ZONE,
    verification_status VARCHAR(50) DEFAULT 'pending', -- pending, verified, failed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Trust scores calculated for artifacts
CREATE TABLE trust_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL UNIQUE REFERENCES artifacts(id) ON DELETE CASCADE,
    score FLOAT NOT NULL CHECK (score >= 0.0 AND score <= 1.0),
    level VARCHAR(50) NOT NULL, -- unknown, untrusted, low, medium, high, trusted
    factors JSONB, -- Breakdown of scoring factors
    policy_id UUID REFERENCES policies(id) ON DELETE SET NULL,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Policies that define acceptable artifact characteristics
CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    rules JSONB NOT NULL, -- Policy rules in Rego or similar
    scope VARCHAR(50) DEFAULT 'global', -- global, namespace, repository
    namespace_id UUID NULL REFERENCES namespaces(id) ON DELETE SET NULL,
    repository_id UUID NULL REFERENCES repositories(id) ON DELETE SET NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (
        (scope = 'global' AND namespace_id IS NULL AND repository_id IS NULL) OR
        (scope = 'namespace' AND namespace_id IS NOT NULL AND repository_id IS NULL) OR
        (scope = 'repository' AND namespace_id IS NULL AND repository_id IS NOT NULL)
    )
);

-- Policy evaluation results (audit trail)
CREATE TABLE policy_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    result VARCHAR(50) NOT NULL, -- pass, fail, warn, error
    details JSONB, -- Detailed evaluation results
    evaluated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    evaluated_by UUID REFERENCES users(id) ON DELETE SET NULL
);
```

### 4. Observability
```sql
-- Metrics collected from services
CREATE TABLE metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(50),
    tags JSONB, -- Key-value pairs for metric dimensions
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for metric queries
CREATE INDEX idx_metrics_service_name ON metrics(service_name);
CREATE INDEX idx_metrics_metric_name ON metrics(metric_name);
CREATE INDEX idx_metrics_timestamp ON metrics(timestamp);
CREATE INDEX idx_metrics_service_metric_time ON metrics(service_name, metric_name, timestamp);

-- Log entries from services
CREATE TABLE log_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    level VARCHAR(50) NOT NULL, -- debug, info, warn, error, fatal
    message TEXT NOT NULL,
    trace_id VARCHAR(255),
    span_id VARCHAR(255),
    attributes JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for log queries
CREATE INDEX idx_log_entries_service_name ON log_entries(service_name);
CREATE INDEX idx_log_entries_level ON log_entries(level);
CREATE INDEX idx_log_entries_timestamp ON log_entries(timestamp);
CREATE INDEX idx_log_entries_trace_id ON log_entries(trace_id);

-- Distributed traces
CREATE TABLE traces (
    id VARCHAR(255) PRIMARY KEY, -- TraceID
    service_name VARCHAR(255) NOT NULL,
    operation_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    duration_ns BIGINT,
    status VARCHAR(50) DEFAULT 'ok', -- ok, error
    tags JSONB
);

-- Indexes for trace queries
CREATE INDEX idx_traces_service_name ON traces(service_name);
CREATE INDEX idx_traces_start_time ON traces(start_time);
CREATE INDEX idx_traces_status ON traces(status);

-- Span entries within traces
CREATE TABLE spans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trace_id VARCHAR(255) NOT NULL REFERENCES traces(id) ON DELETE CASCADE,
    service_name VARCHAR(255) NOT NULL,
    operation_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    duration_ns BIGINT,
    tags JSONB,
    parent_span_id UUID NULL REFERENCES spans(id) ON DELETE SET NULL
);

-- Indexes for span queries
CREATE INDEX idx_spans_trace_id ON spans(trace_id);
CREATE INDEX idx_spans_service_name ON spans(service_name);
CREATE INDEX idx_spans_start_time ON spans(start_time);

-- Alert rules
CREATE TABLE alert_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    condition TEXT NOT NULL, -- PromQL or similar expression
    severity VARCHAR(50) NOT NULL, -- info, warning, error, critical
    notification_channels JSONB, -- Configured notification targets
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Active alerts
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id UUID NOT NULL REFERENCES alert_rules(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL, -- pending, firing, resolved
    severity VARCHAR(50) NOT NULL,
    summary TEXT NOT NULL,
    description TEXT,
    value TEXT, -- Current value that triggered alert
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for alert queries
CREATE INDEX idx_alerts_rule_id ON alerts(rule_id);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_started_at ON alerts(started_at);
```

### 5. GitOps and Automation
```sql
-- Webhook subscriptions
CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    url TEXT NOT NULL,
    events TEXT ARRAY NOT NULL, -- Subscribed event types
    secret_hash TEXT, -- bcrypt hash for HMAC verification
    format VARCHAR(50) DEFAULT 'JSON', -- JSON, form-urlencoded
    headers JSONB DEFAULT '{}', -- HTTP headers to include
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_triggered_at TIMESTAMP WITH TIME ZONE,
    failure_count INTEGER DEFAULT 0,
    next_retry_at TIMESTAMP WITH TIME ZONE
);

-- Webhook delivery attempts
CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event_id VARCHAR(255), -- NATS JetStream message ID or similar
    attempt INTEGER NOT NULL DEFAULT 1,
    status VARCHAR(50) NOT NULL, -- pending, success, failed
    response_code INTEGER,
    response_body TEXT, -- Truncated response body
    attempted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for webhook delivery queries
CREATE INDEX idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_attempted_at ON webhook_deliveries(attempted_at);
```

### 6. Multi-tenancy
```sql
-- Tenants (organizations, teams, etc.)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- Many-to-many relationship between tenants and users
CREATE TABLE tenant_users (
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member', -- owner, admin, member, viewer
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (tenant_id, user_id)
);

-- Many-to-many relationship between tenants and groups
CREATE TABLE tenant_groups (
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (tenant_id, group_id)
);

-- Namespace quotas for resource limits
CREATE TABLE namespace_quotas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    namespace_name VARCHAR(255) NOT NULL, -- Kubernetes namespace name
    hard_limits JSONB NOT NULL, -- Resource quotas (CPU, memory, storage, etc.)
    used_resources JSONB DEFAULT '{}', -- Current usage
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for namespace quota queries
CREATE INDEX idx_namespace_quotas_tenant_id ON namespace_quotas(tenant_id);
CREATE INDEX idx_namespace_quotas_namespace_name ON namespace_quotas(namespace_name);
```

### 7. Audit and Compliance
```sql
-- Audit trail for all privileged operations
CREATE TABLE audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    actor_id UUID NULL REFERENCES users(id) ON DELETE SET NULL, -- NULL for system actions
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(255) NOT NULL,
    resource_id UUID NULL,
    resource_name VARCHAR(255),
    outcome VARCHAR(50) NOT NULL, -- success, failure
    failure_reason TEXT,
    ip_address INET,
    user_agent TEXT,
    request_id VARCHAR(255),
    changes JSONB, -- Before/after values for updates
    tenant_id UUID NULL REFERENCES tenants(id) ON DELETE SET NULL
);

-- Indexes for audit queries
CREATE INDEX idx_audit_events_timestamp ON audit_events(timestamp);
CREATE INDEX idx_audit_events_actor_id ON audit_events(actor_id);
CREATE INDEX idx_audit_events_resource_type ON audit_events(resource_type);
CREATE INDEX idx_audit_events_outcome ON audit_events(outcome);
CREATE INDEX idx_audit_events_tenant_id ON audit_events(tenant_id);
```

## Indexes for Performance

### Composite Indexes for Common Queries
```sql
-- Repository access by namespace
CREATE INDEX idx_repositories_namespace_id ON repositories(namespace_id);

-- Artifact access by repository
CREATE INDEX idx_artifacts_repository_id ON artifacts(repository_id);

-- Tag access by artifact
CREATE INDEX idx_tags_artifact_id ON tags(artifact_id);

-- Blob lookup by size/media type (for GC)
CREATE INDEX idx_blobs_size ON blobs(size_bytes);
CREATE INDEX idx_blobs_media_type ON blobs(media_type);

-- Vulnerability lookup by artifact and severity
CREATE INDEX idx_vulnerabilities_artifact_severity ON vulnerabilities(artifact_id, severity);

-- Trust score lookup
CREATE INDEX idx_trust_scores_artifact_id ON trust_scores(artifact_id);
CREATE INDEX idx_trust_scores_score ON trust_scores(score);
CREATE INDEX idx_trust_scores_level ON trust_scores(level);

-- Policy lookup by scope
CREATE INDEX idx_policies_scope ON policies(scope);
CREATE INDEX idx_policies_namespace_id ON policies(namespace_id);
CREATE INDEX idx_policies_repository_id ON policies(repository_id);

-- Policy evaluation lookup
CREATE INDEX idx_policy_evaluations_artifact_id ON policy_evaluations(artifact_id);
CREATE INDEX idx_policy_evaluations_policy_id ON policy_evaluations(policy_id);
CREATE INDEX idx_policy_evaluations_evaluated_at ON policy_evaluations(evaluated_at);

-- Metric queries
CREATE INDEX idx_metrics_service_metric_time ON metrics(service_name, metric_name, timestamp);
CREATE INDEX idx_metrics_tags_gin ON metrics USING GIN(tags);

-- Log queries
CREATE INDEX idx_log_entries_service_level_time ON log_entries(service_name, level, timestamp);
CREATE INDEX idx_log_entries_attributes_gin ON log_entries USING GIN(attributes);
CREATE INDEX idx_log_entries_trace_span ON log_entries(trace_id, span_id);

-- Trace queries
CREATE INDEX idx_traces_service_operation_time ON traces(service_name, operation_name, start_time);
CREATE INDEX idx_traces_tags_gin ON traces USING GIN(tags);

-- Span queries
CREATE INDEX idx_spans_trace_service_time ON spans(trace_id, service_name, operation_name, start_time);
CREATE INDEX idx_spans_parent_span ON spans(parent_span_id);

-- Webhook delivery queries
CREATE INDEX idx_webhook_deliveries_webhook_status_time ON webhook_deliveries(webhook_id, status, attempted_at);

-- Audit queries
CREATE INDEX idx_audit_events_tenant_time ON audit_events(tenant_id, timestamp);
CREATE INDEX idx_audit_events_action_resource ON audit_events(action, resource_type, resource_id);
```

## Partitioning Strategy

### Time-Based Partitioning for High-Volume Tables
For tables with high write volumes and time-based queries, we implement partitioning:

```sql
-- Partition metrics table by time (monthly partitions)
CREATE TABLE metrics_y2026m07 PARTITION OF metrics
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');

CREATE TABLE metrics_y2026m08 PARTITION OF metrics
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');

-- Similar partitions for log_entries, webhook_deliveries, audit_events
```

### Partition Management
- Automated partition creation for future months
- Partition retention policy (e.g., keep 24 months of data)
- Archiving strategy for older partitions

## Multi-tenancy Implementation

### Tenant Isolation Approaches
1. **Database Level**: Separate schemas or databases per tenant (for strict isolation requirements)
2. **Schema Level**: Shared tables with tenant_id foreign key (current approach)
3. **Application Level**: Tenant context in application logic

### Current Approach (Schema Level)
- All tables include `tenant_id` column where appropriate for multi-tenancy
- Row-level security policies can be enabled for additional protection
- Application enforces tenant context in all queries

### Row-Level Security Example
```sql
-- Enable RLS on tables
ALTER TABLE repositories ENABLE ROW LEVEL SECURITY;

-- Create policy for tenant isolation
CREATE POLICY tenant_isolation ON repositories
    USING (tenant_id = current_setting('app.current_tenant')::uuid);

-- Set tenant context at session start
SET app.current_tenant = 'tenant-uuid-here';
```

## Backup and Recovery Strategy

### Logical Backups
- **Frequency**: Daily full backups, hourly incremental
- **Tool**: `pg_dump` with custom format
- **Retention**: 30 days daily, 12 weeks weekly, 12 months monthly
- **Verification**: Regular restore tests

### Physical Backups
- **Frequency**: Continuous archiving with WAL-E/WAL-G
- **Storage**: Object storage (S3/MinIO) with encryption
- **Point-in-Time Recovery**: Capable of restoring to any point within retention window

### Disaster Recovery
- **RTO**: < 4 hours for critical services
- **RPO**: < 1 hour for user data
- **Standby**: Warm standby in secondary region
- **Failover**: Automated with health checks

## Schema Evolution and Migrations

### Migration Framework
- **Tool**: Custom migration system or Flyway/Liquibase
- **Versioning**: Sequential numbering with timestamps
- **Transactions**: Each migration runs in a transaction
- **Rollback**: Down migrations provided for reversible changes

### Migration Example
```sql
-- V2026071901_create_trust_tables.sql
BEGIN;

CREATE TABLE trust_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL UNIQUE REFERENCES artifacts(id) ON DELETE CASCADE,
    score FLOAT NOT NULL CHECK (score >= 0.0 AND score <= 1.0),
    level VARCHAR(50) NOT NULL,
    factors JSONB,
    policy_id UUID NULL REFERENCES policies(id) ON DELETE SET NULL,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_trust_scores_artifact_id ON trust_scores(artifact_id);
CREATE INDEX idx_trust_scores_score ON trust_scores(score);

COMMIT;
```

### Data Backfilling
For migrations requiring data transformation:
1. Add new column with nullable default
2. Backfill data in batches using background workers
3. Add NOT NULL constraint after backfill complete
4. Update application code to use new column

## Performance Optimization

### Connection Pooling
- **Tool**: PgBouncer in transaction pooling mode
- **Pool Size**: Configured based on expected concurrent connections
- **Timeouts**: Statement and idle timeouts to prevent resource exhaustion

### Query Optimization
- **EXPLAIN ANALYZE**: Regular query performance review
- **Index Usage**: Ensure queries use appropriate indexes
- **Statistics**: Regular ANALYZE to keep planner statistics up to date
- **Vacuum**: Regular VACUUM to prevent bloat and wraparound issues

### Caching Layers
- **L1**: Application-level caching for frequently accessed reference data
- **L2**: Redis for distributed caching of computed values
- **L3**: HTTP/CDN caching for static assets and public blobs

### Read Replicas
- **Purpose**: Offload read queries from primary
- **Configuration**: Asynchronous replication with lag monitoring
- **Usage**: Direct read-only queries to replicas (reports, dashboards)
- **Fallback**: Automatic failover to primary if replica unavailable

## Security Considerations

### Data Protection
- **Encryption at Rest**:
  - Transparent Data Encryption (TDE) at filesystem level
  - Column-level encryption for sensitive data (tokens, secrets)
  - Backup encryption
  
- **Encryption in Transit**:
  - SSL/TLS for all database connections
  - Connection enforcement via `hostssl` in pg_hba.conf
  
- **Secrets Management**:
  - Database passwords stored in external secret manager (Vault, Azure Key Vault)
  - Application credentials use IAM roles or managed identities
  
### Access Control
- **Authentication**:
  - Certificate authentication for service-to-service connections
  - MD5 or SCRAM-SHA-256 for user authentication
  - PAM or LDAP integration for enterprise environments
  
- **Authorization**:
  - Role-based access control within PostgreSQL
  - Minimal privileges principle for application users
  - Schema-level permissions for isolation
  
### Auditing
- **Extension**: pgAudit for comprehensive audit logging
- **Configuration**: Log all DDL and selected DML operations
- **Storage**: Audit logs sent to centralized logging system
- **Retention**: Configurable retention aligned with compliance requirements

## Monitoring and Maintenance

### Key Metrics to Monitor
- **Connection Count**: Active, idle, waiting connections
- **Query Performance**: Slow queries, execution times
- **Cache Hit Ratio**: Buffer cache and index effectiveness
- **Disk Usage**: Storage utilization and growth trends
- **Replication Lag**: For standby replicas
- **Transaction Rates**: Commits, rollbacks, conflicts
- **Lock Waits**: Blocking and deadlock situations

### Maintenance Tasks
- **Vacuum**: Regular autovacuum configuration monitoring
- **Analyze**: Regular statistics updates
- **Backup Verification**: Regular test restores
- **Index Maintenance**: REINDEX for bloated indexes
- **Log Rotation**: PostgreSQL log file rotation
- **Upgrade Planning**: Tested upgrade procedures for minor/major versions

### Health Checks
- **Extension**: pg_isready or custom health check queries
- **Frequency**: Every 30 seconds via monitoring system
- **Checks**:
  - Database accepting connections
  - Replication lag within thresholds
  - No critical errors in logs
  - Disk space above minimum threshold
  - Successful simple query execution

## Data Archiving and Purging

### Archiving Strategy
- **Old Metrics/Logs**: Move to cheaper storage after 90 days
- **Audit Events**: Retain for 7 years (configurable by compliance needs)
- **Webhook Deliveries**: Keep 90 days of delivery history
- **Artifact Metadata**: Permanent (unless explicitly deleted)

### Purging Procedures
- **Soft Delete**: Most entities use soft deletes with `deleted_at` timestamp
- **Hard Delete**: 
  - Blobs: During garbage collection when unreferenced
  - Temporary data: Based on TTL (e.g., failed webhook attempts after 7 days)
  - Archived data: After retention period expires

### Archival Implementation
```sql
-- Archive old metrics to cold storage
INSERT INTO metrics_archive SELECT * FROM metrics 
    WHERE timestamp < NOW() - INTERVAL '90 days';

DELETE FROM metrics 
    WHERE timestamp < NOW() - INTERVAL '90 days';
```

## Integration with Other Systems

### Event Sourcing Considerations
While PostgreSQL is the primary store, Kyros uses an event-driven architecture:
- **Events Source**: NATS JetStream as the system of record for events
- **Event Consumers**: Services update PostgreSQL as read projections
- **CQRS**: Commands modify event store, queries read from PostgreSQL
- **Rebuilding**: Ability to rebuild PostgreSQL projections from event stream

### Search and Analytics
- **Primary Search**: PostgreSQL for structured data queries
- **Full-Text Search**: Elasticsearch for log and audit trail search
- **Analytics**: Data warehouse (Snowflake, BigQuery, etc.) for business intelligence
- **ETL**: Change data capture (Debezium) or batch exports for analytics pipelines

### Cache Integration
- **Write-Through**: Critical data written to both PostgreSQL and cache
- **Read-Through**: Cache populated on cache miss from PostgreSQL
- **Invalidation**: Event-driven cache invalidation when data changes
- **Warming**: Pre-load frequently accessed data during startup

## Limitations and Trade-offs

### Current Limitations
1. **Vertical Scaling**: Single PostgreSQL instance limits write throughput
2. **Geographic Distribution**: Single region deployment increases latency for distant users
3. **Complex Queries**: Joins across large tables can impact performance
4. **Blob Storage**: Not ideal for storing large binary objects (handled by MinIO/S3)

### Mitigation Strategies
1. **Read Scaling**: Read replicas for horizontal read scaling
2. **Connection Pooling**: PgBouncer to manage connection counts
3. **Query Optimization**: Proper indexing and query design
4. **Blob Offloading**: Large binaries stored in object storage with metadata in DB
5. **Caching**: Redis layer to reduce database load
6. **Future Sharding**: Plan for Citus or similar for horizontal write scaling

### Future Enhancements
1. **Multi-Master Replication**: For geographic distribution and write scaling
2. **Columnar Storage**: For analytics workloads (extension like cstore_fdw)
3. **In-Memory Options**: For caching layers (Redis or PostgreSQL extension)
4. **Time-Series Optimization**: Specialized storage for metrics (TimescaleDB extension)
5. **Document Storage**: JSONB improvements or document store extension for flexible schemas

## Diagrams Reference
See [MERMAID.md](MERMAID.md) for detailed Mermaid diagrams including:
- Entity Relationship Diagrams
- Table Relationships
- Indexing Strategy
- Partitioning Scheme
- Multi-tenancy Implementation
- Backup and Recovery Architecture