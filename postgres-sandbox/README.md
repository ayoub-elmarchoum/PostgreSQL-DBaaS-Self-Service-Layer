# postgres 

Docs to read :[DBAAS - Databases](https://confluence.worldline-solutions.com/display/GINFRA/DBAAS+-+Databases)

## Getting started
```
git clone git@gitlab.global.ingenico.com:dbaas/giservices/emea/database/postgres.git
git branch
git checkout branch -b name_branch
git push
```

## Add your files

```
module "dbaas-pg-db_module" {
  source     = "database/module_dbaas/database"
  token      = var.token # token that we give to the tenant and should be added to variable cicd 
  version    = "1.0.9" # don't change it's the version of module 
  tenant     = "TenantName"  
  dbname     = "DbName"
  dbsize     = "50"
  dbconn     = "30"
  db_release = "15" # 15 by default
  db_win     = "3"  # maintenance window Possible values are  " 1: [ 21h00 - 00h00 ] or  2: [ 00h00 - 03h00 ] or  3: [ 03h00 - 05h00 ] or 4: [ 05h00 - 08h00 ] "

  role_map = [
    {
      rol_type    = "user"
      rol_name    = "rolename__admin"
      rol_group   = "roleGroup_rw" # rw or ro or adm
      rol_conn    = 20 # It shouldn't be greater than dbconn !
      rol_timeout = 300
    }
  ]

  extension_map = [
    {
      ext_name = "pgbouncer"
      ext_opt_map = {
        pgb_mode    = "transaction"
        pgb_pool    = "20"
        pgb_idle_to = "300"
      }
    },
    ...
  ]

  hba_map = [
    {
      hba_role       = "mydatabase_access"
      hba_addr       = "192.168.1.0/24"
      hba_src_tenant = "my_network"
      hba_auth       = "xxxx"
    },
    ...
  ]
}

```
## Integrate with your tools


## Collaborate with your team


## Test and Deploy



## Name

## Description

## Badges

## Visuals

## Installation

## Usage

## Support
## Roadmap
If you have ideas for releases in the future, it is a good idea to list them in the README.

## Contributing
##Authors and acknowledgment
 
## License

## Project status
