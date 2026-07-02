resource "azurerm_user_assigned_identity" "identity" {
  for_each            = var.identities
  name                = "${each.key}-identity"
  resource_group_name = var.resource_group_name
  location            = var.location

  tags = {
    environment = "homelab"
    managed-by  = "terraform"
    workload    = each.key
  }
}

resource "azurerm_federated_identity_credential" "federated_credential" {
  for_each            = var.identities
  name                = "${each.key}-fic"
  audience            = ["api://AzureADTokenExchange"]
  issuer                    = var.oidc_issuer_url
  user_assigned_identity_id = azurerm_user_assigned_identity.identity[each.key].id
  subject                   = "system:serviceaccount:${each.value.namespace}:${each.value.serviceaccount}"
}
