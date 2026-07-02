terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

locals {
  tags = {
    environment = "homelab"
    managed-by  = "terraform"
  }
}

resource "azurerm_resource_group" "rg" {
  name     = var.rg_name
  location = var.location
  tags     = local.tags
}

module "networking" {
  source              = "./modules/networking"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  tags                = local.tags
}

resource "azurerm_consumption_budget_resource_group" "budget" {
  name              = "${var.rg_name}-budget"
  resource_group_id = azurerm_resource_group.rg.id
  amount            = 25
  time_grain        = "Monthly"

  time_period {
    start_date = "2026-07-01T00:00:00Z"
    end_date   = "2027-07-01T00:00:00Z"
  }

  notification {
    enabled        = true
    threshold      = 80.0
    operator       = "EqualTo"
    contact_emails = [var.budget_contact_email]
  }

  notification {
    enabled        = true
    threshold      = 100.0
    operator       = "GreaterThanOrEqualTo"
    contact_emails = [var.budget_contact_email]
  }
}

module "aks" {
  source              = "./modules/aks"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  vnet_id             = module.networking.vnet_id
  system_subnet_id    = module.networking.system_subnet_id
  admin_ip_ranges     = var.admin_ip_ranges
}

data "azurerm_client_config" "current" {}

module "identity" {
  source              = "./modules/identity"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  oidc_issuer_url     = module.aks.oidc_issuer_url
  identities = {
    argocd = {
      namespace      = "argocd"
      serviceaccount = "argocd-repo-server"
    }
    jenkins = {
      namespace      = "jenkins"
      serviceaccount = "jenkins"
    }
    traefik = {
      namespace      = "traefik"
      serviceaccount = "traefik"
    }
    cert-manager = {
      namespace      = "cert-manager"
      serviceaccount = "cert-manager"
    }
  }
}

module "keyvault" {
  source              = "./modules/keyvault"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  vnet_id             = module.networking.vnet_id
  infra_subnet_id     = module.networking.infra_subnet_id
  admin_principal_id  = data.azurerm_client_config.current.object_id
  admin_ip_ranges     = var.admin_ip_ranges
}

resource "azurerm_role_assignment" "kv_secrets_user" {
  for_each             = module.identity.principal_ids
  scope                = module.keyvault.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = each.value
}

module "acr" {
  source              = "./modules/acr"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  tags                = local.tags
}

resource "azurerm_role_assignment" "acr_push_jenkins" {
  scope                = module.acr.id
  role_definition_name = "AcrPush"
  principal_id         = module.identity.principal_ids["jenkins"]
}
