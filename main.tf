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
