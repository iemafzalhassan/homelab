variable "resource_group_name" {
  type        = string
  description = "The name of the resource group in which to create the AKS cluster."
}

variable "location" {
  type        = string
  description = "The location/region where the AKS cluster is created."
}

variable "vnet_id" {
  type        = string
  description = "The ID of the Virtual Network."
}

variable "system_subnet_id" {
  type        = string
  description = "The ID of the subnet for the system node pool."
}

variable "admin_ip_ranges" {
  type        = list(string)
  description = "The IP address ranges (CIDR) allowed to access the AKS API server."
}
