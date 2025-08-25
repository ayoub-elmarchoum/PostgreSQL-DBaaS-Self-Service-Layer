## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.12 |
| <a name="requirement_dbaas-postgres"></a> [dbaas-postgres](#requirement\_dbaas-postgres) | ~> 1.1.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_dbaas-postgres"></a> [dbaas-postgres](#provider\_dbaas-postgres) | ~> 1.1.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [dbaas-postgres_pg_db.db](https://registry.terraform.io/providers/ingenico/dbaas-postgres/latest/docs/resources/pg_db) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_db_release"></a> [db\_release](#input\_db\_release) | Type of the database service | `string` | `"15"` | no |
| <a name="input_db_win"></a> [db\_win](#input\_db\_win) | Database window settings | `number` | n/a | yes |
| <a name="input_dbconn"></a> [dbconn](#input\_dbconn) | Maximum number of connections to the database | `number` | n/a | yes |
| <a name="input_dbname"></a> [dbname](#input\_dbname) | Database name, prefixed with '[tenant]\_db\_' | `string` | n/a | yes |
| <a name="input_dbsize"></a> [dbsize](#input\_dbsize) | Size of the database | `number` | n/a | yes |
| <a name="input_debug"></a> [debug](#input\_debug) | Enable debug mode to trace requests. | `bool` | `false` | no |
| <a name="input_extension_map"></a> [extension\_map](#input\_extension\_map) | n/a | <pre>list(object({<br/>    ext_name    = string<br/>    ext_opt_map = map(string)<br/>  }))</pre> | n/a | yes |
| <a name="input_hba_map"></a> [hba\_map](#input\_hba\_map) | n/a | <pre>list(object({<br/>    hba_role       = string<br/>    hba_addr       = string<br/>    hba_src_tenant = string<br/>    hba_auth       = string<br/>  }))</pre> | n/a | yes |
| <a name="input_insecure"></a> [insecure](#input\_insecure) | Disables TLS verification if using HTTPS. | `bool` | `false` | no |
| <a name="input_name"></a> [name](#input\_name) | n/a | `string` | n/a | yes |
| <a name="input_password"></a> [password](#input\_password) | API password for BASIC authentication | `string` | `""` | no |
| <a name="input_role_map"></a> [role\_map](#input\_role\_map) | n/a | <pre>list(object({<br/>    rol_type    = string<br/>    rol_name    = string<br/>    rol_group   = string<br/>    rol_conn    = number<br/>    rol_timeout = number<br/>  }))</pre> | n/a | yes |
| <a name="input_tenant"></a> [tenant](#input\_tenant) | Tenant for the database | `string` | n/a | yes |
| <a name="input_timeout"></a> [timeout](#input\_timeout) | Timeout in seconds for requests. | `number` | `60` | no |
| <a name="input_token"></a> [token](#input\_token) | API token for authentication | `string` | n/a | yes |
| <a name="input_uri"></a> [uri](#input\_uri) | API URI for dbaas | `string` | `"https://moldapi.services.core.sb.eu.ginfra.net/deds/api/v1/"` | no |
| <a name="input_username"></a> [username](#input\_username) | API username for BASIC authentication | `string` | `""` | no |

## Outputs

No outputs.
