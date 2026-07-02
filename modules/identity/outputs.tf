output "principal_ids" {
  description = "A map of identity names to their Principal IDs."
  value = {
    for k, v in azurerm_user_assigned_identity.identity : k => v.principal_id
  }
}

output "client_ids" {
  description = "A map of identity names to their Client IDs."
  value = {
    for k, v in azurerm_user_assigned_identity.identity : k => v.client_id
  }
}

output "identity_ids" {
  description = "A map of identity names to their Resource IDs."
  value = {
    for k, v in azurerm_user_assigned_identity.identity : k => v.id
  }
}
