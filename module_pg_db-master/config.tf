terraform {
  required_version = ">= 0.12"
  required_providers {
    dbaas-postgres = {
      source  = "ingenico/dbaas-postgres"
      version = "~> 1.1.0"
    }
  }
}
