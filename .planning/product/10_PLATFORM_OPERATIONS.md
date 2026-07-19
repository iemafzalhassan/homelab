# Kyros Platform Operations

## Overview
This document defines the operational aspects of running the Kyros platform in production, including Service Level Objectives (SLOs), Service Level Indicators (SLIs), error budgets, capacity planning, backup and disaster recovery, upgrade strategies, incident response, and runbooks. It serves as a guide for platform administrators and site reliability engineers (SREs) to ensure the platform is reliable, scalable, and maintainable.

## Service Level Objectives (SLOs) and Service Level Indicators (SLIs)

### Availability SLO
- **SLO**: The Kyros platform (core services: API, Registry, Auth) shall be available 99.9% of the time over a rolling 30-day window.
- **SLI**: 
  - **Good Request**: An HTTP request to a core service endpoint that returns a successful status code (2xx) or a client error (4xx) that is not due to server-side issues (e.g., 401 Unauthorized, 403 Forbidden, 404 Not Found, 429 Too Many Requests are considered good if they are expected responses).
  - **Bad Request**: An HTTP request that returns a server error (5xx) or a timeout.
  - **Availability** = (Good Requests) / (Total Requests) * 100%
- **Measurement**: 
  - Requests are measured at the ingress (e.g., Traefik or cloud load balancer).
  - Exclude scheduled maintenance windows (with prior notice).

### Latency SLO
- **SLO**: 95% of requests to core services shall have a latency of less than 200ms over a rolling 30-day window.
- **SLI**:
  - **Latency**: The time from when the request is received by the ingress to when the response is fully sent.
  - **Good Request**: Request with latency <= 200ms.
  - **Total Requests**: All requests to core services (excluding health check endpoints).
- **Measurement**:
  - Percentile calculated over a sliding window.
  - Exclude requests during known degradation periods (if any).

### Durability SLO
- **SLO**: Stored objects (blobs, metadata) in Kyros shall have an annual durability of 99.999999999% (11 nines).
- **SLI**:
  - **Durability Event**: Loss of a stored object due to storage medium failure.
  - **Durability** = 1 - (Number of lost objects / Total object-months) * (12 months / observation period in months)
- **Measurement**:
  - Based on storage backend guarantees (e.g., Amazon S3 offers 11 nines).
  - Kyros uses multiple storage backends and erasure coding where applicable to achieve this target.
  - Regular durability audits are performed.

### Freshness SLO (Trust Score SLO: The trust score for an artifact shall be calculated and 5 minutes of the artifact being pushed to the registry.
- SLI:
  - Freshness Event: An artifact push event that results in a trust score being calculated and stored within 5 minutes.
  - Stale Event: An artifact push event that does not have a trust score calculated and stored within 5 minutes.
  - Freshness = (Number of Freshness Events) / (Total Artifact Push Events) * 100%
- Measurement:
  - Measured by the Trust Score Service by tracking the time between the `registry.artifact.pushed` event and the `trustscore.calculated` event for each artifact.

### Correctness SLO)
- SLO: The trust score SLO
- SLO: 99% of trust score calculations shall complete within 2 minutes.
- SLI:
  - Calculation Time: The time from when the Trust Score Service starts processing an artifact to when it stores the trust score.
  - Good Calculation: Calculation time <= 120 seconds.
  - Total Calculations: All trust score calculations attempted.
- Measurement:
  - Instrumented in the Trust Score Service with histograms.

## Error Budgets

### Definition
The error budget is the allowable amount of slack (unreliability) within an SLO. It is calculated as:
- Error Budget = 100% - SLO Target

### Example
For the Availability SLO of 99.9%:
- Error Budget = 0.1% over a 30-day window.
- In a 30-day month (43,200 minutes), this allows for 43.2 minutes of downtime or errors.

### Error Budget Policy
1. **Green Status**: Error budget consumption < 50% of monthly budget -> No restrictions on releases or changes.
2. **Yellow Status**: Error budget consumption between 50% and 90% -> 
   - Non-essential changes require approval.
   - Focus on reliability work (e.g., reducing technical debt, improving monitoring).
3. **Red Status**: Error budget consumption > 90% ->
   - Freeze all non-essential changes.
   - Prioritize incident response and postmortems.
   - Only emergency fixes and critical security patches are allowed.

### Tracking
- Error budget is tracked per SLO and aggregated for overall service health.
- Alerts trigger when error budget burn rate exceeds thresholds (e.g., burning through the monthly budget in less than 5 days).

## Capacity Planning

### Metrics to Monitor
- **Request Rate**: Requests per second (RPS) to core services.
- **Resource Utilization**: CPU, memory, disk I/O, network bandwidth per service instance.
- **Queue Depth**: For event processing (NATS JetStream consumers) and background jobs.
- **Storage Usage**: Blob storage and database growth rates.
- **Latency**: As defined in SLIs.
- **Error Rates**: As defined in SLIs.

### Forecasting
- Use historical data (minimum 90 days) to forecast growth.
- Apply growth rates (monthly, quarterly) to predict future resource needs.
- Account for seasonal variations and planned events (e.g., major releases, marketing campaigns).

### Scaling Triggers
- **Horizontal Pod Autoscaler (HPA)**: 
  - Scale based on CPU utilization (target 60%) or custom metrics (request rate, queue length).
- **Vertical Pod Autoscaler (VPA)**: 
  - Adjust resource requests and limits based on observed usage.
- **Cluster Autoscaler**: 
  - Adjust the number of nodes in the Kubernetes cluster based on pod scheduling needs.

### Capacity Review
- Conduct capacity planning reviews quarterly.
- Simulate peak load scenarios (e.g., Black Friday, product launches) to ensure adequate headroom.
- Maintain a minimum of 50% headroom for unexpected traffic spikes.

## Backup and Disaster Recovery

### Backup Strategy
#### 1. Metadata Database (PostgreSQL)
- **Logical Backups**: 
  - Daily `pg_dump` of all databases.
  - Retained for 30 days.
- **Physical Backups**:
  - Continuous archiving of WAL files using `wal-e` or `wal-g`.
  - Base backups taken weekly.
  - Point-in-time recovery (PITR) available for any point in the last 30 days.
- **Retention**:
  - Daily backups: 30 days
  - Weekly backups: 12 weeks
  - Monthly backups: 12 months
  - Yearly backups: 5 years (for compliance)

#### 2. Blob Storage (MinIO/S3)
- **Versioning**: Enable object versioning to protect against accidental deletion.
- **Replication**: 
  - Cross-Region Replication (CRR) for disaster recovery.
  - Minimum 2 copies in geographically separate regions.
- **Lifecycle Rules**:
  - Transition infrequently accessed objects to cheaper storage after 90 days.
  - Expire deleted objects (delete markers) after 365 days.
- **Backup**: 
  - Not typically needed due to replication and versioning, but periodic snapshots can be taken for regulatory compliance.

#### 3. Configuration and Secrets
- **GitOps**: All configuration stored in a Git repository (the source of truth).
- **Secrets**: 
  - Stored in a secrets manager (e.g., HashiCorp Vault, Azure Key Vault, AWS Secrets Manager).
  - Backed up by the secrets manager's native backup mechanisms.
- **Manual Backups**: 
  - Export of non-secrets configuration (e.g., ConfigMaps) as part of GitOps.

#### 4. Event Stream (NATS JetStream)
- **Persistence**: 
  - JetStream stores events on disk by default.
  - Ensure the underlying disk is backed up (if using persistent volumes).
- **Replication**: 
  - JetStream clusters can be run in replicated mode (requires NATS clustering).
  - For high availability, run a JetStream cluster with at least 3 nodes.
- **Backup**: 
  - Not typically backed up as events are ephemeral; rely on replication for durability.
  - For long-term retention of events, use the audit stream with appropriate retention policies.

### Disaster Recovery Plan
#### 1. Recovery Point Objective (RPO)
- **Target RPO**: 15 minutes for user data (metadata, blob storage).
- **Achieved By**:
  - Database: WAL archiving allows recovery to any point in time.
  - Blob Storage: Replication lag is typically under 15 minutes.
  - Configuration: GitOps ensures immediate availability of configuration.

#### 2. Recovery Time Objective (RTO)
- **Target RTO**: 4 hours for full platform recovery.
- **Achieved By**:
  - Automated failover to secondary region.
  - Pre-configured infrastructure in standby mode.
  - Runbooks for manual intervention if needed.

#### 3. Failover Procedure
1. **Detection**: 
   - Monitoring detects loss of primary region (e.g., via health checks, latency spikes).
   - Manual confirmation of outage.
2. **Decision**: 
   - Incident Commander declares a disaster and initiates failover.
3. **DNS Failover**: 
   - Update DNS records to point to secondary region endpoints (TTL set to 60 seconds for fast failover).
4. **Infrastructure Activation**: 
   - Ensure secondary region has sufficient capacity (pre-warmed or auto-scaling).
5. **Data Synchronization**: 
   - Verify that data replication is up to date (check replication lag).
   - If necessary, initiate a final sync.
6. **Service Startup**: 
   - Start services in secondary region in order of dependencies (e.g., database first, then auth, then registry, then API).
7. **Validation**: 
   - Perform smoke tests (login, push/pull image, trust score calculation).
   - Monitor metrics and logs for anomalies.
8. **Notification**: 
   - Inform stakeholders of failover completion.
4. **Failback Procedure** (when primary region is restored):
   - Similar to failover but in reverse.
   - Ensure data is synchronized from secondary to primary before switching back.

### Regular Drills
- Conduct disaster recovery drills semi-annually.
- Test both failover and failback procedures.
- Update runbooks based on lessons learned.

## Upgrade Strategy

### Versioning
Kyros uses semantic versioning: MAJOR.MINOR.PATCH
- **MAJOR**: Incompatible API changes.
- **MINOR**: Backward-compatible functionality.
- **PATCH**: Backward-compatible bug fixes.

### Upgrade Types
1. **Patch Upgrade**: 
   - In-place rolling update.
   - No downtime expected.
   - Can be automated.
2. **Minor Upgrade**: 
   - May require database schema migrations.
   - Rolling update with compatibility checks.
   - Short downtime possible during schema migration (aim for < 5 minutes).
3. **Major Upgrade**: 
   - May require significant changes (data migration, breaking API changes).
   - Planned maintenance window required.
   - Blue/green or canary release strategy recommended.

### Upgrade Process
#### Pre-Upgrade
1. **Read Release Notes**: 
   - Review breaking changes, deprecations, and migration steps.
2. **Backup**: 
   - Take a full backup of the current state (database, configuration, etc.).
3. **Test in Staging**: 
   - Deploy the new version to a staging environment that mirrors production.
   - Run automated and manual tests.
4. **Notify Stakeholders**: 
   - Inform users of upcoming maintenance (if applicable).
5. **Prepare Rollback Plan**: 
   - Ensure ability to revert to previous version.

#### During Upgrade (for minor/patch)
1. **Update Control Plane**: 
   - Update the Kyros Operator (if used) to the new version.
2. **Update Custom Resources**: 
   - Update any Kyros-specific CRDs if required.
3. **Rolling Update of Services**: 
   - Update services one by one or in batches (e.g., update auth first, then registry, then API).
   - Use Kubernetes rolling update with maxSurge and maxUnavailable set to maintain availability.
4. **Database Migrations**: 
   - Run migration scripts (if any) as a one-time job.
   - Ensure migrations are backward compatible until all instances are upgraded.
5. **Verify**: 
   - Run health checks and smoke tests.

#### Post-Upgrade
1. **Monitor**: 
   - Watch metrics, logs, and error rates for anomalies.
2. **Validate**: 
   - Confirm that new features work as expected.
3. **Clean Up**: 
   - Remove temporary resources used during the upgrade.
4. **Report**: 
   - Document the upgrade process and any issues encountered.

### Rollback Procedure
1. **Detect Issue**: 
   - Identify that the upgrade has caused a problem.
2. **Stop Rollout**: 
   - Halt the rolling update.
3. **Restore Backup**: 
   - If database migration caused issues, restore the database backup.
4. **Revert to Previous Version**: 
   - Roll back the container images to the previous version.
5. **Validate**: 
   - Confirm that the system is stable.
6. **Investigate**: 
   - Perform a postmortem to understand what went wrong and prevent recurrence.

## Incident Response

### Incident Classification
1. **Severity 1 (Critical)**: 
   - Complete platform outage or severe degradation affecting all users.
   - Examples: 
     - All core services unavailable.
     - Data loss or corruption.
     - Security breach.
   - Response Time: 15 minutes.
2. **Severity 2 (High)**: 
   - Significant degradation affecting a large subset of users or critical functionality.
   - Examples: 
     - Registry service unable to accept pushes.
     - Trust score calculation significantly delayed.
     - Authentication service experiencing high latency.
   - Response Time: 30 minutes.
3. **Severity 3 (Medium)**: 
   - Moderate degradation affecting a small subset of users or non-critical functionality.
   - Examples: 
     - Webhook delivery delays.
     - Minor UI issues.
   - Response Time: 1 hour.
4. **Severity 4 (Low)**: 
   - Minor issues with minimal impact.
   - Examples: 
     - Typo in documentation.
     - Non-critical feature not working.
   - Response Time: 24 hours.

### Incident Response Process
1. **Detection**: 
   - Automated alerts from monitoring systems.
   - User reports via support channels.
2. **Triage**: 
   - First responder (on-call engineer) acknowledges the alert and performs initial assessment.
   - Determine severity and impact.
3. **Diagnosis**: 
   - Gather logs, metrics, and traces.
   - Identify root cause (or working hypothesis).
4. **Mitigation**: 
   - Apply immediate fix to reduce impact (e.g., rollback, scale resources, block malicious traffic).
5. **Resolution**: 
   - Implement permanent fix.
   - Verify that the issue is resolved.
6. **Recovery**: 
   - Return to normal operations.
   - Confirm that all services are healthy and metrics are nominal.
7. **Communication**: 
   - Keep stakeholders informed throughout the process.
   - Provide status updates at regular intervals.
8. **Postmortem**: 
   - Conduct a blameless postmortem within 5 business days.
   - Document timeline, root cause, impact, and action items.
   - Track action items to completion.

### Communication Plan
- **Internal**: 
  - Use a dedicated incident channel (e.g., Slack #incident-response).
  - Update status every 15 minutes for Sev1, 30 minutes for Sev2, hourly for Sev3.
- **External (Users)**:
  - Status page (hosted separately) updated with incident details and ETA.
  - Email notifications for subscribed users.
  - Social media updates if appropriate.
- **Post-Incident**: 
  - Publish a postmortem report on the status page or blog.

## Runbooks

### Runbook: High Latency in Registry Service
**Symptoms**: 
- Increased latency for image push/pull operations.
- Alert: Registry service 95th percentile latency > 200ms for 5 minutes.

**Steps**:
1. **Verify**: 
   - Check the latency metric in the monitoring dashboard.
   - Confirm it's not a false positive (e.g., check multiple instances).
2. **Check Resources**: 
   - Look at CPU, memory, disk I/O, and network usage for the registry service pods.
   - If resource saturation is observed, scale the deployment.
3. **Check Dependencies**: 
   - Verify the storage backend (MinIO/S3) is performing well.
   - Check the database for slow queries or connection issues.
4. **Check Logs**: 
   - Look for errors or warnings in the registry service logs.
   - Common causes: 
     - Garbage collection running (check if GC is active).
     - High number of concurrent requests causing thread pool exhaustion.
5. **Mitigation**: 
   - If GC is running, wait for it to complete or adjust schedule.
   - If thread pool exhaustion, increase the number of replicas or adjust thread pool size.
   - If storage backend is slow, investigate storage performance or failover to secondary storage.
6. **Validate**: 
   - Confirm latency returns to acceptable levels.
7. **Document**: 
   - Record actions taken and outcome.

### Runbook: Authentication Service Outage
**Symptoms**: 
- Users unable to log in.
- Auth service health checks failing.

**Steps**:
1. **Verify**: 
   - Check the auth service health endpoint.
   - Confirm multiple instances are affected.
2. **Check Logs**: 
   - Look for errors in the auth service logs.
   - Common causes: 
     - Database connection failure.
     - Keycloak (if used) unavailability.
     - Certificate expiration.
3. **Check Dependencies**: 
   - Verify the database is reachable and healthy.
   - If using Keycloak, check its status.
4. **Failover (if applicable)**: 
   - If using a managed Keycloak service, check its status page.
   - If self-hosted, consider failing over to a secondary Keycloak instance.
5. **Restart**: 
   - If the issue is transient, restart the auth service pods.
6. **Validate**: 
   - Test login functionality with a test account.
7. **Document**: 
   - Record actions taken and outcome.

### Runbook: Trust Score Calculation Backlog
**Symptoms**: 
- Increasing lag in trust score calculation.
- Alert: Trust score service consumer lag > 1000 events for 5 minutes.

**Steps**:
1. **Verify**: 
   - Check the consumer lag metric for the trust score service.
   - Confirm it's increasing over time.
2. **Check Resources**: 
   - Look at CPU and memory usage for the trust score service pods.
   - If resource constrained, scale the deployment.
3. **Check Dependencies**: 
   - Verify the SBOM generator (Syft) and vulnerability scanners (Trivy/Grype) are functioning.
   - Look for errors in the trust score service logs related to these dependencies.
4. **Check Event Stream**: 
   - Verify that the `registry.artifact.pushed` events are being produced at a normal rate.
   - If the event stream is backed up, investigate the registry service.
5. **Optimize**: 
   - If scanners are slow, consider increasing scanner resources or adjusting scan depth.
   - If the trust score calculation is complex, consider simplifying policies or adding caching.
6. **Validate**: 
   - Confirm that the consumer lag decreases and stabilizes.
7. **Document**: 
   - Record actions taken and outcome.

## Conclusion
Effective platform operations are essential for delivering a reliable, scalable, and secure Kyros platform. By defining clear SLOs/SLIs, managing error budgets, planning capacity, implementing robust backup and disaster recovery strategies, following safe upgrade procedures, responding to incidents swiftly, and maintaining detailed runbooks, platform administrators can ensure that Kyros meets the expectations of its users and stakeholders.

Regular review and refinement of these operational practices are necessary to adapt to evolving requirements, technological changes, and lessons learned from operational experience. Continuous improvement in operations directly contributes to the overall success and adoption of the Kyros platform.