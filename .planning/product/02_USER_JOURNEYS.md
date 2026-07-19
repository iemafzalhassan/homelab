# Kyros User Journeys

## Overview
This document outlines the key user journeys for different personas interacting with the Kyros platform. Each journey represents an end-to-end workflow that captures the user's goals, steps, touchpoints, and emotional experience.

## Personas

### 1. Platform Engineer
Responsible for platform operations, security, reliability, and performance. Focuses on keeping the platform running smoothly and securely.

### 2. Developer
Builds and deploys containerized applications. Focuses on rapid development cycles and reliable deployment pipelines.

### 3. Security Engineer
Ensures supply chain security, compliance, and risk management. Focuses on preventing vulnerabilities and enforcing security policies.

### 4. Organization Admin
Manages users, teams, permissions, and organizational settings. Focuses on access control and resource allocation.

### 5. System Admin
Handles infrastructure, deployment, upgrades, and system-level operations. Focuses on availability, performance, and maintenance.

## User Journeys

### 1. Platform Engineer: Ensure Platform Health and Performance

**Goal**: Monitor platform health, identify performance bottlenecks, and ensure optimal operation.

**Steps**:
1. Login to Kyros platform using corporate SSO
2. Navigate to Dashboard view
3. Review system health indicators (all services green)
4. Check key metrics: request latency, error rates, throughput
5. Navigate to Observability section
6. Examine detailed metrics for each service
7. Identify any elevated latency in registry service during peak hours
8. Drill down to specific endpoint performance
9. Correlate with resource utilization (CPU, memory, disk I/O)
10. Identify that garbage collection is causing periodic latency spikes
11. Adjust garbage collection schedule to off-peak hours
12. Configure alerting for sustained high latency
13. Verify changes improve performance metrics
14. Document findings and adjustments in system runbook

**Touchpoints**: Dashboard, Observability, Settings, Alerting
**Emotional Journey**: Confident → Curious → Concerned → Focused → Satisfied
**Success Criteria**: Platform operates within SLA, proactive issue identification

### 2. Developer: Deploy Application Container Image

**Goal**: Build and deploy a new version of their application container image to the registry.

**Steps**:
1. Clone application repository locally
2. Make code changes and test locally
3. Build container image using Dockerfile
4. Tag image with version number: `registry.example.com/my-app:v1.2.3`
5. Authenticate to Kyros registry using SSO credentials
6. Push image to registry: `docker push registry.example.com/my-app:v1.2.3`
7. Monitor push progress in terminal
8. Receive push completion confirmation
9. Navigate to Kyros web UI to verify image availability
10. Check that image appears in repository with correct tag
11. View trust score for newly pushed image (initially "unscored")
12. Wait for automated trust score calculation to complete
13. Review trust score breakdown and security scan results
14. Confirm image meets minimum trust threshold for deployment
15. Update Kubernetes deployment manifest with new image tag
16. Apply changes to trigger rolling update
17. Monitor deployment progress and verify successful rollout
18. Notify team of successful deployment via team channel

**Touchpoints**: Local development environment, CLI, Web UI, Notification system
**Emotional Journey**: Focused → Anticipating → Relieved → Validated → Confident
**Success Criteria**: Image successfully pushed, scanned, approved, and deployed

### 3. Security Engineer: Investigate Potential Vulnerability

**Goal**: Investigate a reported vulnerability in a critical application container.

**Steps**:
1. Receive alert about CVE-2023-12345 affecting OpenSSL versions < 1.1.1t
2. Login to Kyros security dashboard
3. Navigate to Vulnerabilities section
4. Filter by CVE ID: `CVE-2023-12345`
5. View list of affected images across all repositories
6. Sort by severity and exploitability
7. Identify high-priority target: `registry.example.com/payment-service:latest`
8. Click on image to view detailed vulnerability report
9. Examine affected package: OpenSSL 1.1.1k
10. Review available fix: Upgrade to OpenSSL 1.1.1t
11. Check if fix is available in base image
12. Review build configuration (Dockerfile) for image
13. Determine that base image needs updating
14. Create ticket for application team to update base image
15. Apply temporary policy to block promotion of vulnerable image
16. Monitor compliance with new policy
17. Follow up with application team on remediation progress
18. Verify that new image version passes security checks
19. Remove temporary blocking policy
20. Document incident and lessons learned in security knowledge base

**Touchpoints**: Security Dashboard, Vulnerability Details, Policy Management, Notification System
**Emotional Journey**: Alerted → Investigative → Concerned → Proactive → Relieved
**Success Criteria**: Vulnerability identified, contained, and remediation tracked

### 4. Organization Admin: Onboard New Development Team

**Goal**: Set up a new development team with appropriate access to Kyros resources.

**Steps**:
1. Login to Kyros as organization administrator
2. Navigate to Users & Teams section
3. Click "Create New Team"
4. Enter team name: "mobile-app-team"
5. Add team description: "Responsible for mobile backend services"
6. Save team creation
7. Navigate to Users tab
8. Invite new team members via email (import from corporate directory)
9. Assign new users to "mobile-app-team" group
10. Create new namespace for team: "mobile-apps"
10. Set namespace visibility to "private" (team-only access)
11. Set initial resource quotas for namespace:
    - Storage: 50 GB
    - Bandwidth: 5 TB/month
    - Image pull rate: 1000/hour
12. Create team-specific repository: "mobile-apps/auth-service"
13. Configure repository visibility: "private" (team-only)
14. Set up default permissions for team:
    - Team members: push/pull access to team repositories
    - Team leads: additional repository management permissions
    - Security team: read-only audit access
15. Set up automated trust score policy for team:
    - Minimum trust score: 0.7 for production tags
    - Automatic blocking of pushes below threshold
    - Notification to team lead on policy violations
16. Configure webhook for CI/CD integration:
    - URL: https://ci.company.com/webhook/kyros
    - Events: image.pushed, trust.score.updated
    - Secret: [securely stored]
17. Send onboarding documentation to team lead
18. Schedule walkthrough session with new team
19. Monitor initial usage and provide support as needed
20. Quarterly review of team resource usage and access requirements

**Touchpoints**: User & Team Management, Namespace Management, Repository Management, Policy Management, Webhook Configuration
**Emotional Journey**: Organized → Systematic → Thorough → Confident → Supportive
**Success Criteria**: Team fully onboarded with appropriate access, security policies, and monitoring

### 5. System Administrator: Perform Platform Upgrade

**Goal**: Upgrade Kyros platform to latest version with minimal disruption.

**Steps**:
1. Review release notes for target version
2. Check compatibility matrix for dependencies
3. Schedule maintenance window during low-usage period
4. Notify all stakeholders of upcoming maintenance
5. Backup current configuration and metadata
6. Verify backup integrity
7. Pre-pull new container images for all services
8. Enable maintenance mode (read-only for users)
9. Monitor active operations and wait for quiescence
10. Begin rolling update of Kyros services:
    - Update auth service first (dependency for others)
    - Verify authentication still functional
    - Update registry service
    - Verify image push/pull functionality
    - Update trust score service
    - Verify scanning capabilities
    - Update webhook service
    - Verify delivery functionality
    - Update API service
    - Verify dashboard functionality
    - Update operator service last
11. Monitor health checks throughout upgrade process
12. Verify all services report healthy status
13. Disable maintenance mode
14. Perform smoke tests:
    - Login with test user
    - Push and pull test image
    - Verify trust score calculation
    - Test webhook delivery
    - Check dashboard updates
15. Monitor system metrics for anomalies
16. Compare performance metrics to baseline
17. Validate that all configured policies still function
18. Verify audit log continuity
19. Send post-maintenance report to stakeholders
20. Document any issues and lessons learned for future upgrades

**Touchpoints**: System Administration, Service Management, Backup/Recovery, Monitoring, Notification Systems
**Emotional Journey**: Prepared → Methodical → Attentive → Relieved → Confident
**Success Criteria**: Successful upgrade with zero data loss, minimal downtime, and all functionality preserved

### 6. Developer: Collaborate on Image Security Remediation

**Goal**: Work with security team to resolve a vulnerability in a shared base image.

**Steps**:
1. Receive notification from security team about vulnerability in base image
2. Check inbox for detailed vulnerability report
3. Navigate to affected image in Kyros: `registry.example.com/base-images/nodejs:18`
4. View vulnerability details: CVE-2023-24534 in npm package
5. See that fix is available in newer version of npm
6. Fork the base image repository in internal Git registry
7. Update Dockerfile to use newer Node.js base that includes patched npm
8. Rebuild base image locally
9. Test application build with new base image
10. Push updated base image to staging repository
11. Request security team review of new image
12. Security team scans new image and confirms vulnerability resolved
13. Promote image to production repository via automated pipeline
14. Notify application teams of base image update
15. Monitor for any compatibility issues
16. Update base image documentation with new version information
17. Schedule removal of old vulnerable image after grace period
18. Document the remediation process for future reference
19. Participate in retrospective to improve base image update process

**Touchpoints**: Notification System, Image Details, Security Reports, Collaboration Tools (external)
**Emotional Journey**: Notified → Investigative → Proactive → Collaborative → Resolved
**Success Criteria**: Vulnerability resolved, updated image deployed, knowledge shared

### 7. Security Engineer: Implement New Security Policy

**Goal**: Deploy a new organizational security policy for container images.

**Steps**:
1. Review latest threat intelligence and internal security requirements
2. Draft new policy requiring SBOM for all images
3. Policy details:
    - Mandatory SBOM generation for all pushed images
    - SBOM must be in SPDX or CycloneDX format
    - Images without valid SBOR cannot be promoted to production
    - Exception process for legacy applications
4. Navigate to Policy Management in Kyros
5. Click "Create New Policy"
6. Enter policy name: "SBOM Requirement for Production"
7. Select policy scope: "global"
8. Define policy rules using Rego:
    ```rego
    package kyros.policy.sbom_requirement
    
    deny {
        input.event.type == "registry.artifact.pushed"
        not input.data.artifact.sbom_present
        input.data.repository.namespace == "production"
        reason := "Production images must have SBOM"
    }
    
    warn {
        input.event.type == "registry.artifact.pushed"
        not input.data.artifact.sbom_present
        input.data.repository.namespace != "production"
        suggestion := "Consider adding SBOM for better traceability"
    }
    ```
9. Save policy and set to "enabled" mode
10. Configure policy evaluation timing:
    - Pre-push: Optional (for immediate feedback)
    - Post-push: Required (enforces blocking)
    - On-demand: Available for manual checks
11. Set up notifications for policy violations:
    - Email to image uploader
    - Slack notification to #security-alerts channel
    - Ticket creation in JIRA for tracking
12. Communicate policy change to all development teams
13. Provide grace period for compliance (2 weeks)
14. Monitor policy evaluation metrics:
    - Number of images processed
    - Number of violations blocked
    - Number of warnings issued
15. Provide support to teams struggling with compliance
16. After grace period, switch to strict enforcement mode
17. Conduct quarterly review of policy effectiveness
18. Update policy based on feedback and evolving requirements

**Touchpoints**: Policy Management, Notification Configuration, Communication Channels
**Emotional Journey**: Analytical → Deliberative → Authoritative → Supportive → Vigilant
**Success Criteria**: Policy successfully implemented, monitored, and enforced with appropriate flexibility

### 8. System Administrator: Disaster Recovery Drill

**Goal**: Validate disaster recovery procedures and team readiness.

**Steps**:
1. Schedule DR drill with advance notice to stakeholders
2. Notify participants of drill timing and objectives
3. Simulate primary region outage (network partition)
4. Activate incident response plan
5. Verify monitoring systems detect service degradation
6. Confirm alert routing to on-call team
7. Validate failover procedures to secondary region:
    - DNS failover to secondary region endpoints
    - Verify data replication lag is within RPO
    - Check that all services start in secondary region
    - Confirm external dependencies (Keycloak, etc.) are available
8. Test critical user journeys in secondary region:
    - Login and authentication
    - Image push and pull operations
    - Trust score calculation and retrieval
    - Webhook delivery to configured endpoints
9. Measure RTO (Recovery Time Objective) and RPO (Recovery Point Objective)
10. Validate data integrity after failover:
    - Spot-check critical metadata
    - Verify trust scores match pre-failover state
    - Check audit log continuity
11. Test notification systems in secondary region
12. Simulate gradual recovery of primary region
13. Execute failback procedures to primary region
14. Validate successful return to normal operations
15. Conduct post-drill retrospective with all participants
16. Document lessons learned and update runbooks
17. Identify gaps in monitoring, procedures, or training
18. Schedule improvements and follow-up actions
19. Communicate results to executive stakeholders
20. Update DR schedule based on findings and compliance requirements

**Touchpoints**: Monitoring Systems, Incident Response Tools, Failover Mechanisms, Validation Procedures
**Emotional Journey**: Prepared → Alert → Coordinated → Focused → Relieved → Reflective
**Success Criteria**: Successful failover and failback with measured RTO/RPO within targets, validated procedures, and improved readiness

## Journey Mapping Guidelines

### Elements of a Good User Journey
1. **Clear Goal**: Each journey starts with a specific, measurable objective
2. **Persona-Centric**: Written from the perspective of a specific user type
3. **End-to-End**: Covers complete workflow from initiation to resolution
4. **Touchpoint Identification**: Highlights all interactions with the system
5. **Emotional Arc**: Captures the user's emotional experience throughout
6. **Success Criteria**: Defines what constitutes a successful outcome
7. **Pain Points**: Identifies potential frustrations or difficulties
8. **Opportunities**: Notes areas for improvement or enhancement

### Using These Journeys
- **Design Validation**: Verify that proposed designs support these journeys
- **Gap Analysis**: Identify missing features or capabilities
- **Training Development**: Create role-based training materials
- **Process Improvement**: Refine workflows based on identified friction points
- **Metrics Definition**: Establish KPIs based on journey success criteria
- **Communication**: Align teams around user-centered objectives and user experiences

## Conclusion

These user journeys represent the core workflows that the Kyros platform must support to meet the needs of its diverse user base. By designing for these journeys, Kyros ensures that it delivers value to platform engineers, developers, security engineers, organization administrators, and system administrators alike.

The journeys highlight the importance of seamless integration between different platform components, clear communication pathways, and intuitive user interfaces that reduce cognitive load and enable efficient task completion.