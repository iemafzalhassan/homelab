# Kyros Event Catalog

## Overview
This document catalogs all events in the Kyros platform, including their producers, consumers, payloads, schemas, ordering guarantees, retry strategies, dead-letter handling, idempotency considerations, and versioning. Events are the primary communication mechanism between services in Kyros' event-driven architecture.

## Event Naming Convention
Events follow the format: `domain.entity.action`
- **domain**: Business domain (e.g., registry, trustscore, webhook, auth, audit)
- **entity**: Primary resource or object (e.g., artifact, repository, user, policy)
- **action**: What happened (e.g., pushed, pulled, created, updated, deleted, calculated, verified)

## Event Schema
All events adhere to the CloudEvents specification version 1.0 with the following envelope structure:

```json
{
  "specversion": "1.0",
  "id": "unique-event-id",
  "source": "service-name.instance-id",
  "type": "domain.entity.action",
  "time": "ISO-8601 timestamp",
  "datacontenttype": "application/json",
  "dataschema": "schema-url-or-reference",
  "subject": "resource-identifier (optional)",
  "kyrosversion": "1.0.0",
  "traceid": "trace-id-for-distributed-tracing",
  "spanid": "span-id-for-distributed-tracing",
  "tenantid": "tenant-uuid (optional)",
  "userid": "user-uuid (if applicable)",
  "data": {
    // Event-specific payload
  }
}
```

## Event Catalog

### 1. Registry Events
Events related to container image and artifact operations.

#### 1.1 registry.artifact.pushed
- **Producer**: Registry Service
- **Consumers**: 
  - Trust Score Service (triggers scoring)
  - Webhook Service (delivers to subscribers)
  - Audit Service (logs for compliance)
  - API Service (updates search indexes)
- **Payload**:
  ```json
  {
    "artifact": {
      "id": "uuid",
      "digest": "sha256:...",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "size": 1234,
      "repositoryId": "uuid",
      "uploadedBy": {
        "id": "uuid",
        "username": "string"
      },
      "createdAt": "ISO-8601 timestamp"
    },
    "manifest": {
      "schemaVersion": 2,
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": 1469,
        "digest": "sha256:..."
      },
      "layers": [
        {
          "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
          "size": 3265402,
          "digest": "sha256:..."
        }
      ]
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/registry/artifact-pushed.json`
- **Subject**: Artifact digest (e.g., `sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`)
- **Ordering**: Per-artifact (events for same artifact digest are ordered)
- **Retry Strategy**: 
  - Initial attempt + 3 retries with exponential backoff (1s, 2s, 4s)
  - Jitter added to prevent thundering herd
  - After max retries, sent to dead letter queue
- **Dead Letter Handling**: 
  - Sent to `dlq.registry.events` stream
  - Retained for 14 days for inspection
  - Can be replayed manually after issue resolution
- **Idempotency**: 
  - Consumers should track processed event IDs
  - Registry service is idempotent for artifact pushes (same digest)
- **Versioning**: 
  - Schema version in `dataschema` field
  - Backward compatible changes: add optional fields
  - Breaking changes: new event type (e.g., `registry.v2.artifact.pushed`)

#### 1.2 registry.artifact.pulled
- **Producer**: Registry Service
- **Consumers**: 
  - Analytics Service (updates pull metrics)
  - Webhook Service (delivers to subscribers)
  - Audit Service (logs for compliance)
- **Payload**:
  ```json
  {
    "artifact": {
      "id": "uuid",
      "digest": "sha256:...",
      "repositoryId": "uuid",
      "pulledBy": {
        "id": "uuid",
        "username": "string"
      },
      "pulledAt": "ISO-8601 timestamp",
      "ipAddress": "string",
      "userAgent": "string"
    },
    "manifest": {
      "schemaVersion": 2,
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/registry/artifact-pulled.json`
- **Subject**: Artifact digest
- **Ordering**: Per-artifact
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: Pull events are naturally idempotent (same pull multiple times)
- **Versioning**: Same as above

#### 1.3 registry.tag.created
- **Producer**: Registry Service
- **Consumers**: 
  - Webhook Service
  - Audit Service
  - Trust Score Service (may trigger scoring for new tag)
- **Payload**:
  ```json
  {
    "tag": {
      "id": "uuid",
      "name": "string",
      "artifactId": "uuid",
      "createdAt": "ISO-8601 timestamp",
      "createdBy": {
        "id": "uuid",
        "username": "string"
      }
    },
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/registry/tag-created.json`
- **Subject**: Tag name (e.g., `my-app:latest`)
- **Ordering**: Per-repository
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Creating same tag twice: second attempt fails (if immutable) or overwrites (if mutable)
  - Consumers should handle both cases
- **Versioning**: Same as above

#### 1.4 registry.tag.deleted
- **Producer**: Registry Service
- **Consumers**: 
  - Webhook Service
  - Audit Service
  - Garbage Collection Service (may trigger GC if tag was last reference)
- **Payload**:
  ```json
  {
    "tag": {
      "id": "uuid",
      "name": "string",
      "artifactId": "uuid",
      "deletedAt": "ISO-8601 timestamp",
      "deletedBy": {
        "id": "uuid",
        "username": "string"
      }
    },
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/registry/tag-deleted.json`
- **Subject**: Tag name
- **Ordering**: Per-repository
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Deleting already-deleted tag: no-op
  - Consumers should be idempotent
- **Versioning**: Same as above

#### 1.5 registry.blob.stored
- **Producer**: Registry Service
- **Consumers**: 
  - Storage Monitoring Service (tracks storage usage)
  - Garbage Collection Service (updates reference counts)
- **Payload**:
  ```json
  {
    "blob": {
      "digest": "sha256:...",
      "mediaType": "string",
      "size": 123456,
      "storedAt": "ISO-8601 timestamp",
      "storedBy": {
        "id": "uuid",
        "username": "string"
      }
    },
    "artifactId": "uuid"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/registry/blob-stored.json`
- **Subject**: Blob digest
- **Ordering**: Per-blob
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Storing same blob twice: no-op (content-addressable storage)
  - Consumers should be idempotent
- **Versioning**: Same as above

#### 1.6 registry.blob.deleted
- **Producer**: Garbage Collection Service
- **Consumers**: 
  - Storage Monitoring Service
  - Audit Service
- **Payload**:
  ```json
  {
    "blob": {
      "digest": "sha256:...",
      "mediaType": "string",
      "size": 123456,
      "deletedAt": "ISO-8601 timestamp"
    },
    "gcRunId": "uuid"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/registry/blob-deleted.json`
- **Subject**: Blob digest
- **Ordering**: Per-blob
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Deleting already-deleted blob: no-op
  - Consumers should be idempotent
- **Versioning**: Same as above

#### 1.7 repository.created
- **Producer**: API Service
- **Consumers**: 
  - Webhook Service
  - Audit Service
  - Search Indexing Service
- **Payload**:
  ```json
  {
    "repository": {
      "id": "uuid",
      "name": "string",
      "namespaceId": "uuid",
      "description": "string",
      "visibility": "public|private|protected",
      "createdAt": "ISO-8601 timestamp",
      "createdBy": {
        "id": "uuid",
        "username": "string"
      }
    },
    "namespace": {
      "id": "uuid",
      "name": "string"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/repository-created.json`
- **Subject**: Repository ID
- **Ordering**: Per-namespace
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Creating same repository twice: second attempt fails (unique constraint)
  - Consumers should handle duplicate creation attempts
- **Versioning**: Same as above

#### 1.8 repository.deleted
- **Producer**: API Service
- **Consumers**: 
  - Webhook Service
  - Audit Service
  - Search Indexing Service (for removal)
- **Payload**:
  ```json
  {
    "repository": {
      "id": "uuid",
      "name": "string",
      "namespaceId": "uuid",
      "deletedAt": "ISO-8601 timestamp",
      "deletedBy": {
        "id": "uuid",
        "username": "string"
      }
    },
    "namespace": {
      "id": "uuid",
      "name": "string"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/repository-deleted.json`
- **Subject**: Repository ID
- **Ordering**: Per-namespace
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Deleting already-deleted repository: no-op
  - Consumers should be idempotent
- **Versioning**: Same as above

### 2. Trust Score Events
Events related to trust score calculation and updates.

#### 2.1 trustscore.calculated
- **Producer**: Trust Score Service
- **Consumers**: 
  - Webhook Service
  - API Service (updates trust score cache)
  - Audit Service (logs for compliance)
- **Payload**:
  ```json
  {
    "trustScore": {
      "id": "uuid",
      "artifactId": "uuid",
      "score": 0.0-1.0,
      "level": "trusted|high|medium|low|untrusted",
      "factors": {
        "vulnerability": { "score": 0.0, "details": { ... } },
        "sbom": { "score": 0.0, "details": { ... } },
        "signature": { "score": 0.0, "details": { ... } },
        "provenance": { "score": 0.0, "details": { ... } },
        "license": { "score": 0.0, "details": { ... } },
        "maintenance": { "score": 0.0, "details": { ... } },
        "policy": { "score": 0.0, "details": { ... } },
        "community": { "score": 0.0, "details": { ... } }
      },
      "policyId": "uuid (optional)",
      "calculatedAt": "ISO-8601 timestamp",
      "version": 1
    },
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-calculated.json`
- **Subject**: Artifact digest
- **Ordering**: Per-artifact
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Recalculating same artifact: updates score (higher version number)
  - Consumers should use highest version for given artifact
- **Versioning**: 
  - Version field in payload
  - Schema version in `dataschema`
  - Backward compatible: add optional fields to factors

#### 2.2 trustscore.updated
- **Producer**: Trust Score Service
- **Consumers**: Same as `trustscore.calculated`
- **Payload**: Identical to `trustscore.calculated`
- **Notes**: 
  - Sent when trust score changes significantly (configurable threshold)
  - Used for efficient updates (avoids sending same score repeatedly)
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-updated.json`
- **Subject**: Artifact digest
- **Ordering**: Per-artifact
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: Same as above
- **Versioning**: Same as above

#### 2.3 trustscore.policy.evaluate
- **Producer**: Trust Score Service or API Service
- **Consumers**: Policy Engine Service (OPA)
- **Payload**:
  ```json
  {
    "evaluationId": "uuid",
    "artifactId": "uuid",
    "policyId": "uuid",
    "triggeredBy": "push|scan|manual|scheduled",
    "triggeredAt": "ISO-8601 timestamp",
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    },
    "context": {
      "trustScore": 0.0-1.0,
      "vulnerabilityCount": { "critical": 0, "high": 0, "medium": 0, "low": 0 },
      "hasSbom": true,
      "hasSignature": false,
      "sbomFormat": "SPDX",
      "licenseCompliance": 0.0-1.0
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-policy-evaluate.json`
- **Subject**: Policy ID
- **Ordering**: Per-policy
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each evaluation has unique ID
  - Safe to process multiple times (same result)
- **Versioning**: Same as above

#### 2.4 trustscore.policy.result
- **Producer**: Policy Engine Service (OPA)
- **Consumers**: 
  - Trust Score Service (updates score with policy factor)
  - Webhook Service (for policy violation alerts)
  - Audit Service (logs policy decisions)
- **Payload**:
  ```json
  {
    "evaluationId": "uuid",
    "policyId": "uuid",
    "artifactId": "uuid",
    "result": "pass|fail|warn",
    "reason": "string (optional)",
    "evaluatedAt": "ISO-8601 timestamp",
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-policy-result.json`
- **Subject**: Policy ID
- **Ordering**: Per-policy
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each result has unique evaluation ID
  - Safe to process multiple times
- **Versioning**: Same as above

#### 2.5 trustscore.sbom.generated
- **Producer**: Trust Score Service
- **Consumers**: 
  - Webhook Service
  - Audit Service
  - SBOM Storage Service
- **Payload**:
  ```json
  {
    "sbom": {
      "id": "uuid",
      "artifactId": "uuid",
      "format": "SPDX|CycloneDX",
      "generatedAt": "ISO-8601 timestamp",
      "generatedBy": "syft",
      "size": 12345
    },
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-sbom-generated.json`
- **Subject**: Artifact digest
- **Ordering**: Per-artifact
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Generating SBOM for same artifact: updates (higher timestamp)
  - Consumers should use latest by timestamp
- **Versioning**: Same as above

#### 2.6 trustscore.vulnerability.found
- **Producer**: Trust Score Service
- **Consumers**: 
  - Webhook Service
  - Audit Service
  - Vulnerability Tracking Service
- **Payload**:
  ```json
  {
    "vulnerability": {
      "id": "uuid",
      "artifactId": "uuid",
      "scanner": "trivy",
      "vulnerabilityId": "CVE-2023-12345",
      "severity": "critical|high|medium|low|unknown",
      "title": "string",
      "description": "string",
      "references": ["string"],
      "fixedVersion": "string (optional)",
      "discoveredAt": "ISO-8601 timestamp"
    },
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-vulnerability-found.json`
- **Subject**: Artifact digest
- **Ordering**: Per-artifact
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Discovering same vulnerability: updates (if new info)
  - Each vulnerability instance has unique ID
- **Versioning**: Same as above

#### 2.7 trustscore.signature.verified
- **Producer**: Trust Score Service
- **Consumers**: 
  - Webhook Service
  - Audit Service
- **Payload**:
  ```json
  {
    "signature": {
      "id": "uuid",
      "artifactId": "uuid",
      "type": "cosign",
      "keyId": "string",
      "verifiedAt": "ISO-8601 timestamp",
      "verificationStatus": "verified"
    },
    "artifact": {
      "digest": "sha256:...",
      "repositoryId": "uuid"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-signature-verified.json`
- **Subject**: Artifact digest
- **Ordering**: Per-artifact
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Verifying same signature: no-op if already verified
  - Each verification attempt has unique ID
- **Versioning**: Same as above

#### 2.8 trustscore.signature.failed
- **Producer**: Trust Score Service
- **Consumers**: Same as `trustscore.signature.verified`
- **Payload**: Identical but with `verificationStatus: "failed"`
- **Schema**: `https://schemas.kyros.example.com/v1/events/trustscore-signature-failed.json`
- **Subject**: Artifact digest
- **Ordering**: Per-artifact
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: Same as above
- **Versioning**: Same as above

### 3. Webhook Events
Events related to webhook management and delivery.

#### 3.1 webhook.created
- **Producer**: API Service
- **Consumers**: 
  - Webhook Service (loads new subscription)
  - Audit Service
- **Payload**:
  ```json
  {
    "webhook": {
      "id": "uuid",
      "name": "string",
      "url": "string",
      "events": ["string"],
      "secretHash": "string (bcrypt)",
      "format": "JSON|form-urlencoded",
      "headers": { "string": "string" },
      "enabled": true,
      "createdAt": "ISO-8601 timestamp",
      "createdBy": {
        "id": "uuid",
        "username": "string"
      }
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/webhook-created.json`
- **Subject**: Webhook ID
- **Ordering**: Per-webhook
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Creating same webhook name: fails (unique constraint)
  - Consumers should handle duplicate creation
- **Versioning**: Same as above

#### 3.2 webhook.updated
- **Producer**: API Service
- **Consumers**: Same as `webhook.created`
- **Payload**: Identical to `webhook.created` (with updated fields)
- **Schema**: `https://schemas.kyros.example.com/v1/events/webhook-updated.json`
- **Subject**: Webhook ID
- **Ordering**: Per-webhook
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Updating same webhook: last write wins
  - Consumers should use latest by timestamp
- **Versioning**: Same as above

#### 3.3 webhook.deleted
- **Producer**: API Service
- **Consumers**: Same as `webhook.created`
- **Payload**:
  ```json
  {
    "webhook": {
      "id": "uuid",
      "name": "string",
      "deletedAt": "ISO-8601 timestamp",
      "deletedBy": {
        "id": "uuid",
        "username": "string"
      }
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/webhook-deleted.json`
- **Subject**: Webhook ID
- **Ordering**: Per-webhook
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Deleting already-deleted webhook: no-op
  - Consumers should be idempotent
- **Versioning**: Same as above

#### 3.4 webhook.triggered
- **Producer**: Webhook Service
- **Consumers**: 
  - Audit Service (logs delivery attempts)
  - Webhook Service (tracks triggering)
- **Payload**:
  ```json
  {
    "webhook": {
      "id": "uuid",
      "name": "string",
      "url": "string"
    },
    "event": {
      "id": "string (NATS JetStream msg ID)",
      "type": "string (original event type)",
      "time": "ISO-8601 timestamp"
    },
    "triggeredAt": "ISO-8601 timestamp"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/webhook-triggered.json`
- **Subject**: Webhook ID
- **Ordering**: Per-webhook
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each trigger has unique event ID
  - Safe to process multiple times
- **Versioning**: Same as above

#### 3.5 webhook.delivery.success
- **Producer**: Webhook Service
- **Consumers**: 
  - Audit Service
  - Webhook Service (resets failure count)
- **Payload**:
  ```json
  {
    "webhook": {
      "id": "uuid",
      "name": "string"
    },
    "delivery": {
      "id": "uuid",
      "attempt": 1,
      "status": "success",
      "responseCode": 200,
      "responseBody": "string (truncated)",
      "deliveredAt": "ISO-8601 timestamp"
    },
    "event": {
      "id": "string",
      "type": "string"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/webhook-delivery-success.json`
- **Subject**: Webhook ID
- **Ordering**: Per-webhook
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each delivery attempt has unique ID
  - Success is idempotent (safe to process multiple times)
- **Versioning**: Same as above

#### 3.6 webhook.delivery.failed
- **Producer**: Webhook Service
- **Consumers**: 
  - Audit Service
  - Webhook Service (increments failure count, schedules retry)
- **Payload**:
  ```json
  {
    "webhook": {
      "id": "uuid",
      "name": "string"
    },
    "delivery": {
      "id": "uuid",
      "attempt": 3,
      "status": "failed",
      "responseCode": 500,
      "responseBody": "string (truncated)",
      "failedAt": "ISO-8601 timestamp"
    },
    "event": {
      "id": "string",
      "type": "string"
    },
    "error": "string"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/webhook-delivery-failed.json`
- **Subject**: Webhook ID
- **Ordering**: Per-webhook
- **Retry Strategy**: 
  - This event is sent after max retries exhausted
  - No further retries attempted
- **Dead Letter Handling**: 
  - This event itself is not retried
  - The original event is in DLQ
- **Idempotency**: 
  - Each failed delivery attempt has unique ID
  - Safe to process multiple times
- **Versioning**: Same as above

#### 3.7 webhook.delivery.retry
- **Producer**: Webhook Service
- **Consumers**: 
  - Audit Service
  - Webhook Service (tracks retry attempts)
- **Payload**:
  ```json
  {
    "webhook": {
      "id": "uuid",
      "name": "string"
    },
    "delivery": {
      "id": "uuid",
      "attempt": 2,
      "status": "pending",
      "retryAt": "ISO-8601 timestamp"
    },
    "event": {
      "id": "string",
      "type": "string"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/webhook-delivery-retry.json`
- **Subject**: Webhook ID
- **Ordering**: Per-webhook
- **Retry Strategy**: 
  - This event indicates a retry is scheduled
  - Actual retry handled by Webhook Service internally
- **Dead Letter Handling**: Not applicable (still in retry cycle)
- **Idempotency**: 
  - Each retry attempt has unique ID
  - Safe to process multiple times
- **Versioning**: Same as above

### 4. Authentication Events
Events related to user authentication and authorization.

#### 4.1 auth.user.login
- **Producer**: Auth Service
- **Consumers**: 
  - Audit Service
  - Analytics Service (login metrics)
  - Webhook Service (for security alerts on suspicious logins)
- **Payload**:
  ```json
  {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string"
    },
    "loginAt": "ISO-8601 timestamp",
    "ipAddress": "string",
    "userAgent": "string",
    "mfaUsed": true,
    "authMethod": "password|saml|oidc"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/auth-user-login.json`
- **Subject**: User ID
- **Ordering**: Per-user
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each login attempt has unique timestamp
  - Multiple logins valid (concurrent sessions)
- **Versioning**: Same as above

#### 4.2 auth.user.logout
- **Producer**: Auth Service
- **Consumers**: 
  - Audit Service
  - Webhook Service
- **Payload**:
  ```json
  {
    "user": {
      "id": "uuid",
      "username": "string"
    },
    "logoutAt": "ISO-8601 timestamp",
    "sessionId": "string"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/auth-user-logout.json`
- **Subject**: User ID
- **Ordering**: Per-user
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each logout has unique timestamp
  - Safe to process multiple times
- **Versioning**: Same as above

#### 4.3 auth.user.failed_login
- **Producer**: Auth Service
- **Consumers**: 
  - Audit Service
  - Security Monitoring Service (brute force detection)
  - Webhook Service (for account lock alerts)
- **Payload**:
  ```json
  {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string"
    },
    "failedAt": "ISO-8601 timestamp",
    "ipAddress": "string",
    "userAgent": "string",
    "failureReason": "invalid_credentials|mfa_required|account_locked"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/auth-user-failed-login.json`
- **Subject**: User ID
- **Ordering**: Per-user
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each failed attempt has unique timestamp
  - Multiple failures valid (separate attempts)
- **Versioning**: Same as above

#### 4.4 auth.user.created
- **Producer**: Auth Service
- **Consumers**: 
  - Audit Service
  - Welcome Email Service
  - Webhook Service
- **Payload**:
  ```json
  {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "displayName": "string",
      "createdAt": "ISO-8601 timestamp",
      "createdBy": {
        "id": "uuid",
        "username": "string"
      },
      "enabled": true,
      "emailVerified": false
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/auth-user-created.json`
- **Subject**: User ID
- **Ordering**: Per-user
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Creating same user twice: fails (unique username/email)
  - Consumers should handle duplicate creation
- **Versioning**: Same as above

#### 4.5 auth.token.issued
- **Producer**: Auth Service
- **Consumers**: 
  - Audit Service
  - Token Monitoring Service
- **Payload**:
  ```json
  {
    "tokenId": "uuid",
    "userId": "uuid",
    "clientId": "uuid (optional)",
    "issuedAt": "ISO-8601 timestamp",
    "expiresAt": "ISO-8601 timestamp",
    "scope": ["string"],
    "tokenType": "access|refresh"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/auth-token-issued.json`
- **Subject**: Token ID
- **Ordering**: Per-token
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each token issuance has unique ID
  - Safe to process multiple times
- **Versioning**: Same as above

#### 4.6 auth.token.revoked
- **Producer**: Auth Service
- **Consumers**: 
  - Audit Service
  - Token Monitoring Service
- **Payload**:
  ```json
  {
    "tokenId": "uuid",
    "userId": "uuid",
    "revokedAt": "ISO-8601 timestamp",
    "reason": "user_logout|admin_revoke|token_expired"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/auth-token-revoked.json`
- **Subject**: Token ID
- **Ordering**: Per-token
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each token revocation has unique timestamp
  - Revoking already-revoked token: no-op
- **Versioning**: Same as above

### 5. Audit Events
Events related to audit logging and compliance.

#### 5.1 audit.event.logged
- **Producer**: Audit Service
- **Consumers**: 
  - Audit Storage Service (persists to immutable storage)
  - Audit Alerting Service (triggers alerts on sensitive events)
- **Payload**:
  ```json
  {
    "auditEvent": {
      "id": "uuid",
      "timestamp": "ISO-8601 timestamp",
      "actorId": "uuid (optional)",
      "actorType": "user|service-account|system",
      "action": "string",
      "resourceType": "string",
      "resourceId": "uuid (optional)",
      "resourceName": "string (optional)",
      "outcome": "success|failure",
      "failureReason": "string (optional)",
      "ipAddress": "string",
      "userAgent": "string",
      "requestId": "string",
      "changes": {
        "before": { ... },
        "after": { ... }
      },
      "tenantId": "uuid (optional)"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/audit-event-logged.json`
- **Subject**: Audit event ID
- **Ordering**: Per-audit-event (sequential by timestamp)
- **Retry Strategy**: 
  - Audit service attempts to persist to storage
  - On failure, buffers in memory and retries
  - After multiple failures, alerts operator
- **Dead Letter Handling**: 
  - Events that cannot be persisted after max retries are logged to local disk
  - Operator must manually recover
- **Idempotency**: 
  - Each audit event has unique ID
  - Duplicate events should be deduplicated by ID
- **Versioning**: Same as above

#### 5.2 audit.retention.expired
- **Producer**: Audit Service
- **Consumers**: 
  - Audit Storage Service (for cleanup)
  - Compliance Reporting Service (for retention metrics)
- **Payload**:
  ```json
  {
    "auditEventId": "uuid",
    "expiredAt": "ISO-8601 timestamp",
    "retentionPolicyId": "uuid"
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/audit-retention-expired.json`
- **Subject**: Audit event ID
- **Ordering**: Per-audit-event
- **Retry Strategy**: Not applicable (event is about expiration)
- **Dead Letter Handling**: Not applicable
- **Idempotency**: 
  - Each expiration event has unique audit event ID
  - Safe to process multiple times
- **Versioning**: Same as above

### 6. System Events
Events related to system operation and health.

#### 6.1 service.started
- **Producer**: Any Kyros Service
- **Consumers**: 
  - Monitoring Service (tracks service uptime)
  - Alerting Service (triggers startup alerts)
  - Dashboard Service (updates service status)
- **Payload**:
  ```json
  {
    "service": {
      "name": "string",
      "instanceId": "string",
      "version": "string",
      "startedAt": "ISO-8601 timestamp",
      "host": "string",
      "port": "number"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/service-started.json`
- **Subject**: Service instance ID
- **Ordering**: Per-service-instance
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Starting same service instance: should not happen (restart creates new instance ID)
  - Multiple start events for same instance indicate restart
- **Versioning**: Same as above

#### 6.2 service.stopped
- **Producer**: Any Kyros Service
- **Consumers**: Same as `service.started`
- **Payload**:
  ```json
  {
    "service": {
      "name": "string",
      "instanceId": "string",
      "stoppedAt": "ISO-8601 timestamp",
      "host": "string",
      "port": "number"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/service-stopped.json`
- **Subject**: Service instance ID
- **Ordering**: Per-service-instance
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Stopping same service instance: should not happen twice without start in between
  - Multiple stops indicate repeated stopping
- **Versioning**: Same as above

#### 6.3 service.failed
- **Producer**: Any Kyros Service
- **Consumers**: 
  - Monitoring Service (triggers failure alerts)
  - Incident Response Service (creates incident)
  - Dashboard Service (updates service status to failed)
- **Payload**:
  ```json
  {
    "service": {
      "name": "string",
      "instanceId": "string",
      "failedAt": "ISO-8601 timestamp",
      "host": "string",
      "port": "number"
    },
    "error": {
      "type": "string",
      "message": "string",
      "stackTrace": "string (optional)"
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/service-failed.json`
- **Subject**: Service instance ID
- **Ordering**: Per-service-instance
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Each failure event has unique timestamp
  - Multiple failures for same instance valid
- **Versioning**: Same as above

#### 6.4 config.updated
- **Producer**: Any Kyros Service
- **Consumers**: 
  - Configuration Service (reloads configuration)
  - Audit Service (logs configuration changes)
  - Affected Services (may need to reconfigure)
- **Payload**:
  ```json
  {
    "config": {
      "service": "string",
      "key": "string",
      "oldValue": "string (optional)",
      "newValue": "string",
      "updatedAt": "ISO-8601 timestamp",
      "updatedBy": {
        "id": "uuid",
        "username": "string"
      }
    }
  }
  ```
- **Schema**: `https://schemas.kyros.example.com/v1/events/config-updated.json`
- **Subject**: Config key (service.key)
- **Ordering**: Per-config-key
- **Retry Strategy**: Same as above
- **Dead Letter Handling**: Same as above
- **Idempotency**: 
  - Updating same config key: last write wins
  - Consumers should use latest by timestamp
- **Versioning**: Same as above

## Event Ordering Guarantees

### Per-Key Ordering
- Events with the same subject key are processed in order by consumers
- Example: All events for artifact `sha256:abc...` are processed in the order they were produced
- Different keys (different artifacts) have no ordering guarantee

### Consumer Group Ordering
- Within a consumer group, events are distributed by key for load balancing
- Each consumer in the group processes a subset of keys
- For a given key, all events go to the same consumer in the group
- This ensures per-key ordering while allowing horizontal scaling

### Global Ordering
- No global ordering guarantee across different keys
- Events for different artifacts may be processed out of order relative to each other
- Applications should not rely on cross-event ordering

## Retry Strategy Details

### Initial Attempt
- First delivery attempt happens immediately after event is stored

### Backoff Algorithm
- Attempt 1: Immediate
- Attempt 2: 1 second + random jitter (0-500ms)
- Attempt 3: 2 seconds + random jitter (0-500ms)
- Attempt 4: 4 seconds + random jitter (0-500ms)
- Attempt 5: 8 seconds + random jitter (0-500ms)
- Maximum 5 attempts (1 initial + 4 retries)

### Jitter
- Random delay added to prevent thundering herd problems
- Uniform distribution between 0 and 500ms

### Exponential Backoff
- Delay doubles with each attempt (after initial)
- Cap at maximum delay (30 seconds) to prevent excessive waits

### Final Attempt
- After max retries, event is sent to dead letter queue
- No further delivery attempts attempted

## Dead Letter Queue Handling

### DLQ Structure
- Each stream has a corresponding dead letter queue stream
- Named `dlq.{original-stream-name}`
- Example: `dlq.registry.events` for `registry.events` stream

### DLQ Retention
- Retained for 14 days (configurable)
- After retention period, events are permanently deleted

### DLQ Processing
- Manual process: Operators inspect DLQ events
- Automated replay: After fixing underlying issue, events can be replayed
- Selective replay: Ability to replay specific events or time ranges
- Transformation: Events can be modified before replay (e.g., update credentials)

### Monitoring
- Alert on DLQ depth > 0
- Alert on DLQ growth rate > X events/hour
- Dashboard showing DLQ contents and age

## Idempotency Considerations

### Why Idempotency Matters
- At-least-once delivery means events may be delivered more than once
- Network retries, consumer crashes during processing, etc. can cause duplicates
- Non-idempotent processing can lead to inconsistent state

### Idempotency Patterns

#### 1. Event ID Tracking
- Store processed event IDs in a durable store (Redis, database)
- Before processing, check if event ID already processed
- Skip if already processed
- TTL on stored IDs to prevent unbounded growth (e.g., 24 hours)

#### 2. State-Based Processing
- Check current state before performing action
- Example: Before creating repository, check if it already exists
- Only perform action if state transition is valid

#### 3. Version Vectors/Timestamps
- For updatable resources, use version numbers or timestamps
- Only process event if its version is newer than current version
- Discard or merge older events

#### 4. External System Idempotency
- Leverage idempotency of external systems
- Example: Use PUT instead of POST for creating resources when possible
- Use unique idempotency keys for payment processing

### Implementation Example
```go
func (s *TrustScoreService) handleArtifactPushed(event *events.Event) error {
    // Check if we've already processed this event
    if s.eventProcessor.IsProcessed(event.ID) {
        s.logger.Debug("Skipping duplicate event",
                      zap.String("event_id", event.ID),
                      zap.String("event_type", event.Type))
        return nil
    }
    
    // Process the event
    digest := extractDigest(event.Data)
    if err := s.calculateTrustScore(digest, event); err != nil {
        s.logger.Error("Failed to calculate trust score",
                      zap.Error(err),
                      zap.String("event_id", event.ID))
        return nil // Don't mark as processed on failure to allow retry
    }
    
    // Mark event as processed
    if err := s.eventProcessor.MarkProcessed(event.ID); err != nil {
        s.logger.Error("Failed to mark event as processed",
                      zap.Error(err),
                      zap.String("event_id", event.ID))
        // Don't return error - event was processed successfully
    }
    
    return nil
}
```

## Event Versioning

### Backward Compatible Changes
- Add new optional fields to event payload
- Add new enum values
- Make required fields optional (with sensible defaults)
- Add new event types (doesn't break existing consumers)

### Forward Compatible Changes
- Remove optional fields (if consumers ignore unknown fields)
- Make optional fields required (if producers already provide them)
- Change field names (if consumers use aliasing or mapping)

### Breaking Changes
- Change event type semantics (e.g., `artifact.pushed` now means something different)
- Remove required fields
- Change data types in incompatible ways (string to integer)
- Change subject key format

### Versioning Strategy
1. **Schema Version**: Included in `dataschema` field (e.g., `v1`, `v2`)
2. **Event Type Version**: Included in event type for major changes (e.g., `registry.v2.artifact.pushed`)
3. **Payload Version**: Version field in payload for internal versioning
4. **Consumer Adaptation**: 
   - Consumers declare which schema versions they support
   - Unsupported versions sent to DLQ or logged
   - Version negotiation during startup

### Deprecation Policy
- Deprecated event types supported for 6 months after deprecation notice
- Clear migration guide provided
- Automated tooling to help migrate consumers
- After deprecation period, unsupported events rejected at source

## Security Considerations

### Transport Security
- All NATS connections encrypted with TLS 1.2+
- Mutual TLS authentication for service-to-service communication
- JWT or NKEY authentication for service identification
- Short-lived credentials with automatic rotation

### Authentication and Authorization
- Services authenticate to NATS using service accounts
- JWT tokens or NKEY pairs for identification
- Stream/subject-based permissions:
  - Publish/subscribe permissions granted separately
  - Wildcard subject matching for efficient permission definition
- Claims-based authorization using JWT claims
- Tenant and user information extracted from tokens
- Event producers can only publish events for their tenant/user

### Event Integrity
- Optional schema validation against registered schemas
- Optional cryptographic signing of events for high-security domains
- Optional hash chaining for tamper evidence
- Timestamp rejection for events with skewed clocks ( >5min skew)

### Audit and Compliance
- Separate immutable audit stream for audit events (`audit.events`)
- Access logging: who published/consumed what events
- Configurable retention with legal hold capabilities
- Export capabilities for external audit

## Performance and Scaling

### Publish Latency
- Local: <1ms
- Cross-AZ: <10ms
- 99th percentile: <50ms

### Consume Latency
- Local: <1ms
- Cross-AZ: <10ms
- 99th percentile: <50ms

### Throughput
- 100,000+ msg/sec per cluster node (dependent on hardware)
- Horizontal scaling of producers and consumers
- Batching for improved throughput

### Resource Utilization
#### Memory
- Connection buffers: minimal per-connection overhead
- In-memory streams: configurable memory allocation
- Consumer state: minimal (mainly acknowledgment tracking)
- Batch buffers: configurable batch size

#### Disk
- File-based streams: disk space proportional to retained event size
- Minimal indexing overhead for subject-based lookup
- Automatic compaction of deleted events
- Segment-based rotation for efficient disk usage

### Monitoring Metrics
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

## Implementation Best Practices

### For Event Producers
1. **Publish After Commit**: Only publish events after business transaction is committed
2. **Include Context**: Always include trace IDs, tenant IDs, and user IDs when available
3. **Validate Before Publishing**: Validate business data before creating events
4. **Handle Publishing Failures**: Log errors but don't fail business transactions for event publishing
5. **Use Appropriate Subjects**: Use hierarchical subjects for effective filtering (e.g., `registry.artifact.pushed`)
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
This event catalog provides a comprehensive reference for all events in the Kyros platform. By adhering to the specifications outlined here, developers can ensure consistent, reliable, and secure event-driven communication between services. The event-driven architecture enables loose coupling, scalability, and resilience while maintaining clear semantics and guarantees for event processing.

Regular review and updates to this catalog are recommended as the platform evolves. New events should be added following the same patterns and guidelines established here.