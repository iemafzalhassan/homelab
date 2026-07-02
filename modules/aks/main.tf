resource "azurerm_kubernetes_cluster" "aks" {
  name                = "homelab-aks"
  location            = var.location
  resource_group_name = var.resource_group_name
  dns_prefix          = "homelab"

  kubernetes_version = "1.36" # Will use latest patch

  oidc_issuer_enabled       = true
  workload_identity_enabled = true

  network_profile {
    network_plugin      = "azure"
    network_plugin_mode = "overlay"
    pod_cidr            = "10.244.0.0/16"
    service_cidr        = "10.96.0.0/16"
    dns_service_ip      = "10.96.0.10"
  }

  default_node_pool {
    name                         = "system"
    vm_size                      = "Standard_D2s_v3"
    zones                        = ["1"]
    node_count                   = 1
    os_disk_type                 = "Managed"
    os_disk_size_gb              = 30
    vnet_subnet_id               = var.system_subnet_id
    only_critical_addons_enabled = true # Best practice for system node pool
  }

  identity {
    type = "SystemAssigned"
  }

  api_server_access_profile {
    authorized_ip_ranges = var.admin_ip_ranges
  }

  tags = {
    environment = "homelab"
    managed-by  = "terraform"
  }
}
