variable "resource_group_name" {
  type        = string
  description = "The name of the resource group."
}

variable "location" {
  type        = string
  description = "The location/region where the Key Vault should be created."
}

variable "vnet_id" {
  type        = string
  description = "The ID of the Virtual Network."
}

variable "infra_subnet_id" {
  type        = string
  description = "The ID of the infrastructure subnet for the private endpoint."
}

variable "admin_principal_id" {
  type        = string
  description = "The Object ID of the admin user running Terraform (to manage placeholder secrets)."
}

variable "admin_ip_ranges" {
  type        = list(string)
  description = "The IP address ranges (CIDR) allowed to access the Key Vault data plane."
}
