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

resource "azurerm_kubernetes_cluster_node_pool" "spot" {
  name                  = "spot"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.aks.id
  vm_size               = "Standard_D2as_v5"
  
  priority              = "Spot"
  eviction_policy       = "Delete"
  spot_max_price        = -1
  
  zones                 = ["1", "2"]
  
  auto_scaling_enabled  = true
  min_count             = 2
  max_count             = 4
  
  os_disk_type          = "Ephemeral"
  os_disk_size_gb       = 30
  vnet_subnet_id        = var.system_subnet_id

  node_taints = [
    "kubernetes.azure.com/scalesetpriority=spot:NoSchedule"
  ]
  
  node_labels = {
    "kubernetes.azure.com/scalesetpriority" = "spot"
  }

  tags = {
    environment = "homelab"
    managed-by  = "terraform"
  }
}
