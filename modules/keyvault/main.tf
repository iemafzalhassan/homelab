resource "random_string" "kv_suffix" {
  length  = 4
  special = false
  upper   = false
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                = "homelab-kv-${random_string.kv_suffix.result}"
  location            = var.location
  resource_group_name = var.resource_group_name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  # Use Azure RBAC instead of Access Policies
  rbac_authorization_enabled    = true

  # Public network access must be true for IP rules to work
  public_network_access_enabled = true

  network_acls {
    default_action = "Deny"
    bypass         = "AzureServices"
    ip_rules       = [for ip in var.admin_ip_ranges : split("/", ip)[0]]
  }

  # Needed for Secret Store CSI
  purge_protection_enabled = false

  tags = {
    environment = "homelab"
    managed-by  = "terraform"
  }
}

# Private Endpoint
resource "azurerm_private_endpoint" "kv_pe" {
  name                = "${azurerm_key_vault.kv.name}-pe"
  location            = var.location
  resource_group_name = var.resource_group_name
  subnet_id           = var.infra_subnet_id

  private_service_connection {
    name                           = "${azurerm_key_vault.kv.name}-psc"
    private_connection_resource_id = azurerm_key_vault.kv.id
    subresource_names              = ["vault"]
    is_manual_connection           = false
  }

  private_dns_zone_group {
    name                 = "kv-dns-zone-group"
    private_dns_zone_ids = [azurerm_private_dns_zone.kv_dns.id]
  }
}

# Private DNS Zone
resource "azurerm_private_dns_zone" "kv_dns" {
  name                = "privatelink.vaultcore.azure.net"
  resource_group_name = var.resource_group_name
}

# VNet Link
resource "azurerm_private_dns_zone_virtual_network_link" "kv_dns_link" {
  name                  = "kv-vnet-link"
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.kv_dns.name
  virtual_network_id    = var.vnet_id
}

# Grant the terraform admin access to create placeholder secrets
resource "azurerm_role_assignment" "admin_kv_access" {
  scope                = azurerm_key_vault.kv.id
  role_definition_name = "Key Vault Secrets Officer"
  principal_id         = var.admin_principal_id
}

# Placeholder secret for Cloudflare API Token
resource "azurerm_key_vault_secret" "cloudflare_api_token" {
  name         = "cloudflare-api-token"
  value        = ""
  key_vault_id = azurerm_key_vault.kv.id

  # Ensure the admin has access before trying to create the secret
  depends_on = [
    azurerm_role_assignment.admin_kv_access
  ]
}
