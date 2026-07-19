# Kyros Event Architecture

## Overview
Kyros implements an event-driven architecture using NATS JetStream as the central event streaming platform. This enables loose coupling between services, high scalability, resilience, and real-time responsiveness. Events represent significant business occurrences that other services may need to react to.

## Core Concepts

### Event-Driven Architecture Principles
1. **Loose Coupling**: Services communicate through events, not direct calls
2. **Real-Time Responsiveness**: Immediate reaction to business events
3. **Scalability**: Independent scaling of event producers and consumers
4. **Resilience**: Event buffering during consumer downtime
5. **Auditability**: Immutable event log for compliance and debugging
6. **Replayability**: Ability to reprocess events for recovery or new consumers

### NATS JetStream Selection
Kyros chose NATS JetStream for its event streaming backbone because:
- **High Performance**: Low latency, high throughput messaging
- **At-Least-Once Delivery**: Guaranteed event delivery with deduplication
- **Persistent Storage**: Disk-based storage with memory caching
- **Scalability**: Horizontal scaling of consumers and producers
- **Built-in Replay**: Ability to replay events from any point in time
- **Schema Support**: Schema validation for event contracts
- **Flow Control**: Consumer flow control to prevent overload
- **Mirroring**: Event stream mirroring for geo-distribution
- **TLS Support**: Built-in encryption for secure communication
- **Lightweight**: Minimal resource footprint
- **Cloud Native**: Designed for Kubernetes deployment

## Event Streaming Architecture

### NATS JetStream Components
```mermaid
graph TD
    subgraph Kyros Services
        A[Registry Service] -->|Publishes Events| B[NATS JetStream]
        C[API Service] -->|Publishes Events| B
        D[Trust Score Service] -->|Publishes Events| B
        E[Webhook Service] -->|Publishes Events| B
        F[Auth Service] -->|Publishes Events| B
        G[Operator Service] -->|Publishes Events| B
    end
    
    subgraph NATS JetStream
        B --> H[Streams]
        H --> I[Object Store]
        H --> J[Key Value Store]
        H --> K[Consumers]
    end
    
    subgraph Service Consumers
        L[Registry Service] <--|Consumes Events| K
        M[API Service] <--|Consumes Events| K
        N[Trust Score Service] <--|Consumes Events| K
        O[Webhook Service] <--|Consumes Events| K
        P[Auth Service] <--|Consumes Events| K
        Q[Operator Service] <--|Consumes Events| K
        R[External Systems] <--|Consumes Events| K
    end
```

### Core JetStream Features Used
1. **Streams**: Logical grouping of related events with retention policies
2. **Consumers**: Entities that receive and process events from streams
3. **Push/Pull Consumers**: Both asynchronous push and synchronous pull modes
4. **Durable Consumers**: Named consumers that maintain processing state
5. **Ephemeral Consumers**: Temporary consumers for transient subscriptions
6. **Acknowledgments**: Explicit ack/nack for reliable processing
7. **Flow Control**: Consumer-initiated flow control to prevent overload
8. **Schema Validation**: Optional validation of event structure
9. **Mirroring**: Stream replication for disaster recovery and geo-distribution
10. **Consumer Groups**: Load balancing across multiple consumer instances

## Event Streams

Kyros organizes events into logical streams based on domain and functionality:

### 1. Registry Events Stream (`REGISTRY_EVENTS`)
Events related to container image and artifact operations:
- **Subject Prefix**: `registry.>`
- **Retention**: Limits-based (max 1GB) or age-based (24 hours)
- **Storage**: File-based
- **Replicas**: 1 (can be increased for HA)
- **Schema**: Enforced for event validation
- **Max Consumers**: -1 (unlimited)
- **Max Ack Pending**: -1 (unlimited)

#### Event Types:
- `registry.artifact.pushed` - New artifact pushed to registry
- `registry.artifact.pulled` - Artifact pulled from registry
- `registry.artifact.deleted` - Artifact deleted (soft or hard)
- `registry.tag.created` - New tag created
- `registry.tag.deleted` - Tag deleted
- `registry.blob.stored` - Blob stored in blob storage
- `registry.blob.deleted` - Blob deleted during garbage collection
- `registry.repository.created` - New repository created
- `repository.repository.deleted` - Repository deleted
- `registry.namespace.created` - New namespace created
- `registry.namespace.deleted` - Namespace deleted

### 2. Trust Score Events Stream (`TRUST_SCORE_EVENTS`)
Events related to trust score calculation and updates:
- **Subject Prefix**: `trustscore.>`
- **Retention**: Limits-based (max 500MB) or age-based (1 hour)
- **Storage**: Memory-based (for fast access)
- **Replicas**: 1
- **Schema**: Enforced
- **Max Consumers**: -1
- **Max Ack Pending**: 1000

#### Event Types:
- `trustscore.calculated` - Trust score calculated for artifact
- `trustscore.updated` - Trust score updated for artifact
- `trustscore.policy.evaluate` - Policy evaluation requested
- `trustscore.policy.result` - Policy evaluation completed
- `trustscore.sbom.generated` - SBOM generated for artifact
- `trustscore.vulnerability.found` - Vulnerability discovered
- `trustscore.signature.verified` - Signature verified as valid
- `trustscore.signature.failed` - Signature verification failed

### 3. Webhook Events Stream (`WEBHOOK_EVENTS`)
Events related to webhook management and delivery:
- **Subject Prefix**: `webhook.>`
- **Retention**: Limits-based (max 200MB) or age-based (6 hours)
- **Storage**: File-based
- **Replicas**: 1
- **Schema**: Enforced
- **Max Consumers**: -1
- **Max Ack Pending**: 500

#### Event Types:
- `webhook.created` - New webhook subscription created
- `webhook.updated` - Webhook subscription updated
- `webhook.deleted` - Webhook subscription deleted
- `webhook.triggered` - Webhook triggered by event
- `webhook.delivery.success` - Webhook delivery successful
- `webhook.delivery.failed` - Webhook delivery failed after retries
- `webhook.delivery.retry` - Webhook delivery attempt retried

### 4. Authentication Events Stream (`AUTH_EVENTS`)
Events related to user authentication and authorization:
- **Subject Prefix**: `auth.>`
- **Retention**: Limits-based (max 300MB) or age-based (7 days)
- **Storage**: File-based
- **Replicas**: 1
- **Schema**: Enforced
- **Max Consumers**: -1
- **Max Ack Pending**: 200

#### Event Types:
- `auth.user.login` - User successfully authenticated
- `auth.user.logout` - User session ended
- `auth.user.failed_login` - Failed authentication attempt
- `auth.user.created` - New user account created
- `auth.user.updated` - User account updated
- `auth.user.deleted` - User account deleted
- `auth.user.disabled` - User account disabled
- `auth.user.enabled` - User account re-enabled
- `auth.mfa.enrolled` - MFA enrolled for user
- `auth.mfa.verified` - MFA verification successful
- `auth.mfa.challenge` - MFA challenge issued
- `auth.token.issued` - Access token issued
- `auth.token.refreshed` - Access token refreshed
- `auth.token.revoked` - Token revoked
- `auth.role.assigned` - Role assigned to user/group
- `auth.role.revoked` - Role revoked from user/group

### 5. Audit Events Stream (`AUDIT_EVENTS`)
Immutable audit trail for compliance and security:
- **Subject Prefix**: `audit.>`
- **Retention**: Limits-based (max 2GB) or age-based (365 days)
- **Storage**: File-based with compression
- **Replicas**: 3 (for HA and durability)
- **Schema**: Strictly enforced
- **Max Consumers**: -1
- **Max Ack Pending**: -1
- **No Delete**: True (immutable audit log)
- **Max Msg Size**: 1MB

#### Event Types:
- `audit.event.logged` - Audit event recorded
- `audit.retention.expired` - Audit event deleted due to retention
- `audit.integrity.check` - Audit log integrity verification
- `audit.access.granted` - Audit log access granted
- `audit.access.denied` - Audit log access denied

### 6. System Events Stream (`SYSTEM_EVENTS`)
Internal system and operational events:
- **Subject Prefix**: `system.>`
- **Retention**: Limits-based (max 100MB) or age-based (1 hour)
- **Storage**: Memory-based
- **Replicas**: 1
- **Schema**: Enforced
- **Max Consumers**: -1
- **Max Ack Pending**: 100

#### Event Types:
- `service.started` - Service started successfully
- `service.stopped` - Service stopped gracefully
- `service.failed` - Service failed unexpectedly
- `service.health.changed` - Service health status changed
- `config.updated` - Service configuration updated
- `schema.migrated` - Database schema migrated
- `backup.completed` - Backup operation completed
- `backup.failed` - Backup operation failed
- `gc.started` - Garbage collection started
- `gc.completed` - Garbage collection completed
- `gc.failed` - Garbage collection failed

## Event Schema Definition

All Kyros events follow a common schema structure for consistency and interoperability:

### Common Event Envelope
```json
{
  "specversion": "1.0",
  "id": "unique-event-id",
  "source": "service-name.instance-id",
  "type": "domain.event.action",
  "time": "ISO-8601 timestamp",
  "datacontenttype": "application/json",
  "dataschema": "schema-url-or-reference",
  "subject": "resource-identifier",
  "kyrosversion": "1.0.0",
  "traceid": "trace-id-for-distributed-tracing",
  "spanid": "span-id-for-distributed-tracing",
  "tenantid": "tenant-uuid",
  "userid": "user-uuid (if applicable)",
  "data": {
    // Event-specific payload
  }
}
```

### Field Descriptions
- **specversion**: CloudEvents specification version (always "1.0")
- **id**: Unique identifier for the event (UUID v4 recommended)
- **source**: Identifier of the event producer (service name + instance ID)
- **type**: Event type in format `domain.event.action` (e.g., `registry.artifact.pushed`)
- **time**: Timestamp of when the event occurred (ISO 8601 UTC)
- **datacontenttype**: MIME type of the data payload (usually application/json)
- **dataschema**: Reference to the schema definition for the data payload
- **subject**: Resource identifier that the event pertains to (optional)
- **kyrosversion**: Version of the Kyros service that produced the event
- **traceid**: Distributed tracing trace ID for correlation
- **spanid**: Distributed tracing span ID for correlation
- **tenantid**: Tenant identifier for multi-tenancy (optional)
- **userid**: User identifier if event is user-initiated (optional)
- **data**: Event-specific payload following the defined schema

### Example Event: Registry Artifact Pushed
```json
{
  "specversion": "1.0",
  "id": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "source": "registry-service.registry-7b9c5f4d6d6d6d-kslx2",
  "type": "registry.artifact.pushed",
  "time": "2023-07-19T15:04:05.123456Z",
  "datacontenttype": "application/json",
  "dataschema": "https://schemas.kyros.example.com/v1/events/registry/artifact-pushed.json",
  "subject": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "kyrosversion": "1.0.0",
  "traceid": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "spanid": "b2c3d4e5-f6g7-8901-h2i3-j4k5l6m7n8o9",
  "tenantid": "123e4567-e89b-12d3-a456-426614174000",
  "userid": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "data": {
    "artifact": {
      "digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "size": 1234,
      "repository": {
        "id": "b2c3d4e5-f6g7-8901-h2i3-j4k5l6m7n8o9",
        "name": "my-app",
        "namespace": "production"
      },
      "tags": [
        {
          "name": "latest",
          "createdAt": "2023-07-19T15:04:05.123Z"
        }
      ],
      "uploadedBy": {
        "id": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
        "username": "admin-user"
      },
      "createdAt": "2023-07-19T15:04:05.123Z"
    },
    "manifest": {
      "schemaVersion": 2,
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": 1469,
        "digest": "sha256:9c6d842013... (truncated)"
      },
      "layers": [
        {
          "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
          "size": 3265402,
          "digest": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
        }
      ]
    }
  }
}
```

## Event Producer Implementation

Services publish events using the Kyros Event Publisher library:

### Publisher Interface
```go
type Publisher interface {
    // PublishEvent publishes an event to the specified stream
    PublishEvent(ctx context.Context, event *Event) error
    
    // PublishEventWithSubject publishes an event with explicit subject
    PublishEventWithSubject(ctx context.Context, subject string, event *Event) error
    
    // PublishBatch publishes multiple events in a batch
    PublishBatch(ctx context.Context, events []*Event) error
    
    // Close closes the publisher and releases resources
    Close() error
}
```

### Event Structure
```go
type Event struct {
    SpecVersion   string                 `json:"specversion"`
    ID            string                 `json:"id"`
    Source        string                 `json:"source"`
    Type          string                 `json:"type"`
    Time          time.Time              `json:"time"`
    DataContentType string               `json:"datacontenttype,omitempty"`
    DataSchema    string                 `json:"dataschema,omitempty"`
    Subject       string                 `json:"subject,omitempty"`
    KyrosVersion  string                 `json:"kyrosversion"`
    TraceID       string                 `json:"traceid,omitempty"`
    SpanID        string                 `json:"spanid,omitempty"`
    TenantID      string                 `json:"tenantid,omitempty"`
    UserID        string                 `json:"userid,omitempty"`
    Data          map[string]interface{} `json:"data"`
}
```

### Publishing Example
```go
func (s *RegistryService) handleArtifactPush(ctx context.Context, req *registry.PushRequest) error {
    // ... process the push ...
    
    // Create and publish event
    event := &events.Event{
        ID:           uuid.NewString(),
        Source:       fmt.Sprintf("registry-service.%s", s.instanceID),
        Type:         "registry.artifact.pushed",
        Time:         time.Now().UTC(),
        DataContentType: "application/json",
        DataSchema:   "https://schemas.kyros.example.com/v1/events/registry/artifact-pushed.json",
        Subject:      artifactDigest,
        KyrosVersion: s.version,
        TraceID:      traceIDFromContext(ctx),
        SpanID:       spanIDFromContext(ctx),
        TenantID:     tenantIDFromContext(ctx),
        UserID:       userIDFromContext(ctx),
        Data: map[string]interface{}{
            "artifact": artifactInfo,
            "manifest": manifestInfo,
        },
    }
    
    if err := s.eventPublisher.PublishEvent(ctx, event); err != nil {
        // Log error but don't fail the operation
        s.logger.Error("Failed to publish artifact pushed event", zap.Error(err))
    }
    
    return nil
}
```

## Event Consumer Implementation

Services consume events using the Kyros Event Consumer library:

### Consumer Interface
```go
type Consumer interface {
    // StartConsuming begins consuming events from the specified stream
    StartConsuming(ctx context.Context, opts *ConsumerOptions) (<-chan *Event, <-chan error, error)
    
    // StopConsuming stops consuming events
    StopConsuming() error
    
    // Close closes the consumer and releases resources
    Close() error
}
```

### Consumer Options
```go
type ConsumerOptions struct {
    // Stream to consume from
    Stream string
    
    // Subjects to filter (wildcards supported)
    Subjects []string
    
    // Consumer configuration
    Durable   string // Durable consumer name (empty for ephemeral)
    Group     string // Consumer group name for load balancing
    AckPolicy  AckPolicy // Explicit, All, None
    FlowControl bool // Enable flow control
    Heartbeat  time.Duration // Heartbeat interval
    MaxAckPending int // Maximum unacknowledged messages
    
    // Processing configuration
    BatchSize  int // Number of events to batch
    BatchTimeout time.Duration // Max time to wait for batch
    
    // Error handling
    MaxRetries int // Maximum retry attempts
    RetryDelay time.Duration // Delay between retries
    DeadLetterEnabled bool // Send failed events to dead letter
}
```

### Consuming Example
```go
func (s *TrustScoreService) StartEventConsuming(ctx context.Context) error {
    eventsChan, errorsChan, err := s.eventConsumer.StartConsuming(ctx, &events.ConsumerOptions{
        Stream:   "REGISTRY_EVENTS",
        Subjects: []string{"registry.artifact.pushed"},
        Durable:  "trustscore-service",
        Group:    "trustscore-workers",
        AckPolicy: events.AckExplicit,
        FlowControl: true,
        MaxAckPending: 1000,
        BatchSize: 10,
        BatchTimeout: time.Second,
        MaxRetries: 3,
        RetryDelay: time.Second * 5,
        DeadLetterEnabled: true,
    })
    
    if err != nil {
        return fmt.Errorf("failed to start event consumer: %w", err)
    }
    
    // Process events in goroutine
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case event := <-eventsChan:
                s.processRegistryEvent(event)
            case err := <-errorsChan:
                s.logger.Error("Error consuming events", zap.Error(err))
                // Depending on error type, may want to restart consumer
            }
        }
    }()
    
    return nil
}

func (s *TrustScoreService) processRegistryEvent(event *events.Event) {
    switch event.Type {
    case "registry.artifact.pushed":
        s.handleArtifactPushed(event)
    // ... other event types ...
    }
}

func (s *TrustScoreService) handleArtifactPushed(event *events.Event) {
    // Extract artifact information from event data
    artifactData := event.Data["artifact"].(map[string]interface{})
    digest := artifactData["digest"].(string)
    
    // Initiate trust score calculation
    s.calculateTrustScore(digest, event)
}
```

## Event Processing Guarantees

### Delivery Semantics
Kyros uses **at-least-once** delivery semantics:
- Events may be delivered more than once but never lost
- Consumers must be idempotent to handle duplicate events
- NATS JetStream provides deduplication at the stream level when configured

### Idempotency Patterns
Event consumers implement idempotency through:
1. **Event ID Tracking**: Recording processed event IDs to skip duplicates
2. **State-Based Processing**: Checking if action already performed based on current state
3. **Version Vectors**: Using version numbers or timestamps to detect stale events
4. **External Deduplication**: Leveraging external system's idempotency (e.g., database constraints)

#### Example Idempotency Implementation
```go
func (s *TrustScoreService) handleArtifactPushed(event *events.Event) {
    // Check if we've already processed this event
    if s.hasProcessedEvent(event.ID) {
        s.logger.Debug("Skipping duplicate event", 
                      zap.String("event_id", event.ID),
                      zap.String("event_type", event.Type))
        return
    }
    
    // Process the event
    digest := extractDigest(event.Data)
    if err := s.calculateTrustScore(digest, event); err != nil {
        s.logger.Error("Failed to calculate trust score", 
                      zap.Error(err),
                      zap.String("event_id", event.ID))
        // Don't mark as processed on failure to allow retry
        return
    }
    
    // Mark event as processed
    s.markEventProcessed(event.ID)
}
```

### Ordering Guarantees
- **Per-Key Ordering**: Events with same subject key are processed in order
- **Global Ordering**: No global ordering guarantee across different keys
- **Consumer Group Ordering**: Within a consumer group, ordering is maintained per key
- **Application-Level Ordering**: Applications can implement sequencing when needed

## Error Handling and Retry Logic

### Consumer-Level Retry
Consumers implement retry logic for transient failures:
```go
func (s *TrustScoreService) processEventWithRetry(event *events.Event) error {
    var lastErr error
    for attempt := 0; eventConsumerOptions.MaxRetries >= attempt; attempt++ {
        err := s.processEvent(event)
        if err == nil {
            // Success - acknowledge event
            return nil
        }
        
        lastErr = err
        s.logger.Warn("Event processing failed, retrying",
                     zap.Int("attempt", attempt+1),
                     zap.Int("max_attempts", eventConsumerOptions.MaxRetries+1),
                     zap.Error(err),
                     zap.String("event_id", event.ID))
        
        if attempt < eventConsumerOptions.MaxRetries {
            // Wait before retry (with exponential backoff + jitter)
            delay := calculateBackoffDelay(attempt, eventConsumerOptions.RetryDelay)
            time.Sleep(delay)
        }
    }
    
    // All retries exhausted
    return fmt.Errorf("failed to process event after %d attempts: %w", 
                     eventConsumerOptions.MaxRetries+1, lastErr)
}
```

### Dead Letter Queue
Failed events can be routed to a dead letter queue for later inspection:
```go
// In consumer options
DeadLetterEnabled: true
DeadLetterStream: "DLQ_REGISTRY_EVENTS"
DeadLetterSubject: "dlq.registry.events"

// On repeated failures:
// 1. Event is not acknowledged
// 2. After max retries, JetStream redirects to dead letter stream
// 3. Separate consumer processes dead letter events for analysis
```

### Circuit Breaker Pattern
To prevent cascading failures:
```go
// Circuit breaker state
type CircuitBreaker struct {
    state      State // Closed, Open, HalfOpen
    failureCount int
    lastFailureTime time.Time
    successThreshold uint32
    failureThreshold uint32
    timeout time.Duration
}

// Before processing event
if cb.State() == StateOpen {
    if time.Since(cb.LastFailureTime()) > cb.Timeout {
        cb.SetState(StateHalfOpen)
    } else {
        return ErrCircuitBreakerOpen
    }
}

// After processing event
if err == nil {
    cb.OnSuccess()
} else {
    cb.OnFailure()
}
```

## Schema Management and Evolution

### Event Schema Registry
Kyros maintains a schema registry for event validation:
- **Storage**: ConfigMap or external registry (Apicurio, Confluent Schema Registry, etc.)
- **Versioning**: Semantic versioning for schema evolution
- **Compatibility**: Backward and forward compatibility checking
- **Validation**: Optional schema validation at publish time
- **Documentation**: Generated documentation from schemas

### Schema Definition Example
```json
{
  "$id": "https://schemas.kyros.example.com/v1/events/registry/artifact-pushed.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Registry Artifact Pushed Event",
  "type": "object",
  "required": ["specversion", "id", "source", "type", "time", "data"],
  "properties": {
    "specversion": {
      "type": "string",
      "const": "1.0"
    },
    "id": {
      "type": "string",
      "format": "uuid"
    },
    "source": {
      "type": "string",
      "pattern": "^registry-service\\.[a-zA-Z0-9-]+$"
    },
    "type": {
      "type": "string",
      "pattern": "^registry\\.artifact\\.pushed$"
    },
    "time": {
      "type": "string",
      "format": "date-time"
    },
    "datacontenttype": {
      "type": "string",
      "enum": ["application/json"]
    },
    "dataschema": {
      "type": "string",
      "pattern": "^https://schemas\\.kyros\\.example\\.com/v1/events/"
    },
    "subject": {
      "type": "string",
      "pattern": "^sha256:[a-f0-9]{64}$"
    },
    "kyrosversion": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+\\.\\d+$"
    },
    "traceid": {
      "type": ["string", "null"],
      "format": "uuid"
    },
    "spanid": {
      "type": ["string", "null"],
      "format": "uuid"
    },
    "tenantid": {
      "type": ["string", "null"],
      "format": "uuid"
    },
    "userid": {
      "type": ["string", "null"],
      "format": "uuid"
    },
    "data": {
      "type": "object",
      "required": ["artifact": { "$ref": "https://schemas.kyros.example.com/v1/schemas/registry/artifact.json"
    }
  }
}
```

### Schema Evolution Strategies
1. **Backward Compatible Changes**:
   - Add new optional fields
   - Add new enum values
   - Make required fields optional (with default values)
   
2. **Forward Compatible Changes**:
   - Remove optional fields (if consumers ignore unknown fields)
   - Make optional fields required (if producers already provide them)
   
3. **Breaking Changes**:
   - Require new stream or subject
   - Change event type semantics
   - Remove required fields
   - Change data types in incompatible ways

### Versioning Approach
- **Event Type Versioning**: Include version in event type (`registry.v2.artifact.pushed`)
- **Schema Reference Versioning**: Reference specific schema versions in `dataschema`
- **Content Versioning**: Include version in data payload
- **Consumer Adaptation**: Consumers handle multiple versions of same event type

## Security Considerations

### Transport Security
- **TLS Encryption**: All NATS connections encrypted with TLS 1.2+
- **Authentication**: JWT or NKEY authentication for service connections
- **Authorization**: Role-based access control for stream/subject permissions
- **Message Encryption**: Optional end-to-end encryption for sensitive events

### Authentication and Authorization
1. **Connection Authentication**:
   - Services authenticate to NATS using service accounts
   - JWT tokens or NKEY pairs for identification
   - Short-lived credentials with rotation
   
2. **Stream/Subject Authorization**:
   - Permissions defined per stream or subject pattern
   - Publish/subscribe permissions granted separately
   - Wildcard subject matching for efficient permission definition
   
3. **Claims-Based Access**:
   - JWT claims used for fine-grained authorization
   - Tenant and user information extracted from tokens
   - Event producers can only publish events for their tenant/user

### Event Integrity
- **Schema Validation**: Optional validation against registered schemas
- **Message Signing**: Optional cryptographic signing of events
- **Hash Chaining**: Optional hash chaining for tamper evidence
- **Timestamp Validation**: Rejection of events with skewed timestamps

### Audit and Compliance
- **Immutable Audit Stream**: Separate stream for audit events with no deletion
- **Access Logging**: Logging of who published/consumed what events
- **Retention Policies**: Configurable retention with legal hold capabilities
- **Export Capabilities**: Ability to export event streams for external audit

## Performance and Scaling

### Throughput Characteristics
- **Publish Latency**: <1ms for local, <10ms for cross-availability zone
- **Consume Latency**: <1ms for local, <10ms for cross-availability zone
- **Maximum Throughput**: 100,000+ msg/sec per cluster node (dependent on hardware)
- **Event Size**: Efficient for events up to 1MB (larger events use object store)

### Scaling Patterns
#### Producer Scaling
- **Stateless Publishers**: Any service instance can publish events
- **Connection Pooling**: Efficient reuse of NATS connections
- **Bulk Publishing**: Batch multiple events for efficiency
- **Asynchronous Publishing**: Non-blocking event publishing

#### Consumer Scaling
- **Consumer Groups**: Multiple instances share workload via load balancing
- **Dynamic Scaling**: Add/remove consumers based on lag and throughput
- **Prefetching**: Configureable number of events to prefetch for efficiency
- **Ack Strategies**: Explicit, all, or none acknowledgment policies

### Resource Utilization
#### Memory Usage
- **Connection Buffers**: Minimal per-connection memory overhead
- **In-Memory Streams**: Configurable memory allocation for fast streams
- **Consumer State**: Minimal state per consumer (mainly acknowledgment tracking)
- **Batch Buffers**: Configurable batch sizes for memory/throughput tradeoff

#### Disk Usage
- **File-Based Streams**: Disk space proportional to retained event size
- **Indexing**: Minimal indexing overhead for subject-based lookup
- **Compaction**: Automatic compaction of deleted events
- **Rotation**: Segment-based rotation for efficient disk usage

### Monitoring and Metrics
Key metrics to monitor for event streaming health:
```prometheus
# NATS JetStream metrics (if exposed)
nats_js_stream_messages_total{stream="REGISTRY_EVENTS",dir="in"} 12450
nats_js_stream_bytes_total{stream="REGISTRY_EVENTS",dir="in"} 45.2MB
nats_js_stream_messages_total{stream="REGISTRY_EVENTS",dir="out"} 12430
nats_js_stream_bytes_total{stream="REGISTRY_EVENTS",dir="out"} 45.0MB
nats_js_stream_consume_latency_seconds{stream="REGISTRY_EVENTS"} 0.002
nats_js_stream_publish_latency_seconds{stream="REGISTRY_EVENTS"} 0.001
nats_js_stream_ack_wait_seconds{stream="REGISTRY_EVENTS"} 0.005
nats_js_stream_num_consumers{stream="REGISTRY_EVENTS"} 5
nats_js_stream_memory_bytes{stream="REGISTRY_EVENTS"} 10.5MB
nats_js_stream_store_bytes{stream="REGISTRY_EVENTS"} 40.1MB

# Application-level metrics
events_published_total{service="registry",type="registry.artifact.pushed"} 1245
events_consumed_total{service="trustscore",type="registry.artifact.pushed"} 1240
events_failed_total{service="trustscore",type="registry.artifact.pushed",reason="timeout"} 3
events_retried_total{service="trustscore",type="registry.artifact.pushed"} 2
events_dlq_total{stream="REGISTRY_EVENTS"} 2
processing_latency_seconds{service="trustscore",type="registry.artifact.pushed"} 0.85
```

## Integration Patterns

### 1. Service-to-Service Communication
```mermaid
sequenceDiagram
    participant Registry
    participant NATS
    participant TrustScore
    
    Registry->>NATS: Publish registry.artifact.pushed
    NATS->>TrustScore: Deliver to trustscore consumers
    TrustScore->>Registry: Process and calculate trust score
    TrustScore->>NATS: Publish trustscore.calculated
    NATS->>Registry: Deliver to registry consumers
    Registry->>NATS: Publish registry.artifact.tagged (with trust score)
```

### 2. Webhook Integration
```mermaid
sequenceDiagram
    participant Registry
    participant NATS
    participant Webhook
    participant External_System
    
    Registry->>NATS: Publish registry.artifact.pushed
    NATS->>Webhook: Deliver matching webhook subscriptions
    Webhook->>External_System: HTTP POST with event payload
    External_System->>Webhook: HTTP 200 OK
    Webhook->>NATS: Publish webhook.delivery.success
    NATS->>Registry: Deliver to interested services
```

### 3. External System Integration
External systems can integrate with Kyros events through:
- **NATS Clients**: Direct connection to Kyros NATS cluster (with proper auth)
- **Webhook Subscriptions**: Kyros delivers events to external HTTP endpoints
- **Event Export**: Periodic export of event streams to external systems
- **Schema Sharing**: Shared schema definitions for event validation

### 4. Geo-Distribution and Mirroring
```mermaid
graph TD
    subgraph Primary Region
        A[Kyros Services] -->|Publish Events| B[NATS Primary]
        B -->|Mirror Stream| C[NATS Secondary]
    end
    
    subgraph Secondary Region
        D[Kyros Services] <--|Consume Events| C
    end
    
    C -->|Mirror Stream| B
```

## Operational Considerations

### Deployment Guidelines
- **NATS Cluster Size**: 
  - Development: 1 node
  - Staging: 3 nodes
  - Production: 3-5 nodes (odd number for Raft quorum)
- **Resource Allocation per Node**:
  - CPU: 2-4 cores
  - Memory: 4-8 GB RAM
  - Storage: 
    - Memory Streams: Configurable (default 1GB)
    - File Streams: Based on retention policy
    - Indexes: Minimal overhead
- **Storage Type**: 
  - SSD recommended for file-based streams
  - NVMe for high-throughput scenarios
- **Network**: 
  - 1 Gbps minimum, 10 Gbps recommended for high volume
  - Low latency inter-node communication (<1ms ideal)

### Configuration Parameters
#### Server Configuration
```conf
# Server binding
listen: 0.0.0.0:4222

# Clustering
cluster {
  listen: 0.0.0.0:6222
  routes: [
    nats-route://user:password@nats1:6222
    nats-route://user:password@nats2:6222
    nats-route://user:password@nats3:6222
  ]
}

# JetStream configuration
jetstream {
  store_dir: /var/lib/jetstream
  max_mem_store: 1GB
  max_file_store: 10GB
}

# Security
authorization {
  user: jwt
  # or nkeys
}

# Limits
max_payload: 1MB
max_pending: 100000
```

#### Client Configuration (Services)
```go
// Connection options
opts := []nats.Option{
    nats.Name("kyros-registry-service"),
    nats.UserCredentials("/creds/operator.creds"),
    nats.ReconnectWait(time.Second * 2),
    nats.MaxReconnects(-1), // Infinite reconnect attempts
    nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
        log.Printf("Disconnected: %v", err)
    }),
    nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Printf("Reconnected [%s]", nc.ConnectedUrl())
    }),
    nats.ClosedHandler(func(nc *nats.Conn) {
        log.Printf("Connection closed: %v", nc.LastError())
    }),
}

// JetStream context
js, err := nc.JetStream(
    nats.Context(context.Background()),
)
```

### Maintenance Procedures
#### Regular Tasks
- **Monitoring**: Check latency, throughput, and resource utilization
- **Log Rotation**: Rotate and archive NATS server logs
- **Disk Usage**: Monitor disk space for file-based streams
- **Memory Usage**: Monitor memory usage for memory-based streams
- **Connection Count**: Track active client connections
- **Consumer Lag**: Monitor consumer processing lag
- **Stream Health**: Check for stuck consumers or failed acknowledgments

#### Periodic Tasks
- **Software Updates**: Apply NATS server updates and security patches
- **Configuration Review**: Review and update stream configurations
- **Capacity Planning**: Analyze growth trends and plan for scaling
- **Backup Verification**: Verify backup and restore procedures
- **Security Audit**: Review authentication and authorization configurations

#### Incident Response
- **High Latency**: 
  - Check consumer processing speed
  - Look for blocked consumers or network issues
  - Consider scaling consumers or optimizing processing
  
- **Message Loss**:
  - Check stream configuration (retention, max messages)
  - Look for producer errors or network issues
  - Verify acknowledgment processing
  
- **Connection Issues**:
  - Check network connectivity and firewall rules
  - Verify authentication credentials
  - Look for server resource exhaustion
  
- **Disk Space Issues**:
  - Review retention policies and stream limits
  - Consider increasing storage or reducing retention
  - Look for stuck consumers preventing acknowledgment
  
- **Memory Pressure**:
  - Check memory-based stream usage
  - Consider converting to file-based or increasing memory
  - Look for memory leaks in clients

### Backup and Recovery
#### JetStream Backup
- **File-Based Streams**: 
  - Copy stream directory while NATS is running (consistent snapshots)
  - Use filesystem snapshots (LVM, ZFS, etc.) for consistency
  - Backup index files along with data files
  
- **Memory-Based Streams**:
  - Not directly backuplable (ephemeral by design)
  - Configure as file-based if persistence required
  - Export important events to external systems before loss
  
#### Recovery Procedures
- **Node Failure**: 
  - NATS Raft protocol handles automatic failover
  - Replace failed node and rejoin cluster
  - Stream data replicated across surviving nodes
  
- **Data Corruption**:
  - Restore from backup if available
  - Rebuild streams from external sources if possible
  - May lose some recent events depending on backup frequency
  
- **Complete Cluster Loss**:
  - Restore from backup
  - Replay events from external archives if available
  - Accept potential data loss for window between backups

## Observability and Tracing Integration

### Distributed Tracing
Kyros integrates event streaming with distributed tracing:
- **Trace Context Propagation**: 
  - Trace ID and Span ID included in event metadata
  - Consumers extract trace context and create child spans
  - Enables end-to-end tracing across service boundaries
  
- **Event Processing Spans**:
  - Publish span: Time to publish event to NATS
  - Receive span: Time to receive event from NATS
  - Process span: Time to process event business logic
  - Acknowledge span: Time to acknowledge event processing
  
- **Correlation**:
  - All spans for a single business operation share trace ID
  - Enables tracing from user action through multiple services
  - Supports root cause analysis and performance optimization

### Metrics Export
Kyros exports Prometheus metrics for event streaming:
```prometheus
# Event publishing metrics
events_published_total{service="registry",stream="REGISTRY_EVENTS",type="registry.artifact.pushed"} 1245
events_publish_latency_seconds{service="registry",stream="REGISTRY_EVENTS"} 0.001

# Event consumption metrics
events_consumed_total{service="trustscore",stream="REGISTRY_EVENTS",type="registry.artifact.pushed"} 1240
events_consumed_latency_seconds{service="trustscore",stream="REGISTRY_EVENTS"} 0.002
events_consumer_lag{service="trustscore",stream="REGISTRY_EVENTS"} 5

# Event processing metrics
events_processing_latency_seconds{service="trustscore",stream="REGISTRY_EVENTS",type="registry.artifact.pushed"} 0.85
events_processing_total{service="trustscore",stream="REGISTRY_EVENTS",type="registry.artifact.pushed",result="success"} 1200
events_processing_total{service="trustscore",stream="REGISTRY_EVENTS",type="registry.artifact.pushed",result="error"} 40

# JetStream metrics (if exposed via exporter)
nats_js_stream_consumers{stream="REGISTRY_EVENTS"} 3
nats_js_stream_pending{stream="REGISTRY_EVENTS"} 10
nats_js_stream_memory{stream="REGISTRY_EVENTS"} 10500000
nats_js_stream_store{stream="REGISTRY_EVENTS"} 40100000
```

### Logging Integration
Events are integrated with structured logging:
```json
{
  "timestamp": "2023-07-19T15:04:05.123Z",
  "level": "info",
  "message": "Event published",
  "service": "registry",
  "traceID": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "spanID": "b2c3d4e5-f6g7-8901-h2i3-j4k5l6m7n8o9",
  "event": {
    "id": "c3d4e5f6-g7h8-9012-i3j4-k5l6m7n8o9p0",
    "type": "registry.artifact.pushed",
    "stream": "REGISTRY_EVENTS",
    "size": 1024
  }
}
```

```json
{
  "timestamp": "2023-07-19T15:04:06.500Z",
  "level": "info",
  "message": "Event consumed",
  "service": "trustscore",
  "traceID": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "spanID": "c3d4e5f6-g7h8-9012-i3j4-k5l6m7n8o9p0",
  "event": {
    "id": "c3d4e5f6-g7h8-9012-i3j4-k5l6m7n8o9p0",
    "type": "registry.artifact.pushed",
    "stream": "REGISTRY_EVENTS",
    "processing_time_ms": 850
  }
}
```

## Future Enhancements

### Near-Term Improvements (0-6 months)
1. **Schema Registry Integration**: Full integration with Apicurio or similar for schema management
2. **Enhanced Monitoring**: Better integration with Prometheus and Grafana dashboards
3. **Event Replay UI**: Interface for replaying events for debugging or reprocessing
4. **Dead Letter Queue Management**: Tools for inspecting and reprocessing DLQ events
5. **Schema Validation**: Optional schema validation at publish time with configuration
6. **Event Compression**: Optional compression of large events to reduce bandwidth
7. **Batch Publishing**: Native support for batch event publishing for efficiency
8. **Flow Control Tuning**: Dynamic flow control based on consumer processing capacity

### Mid-Term Features (6-18 months)
1. **Event Sourcing**: 
   - Events as system of record for certain domains
   - Ability to rebuild state from event stream
   - Event sourcing patterns for audit trails and historical reconstruction
   
2. **Stream Processing**:
   - Integration with stream processing frameworks (Apache Flink, Kafka Streams)
   - Complex event processing (CEP) for pattern detection
   - Real-time analytics and aggregations on event streams
   
3. **Geo-Distribution Enhancements**:
   - Active-active multi-region deployment with conflict resolution
   - Smart event routing based on consumer location
   - Latency-based event delivery optimization
   
4. **Security Enhancements**:
   - End-to-end event encryption for sensitive domains
   - Event signing and verification for integrity
   - Advanced authorization policies based on event content
   
5. **Developer Experience**:
   - Event schema generation from Go structs
   - Mock event servers for testing
   - Event visualization and debugging tools

### Long-Term Vision (18+ months)
1. **Event Mesh Integration**:
   - Integration with service meshes for advanced traffic management
   - Event-based service discovery and routing
   - Event-driven service level objectives (SLOs)
   
2. **Machine Learning Integration**:
   - Anomaly detection in event streams for threat identification
   - Predictive event processing for proactive scaling
   - Automated root cause analysis using event correlations
   
3. **Semantic Event Processing**:
   - Event enrichment with contextual information
   - Semantic querying and reasoning over event streams
   - Event knowledge graph construction and querying
   
4. **Federated Event Architecture**:
   - Federation of event streams across organizations
   - Trust-based event sharing and consumption
   - Event marketplace for commercial and community events
   
5. **Quantum-Resistant Security**:
   - Post-quantum cryptography for event security
   - Quantum-resistant authentication and authorization
   - Quantum-safe event integrity protection

## Implementation Best Practices

### For Event Producers
1. **Publish After Commit**: Only publish events after business transaction is committed
2. **Include Context**: Always include trace IDs, tenant IDs, and user IDs when available
3. **Validate Before Publishing**: Validate business data before creating events
4. **Handle Publishing Failures**: Log errors but don't fail business transactions for event publishing
5. **Use Appropriate Subjects**: Use hierarchical subjects for effective filtering
6. **Consider Event Size**: Keep events reasonably sized; use references for large payloads
7. **Publish Meaningful Events**: Only publish events that represent significant business occurrences
8. **Maintain Backward Compatibility**: Evolve events carefully to avoid breaking consumers

### For Event Consumers
1. **Be Idempotent**: Design consumers to handle duplicate events safely
2. **Acknowledge Properly**: Only acknowledge events after successful processing
3. **Handle Errors Gracefully**: Implement retry logic and dead letter queue handling
4. **Monitor Lag**: Track and alert on consumer processing lag
5. **Respect Flow Control**: Respond to flow control signals from the event stream
6. **Filter Appropriately**: Use subject filtering to reduce unnecessary event processing
7. **Manage Resources**: Limit concurrent processing to prevent resource exhaustion
8. **Maintain Observability**: Include tracing and metrics in event processing

### For Architecture and Operations
1. **Design for Failure**: Assume events may be delayed, duplicated, or arrive out of order
2. **Plan for Scale**: Design consumers and producers to scale horizontally
3. **Monitor Relentlessly**: Implement comprehensive monitoring for event streaming health
4. **Test Thoroughly**: Test failure scenarios, network partitions, and resource exhaustion
5. **Document Events**: Maintain clear documentation of event types, schemas, and semantics
6. **Version Carefully**: Plan event schema evolution with backward compatibility in mind
7. **Secure Appropriately**: Implement proper authentication, authorization, and encryption
8. **Operate Sustainably**: Implement maintenance procedures and capacity planning

## Conclusion

Kyros' event-driven architecture using NATS JetStream provides a robust, scalable, and resilient foundation for building a modern cloud-native platform. By leveraging events as the primary communication mechanism, Kyros achieves loose coupling between services, enabling independent development, deployment, and scaling.

The architecture supports a wide range of integration patterns from simple service-to-service communication to complex event processing and geo-distribution. With proper implementation, monitoring, and operational practices, the event streaming platform becomes a powerful enabler for Kyros' vision as a world-class cloud-native software supply chain platform.

Through careful attention to event schema design, processing guarantees, error handling, and security considerations, Kyros ensures that its event architecture remains reliable, maintainable, and adaptable to evolving requirements.