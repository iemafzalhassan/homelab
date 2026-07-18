---
status: verifying
trigger: "OpenTelemetry Demo frontend pod was not becoming ready and ingress/service validation was needed"
created: 2026-07-18T01:53:00+05:30
updated: 2026-07-18T01:53:00+05:30
---

## Current Focus
hypothesis: The frontend container was listening on port 3000 while the deployment probe and service targeted port 8080.
test: Reapply the manifest with the probe/service port aligned to 3000 and verify pod readiness plus traffic flow.
expecting: The frontend pod should become Ready and respond on the service path.
next_action: Validate service reachability and telemetry ingestion in the monitoring stack.

## Symptoms
expected: The OpenTelemetry Demo frontend should become Ready and serve traffic in the demo-prod namespace.
actual: The pod was scheduled and running but stayed Unready because readiness probes to port 8080 failed with connection refused.
errors: Readiness probe failed: Get "http://10.244.1.41:8080/": dial tcp 10.244.1.41:8080: connect: connection refused
reproduction: Deploy the frontend workload and observe the pod readiness state.
started: 2026-07-18

## Evidence
- The container logs show the Next.js frontend app starting successfully and reporting the app on port 3000.
- The deployment spec and service were updated to target port 3000.
- The frontend pod became Ready after the manifest change.

## Resolution
root_cause: The OpenTelemetry frontend image listens on port 3000, while the deployment had been configured to probe and expose port 8080.
fix: Align the deployment container port, service target port, and readiness probe with the application’s actual port 3000.
verification: Pending end-to-end validation through the service path and monitoring pipeline.
files_changed: [manifests/apps/opentelemetry-demo/prod/frontend.yaml]
