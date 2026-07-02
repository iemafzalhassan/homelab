output "resource_group_name" {
  description = "The name of the resource group"
  value       = azurerm_resource_group.rg.name
}

output "vnet_id" {
  description = "The ID of the virtual network"
  value       = module.networking.vnet_id
}

output "system_subnet_id" {
  description = "The ID of the system subnet"
  value       = module.networking.system_subnet_id
}

output "user_subnet_id" {
  description = "The ID of the user subnet"
  value       = module.networking.user_subnet_id
}

output "infra_subnet_id" {
  description = "The ID of the infra subnet"
  value       = module.networking.infra_subnet_id
}

output "ingress_subnet_id" {
  description = "The ID of the ingress subnet"
  value       = module.networking.ingress_subnet_id
}

output "aks_cluster_name" {
  description = "The name of the AKS cluster"
  value       = module.aks.cluster_name
}

output "kube_config_raw" {
  description = "Raw Kubernetes config to be used by kubectl and other tools"
  value       = module.aks.kube_config_raw
  sensitive   = true
}

output "oidc_issuer_url" {
  description = "The OIDC issuer URL that is associated with the cluster"
  value       = module.aks.oidc_issuer_url
}

output "identity_client_ids" {
  description = "The Client IDs for the created user-assigned managed identities"
  value       = module.identity.client_ids
}

output "keyvault_uri" {
  description = "The URI of the created Azure Key Vault"
  value       = module.keyvault.vault_uri
}

