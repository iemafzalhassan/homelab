# ADR 01: Remove Kargo and Multi-Environment Setup

## Context
During the deployment of the `microservices-demo` application, we encountered severe quota limits on the Azure subscription in the Central India region:
- **Spot vCPU Quota:** Limited to exactly 3 vCPUs. Since our Spot node pool uses `Standard_D2as_v5` (2 vCPUs per node), we are hard-capped at a single spot node. The cluster autoscaler failed with `OperationNotAllowed` when attempting to scale to 2 nodes.
- **Public IP Quota:** Limited to 3 Public IPs in the region, which causes `LoadBalancer` services to fail if we exceed the limit.

## Decision
To fit within these strict budget and quota constraints, we made the following architectural changes:

1. **Removed Kargo:** We completely removed the Kargo multi-environment promotion pipeline (UAT / Prod). Running duplicate environments for a 11-microservice application is impossible on a single 2-core node. We reverted to a single `microservices-demo` namespace managed directly by ArgoCD.
2. **Stripped Resource Limits:** We removed all CPU and memory `requests` and `limits` from the official Google `microservices-demo` manifests. This allows the Kubernetes scheduler to pack all 11 microservices densely onto the single available Spot node without triggering a scale-up event.
3. **ClusterIP for Ingress:** We modified the `frontend-external` service from `LoadBalancer` to `ClusterIP` to avoid requesting a new Public IP. Traffic will be routed through the existing Traefik Gateway API ingress.

## Consequences
- **Positive:** The cluster remains within the $100/month budget and operates successfully within the Azure quota limits.
- **Negative:** We lose the ability to demonstrate multi-stage GitOps promotion (UAT -> Prod) in this homelab. We also risk CPU throttling under heavy load since resource requests are omitted, though this is acceptable for a homelab environment.
