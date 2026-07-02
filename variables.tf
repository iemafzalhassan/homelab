variable "location" {
  type        = string
  description = "Azure region for the homelab"
  default     = "centralindia"
}

variable "rg_name" {
  type        = string
  description = "Name of the resource group"
  default     = "homelab-rg"
}

variable "budget_contact_email" {
  type        = string
  description = "Email address to receive budget alerts"
}

variable "admin_ip_ranges" {
  type        = list(string)
  description = "The IP address ranges (CIDR) allowed to access the AKS API server."
}
