output "vnet_id" {
  description = "The ID of the virtual network"
  value       = azurerm_virtual_network.vnet.id
}

output "system_subnet_id" {
  description = "The ID of the system subnet"
  value       = azurerm_subnet.system.id
}

output "user_subnet_id" {
  description = "The ID of the user subnet"
  value       = azurerm_subnet.user.id
}

output "infra_subnet_id" {
  description = "The ID of the infra subnet"
  value       = azurerm_subnet.infra.id
}

output "ingress_subnet_id" {
  description = "The ID of the ingress subnet"
  value       = azurerm_subnet.ingress.id
}
