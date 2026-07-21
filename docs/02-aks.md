# Track B: Azure AKS + Terraform + ArgoCD Production Guide

Production-patterned deployment guide on Azure Kubernetes Service (AKS) with GitOps.

---

## Infrastructure Overview
- **Azure AKS**: Free-tier control plane with D2as_v5 system node pool and Spot user node pool.
- **ArgoCD**: GitOps declarative deployment.
- **Traefik Gateway API**: Ingress & HTTPRoute management.
- **Grafana Alloy**: Telemetry collection agent shipping logs, metrics, and traces to Grafana Cloud.

---

## Deployment Steps

1. Provision infrastructure via Terraform:
   ```bash
   cd terraform
   terraform init
   terraform apply
   ```

2. Register ArgoCD Root Application:
   ```bash
   kubectl apply -f manifests/apps/root.yaml
   ```

3. ArgoCD automatically reconciles `opentelemetry-demo` and `opentelemetry-demo-extras`.
