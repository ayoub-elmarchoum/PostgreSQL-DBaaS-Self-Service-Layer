variable "uri" {
  description = "API URI for dbaas"
  type        = string
  default     = "https://moldapi.services.core.sb.eu.ginfra.net/deds/api/v1/"
}

variable "tenant" {
  description = "Tenant for the database"
  type        = string
  validation {
    condition     = can(regex("^[a-z]{3,20}$", var.tenant))
    error_message = "Le tenant doit contenir entre 3 et 20 caractères alphabétiques minuscules."
  }
}

variable "token" {
  description = "API token for authentication"
  type        = string
}

variable "username" {
  description = "API username for BASIC authentication"
  type        = string
  default     = ""
}

variable "password" {
  description = "API password for BASIC authentication"
  type        = string
  default     = ""
}

variable "insecure" {
  description = "Disables TLS verification if using HTTPS."
  type        = bool
  default     = false
}

variable "timeout" {
  description = "Timeout in seconds for requests."
  type        = number
  default     = 60
}

variable "debug" {
  description = "Enable debug mode to trace requests."
  type        = bool
  default     = false
}


variable "dbname" {
  description = "Database name, prefixed with '[tenant]_db_'"
  type        = string
}
variable "dbsize" {
  description = "Size of the database"
  type        = number
}

variable "dbconn" {
  description = "Maximum number of connections to the database"
  type        = number
}

variable "db_release" {
  description = "Type of the database service"
  type        = string
  default     = "15"
}

variable "db_win" {
  description = "Database window settings"
  type        = number
}

variable "role_map" {
  type = list(object({
    rol_type    = string
    rol_name    = string
    rol_group   = string
    rol_conn    = number
    rol_timeout = number
  }))
}

variable "extension_map" {
  type = list(object({
    ext_name    = string
    ext_opt_map = map(string)
  }))
}

variable "hba_map" {
  type = list(object({
    hba_role       = string
    hba_addr       = string
    hba_src_tenant = string
    hba_auth       = string
  }))
}

variable "name" {
  type = string
}
