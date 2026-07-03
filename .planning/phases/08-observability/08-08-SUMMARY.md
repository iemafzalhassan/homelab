---
status: complete
phase: 08
plan: 08
---

# Plan 08-08 Summary: ArgoCD Metrics Scraping

## What was built
- Enabled per-component metrics in the ArgoCD Helm chart values (`metrics.enabled: true` for controller, server, repoServer, and applicationSet)
- Deployed four `ServiceMonitor` CRDs mapped to each respective metric service to allow Alloy to scrape `argocd_app_info` and other key metrics.
- Wired the `servicemonitors.yaml` into the `argocd-extras` Application.

## Key Files
### Created
- manifests/bootstrap/argocd/servicemonitors.yaml

### Modified
- manifests/bootstrap/argocd/values.yaml
- manifests/apps/argocd-extras.yaml

## Issues Encountered
- Initial approach of using `global.addPrometheusAnnotations: true` was insufficient as the newer version of the ArgoCD Helm chart hides the actual metrics Services behind component-specific `metrics.enabled: true` flags.
- Labels for some ArgoCD metrics Services share the same `app.kubernetes.io/name` (`argocd-metrics`), so `app.kubernetes.io/component` was required in the ServiceMonitor selectors to correctly target them.

## Verification
- Confirmed the 4 ArgoCD metrics Services are deployed on their correct ports (8080, 8082, 8083, 8084).
- Confirmed the 4 ServiceMonitors are deployed and match the service selectors correctly using `app.kubernetes.io/component` and `app.kubernetes.io/name`.

## Self-Check
- [x] All tasks executed
- [x] Each task committed individually
- [x] SUMMARY.md created in plan directory
- [x] STATE.md updated with position and decisions
- [x] ROADMAP.md updated with plan progress
