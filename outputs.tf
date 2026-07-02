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
