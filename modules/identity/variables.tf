variable "resource_group_name" {
  type        = string
  description = "The name of the resource group."
}

variable "location" {
  type        = string
  description = "The location/region where the identities should be created."
}

variable "oidc_issuer_url" {
  type        = string
  description = "The OIDC issuer URL of the AKS cluster."
}

variable "identities" {
  type = map(object({
    namespace      = string
    serviceaccount = string
  }))
  description = "Map of identity names to their corresponding Kubernetes namespace and serviceaccount."
}
