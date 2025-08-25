resource "dbaas-postgres_pg_db" "db" {
  metadata {
    name       = var.name
    type       = "database-postgres"
    affinity   = "all"
    timeout    = 150
    retry      = true
    wait_retry = 15
  }

  tenant     = var.tenant
  dbname     = var.dbname
  dbsize     = var.dbsize
  dbconn     = var.dbconn
  db_release = var.db_release
  db_win     = var.db_win

  dynamic "role_map" {
    for_each = var.role_map
    content {
      rol_type    = role_map.value.rol_type
      rol_name    = role_map.value.rol_name
      rol_group   = role_map.value.rol_group
      rol_conn    = role_map.value.rol_conn
      rol_timeout = role_map.value.rol_timeout
    }
  }

  dynamic "extension_map" {
    for_each = var.extension_map
    content {
      ext_name    = extension_map.value.ext_name
      ext_opt_map = extension_map.value.ext_opt_map
    }
  }

  dynamic "hba_map" {
    for_each = var.hba_map
    content {
      hba_role       = hba_map.value.hba_role
      hba_addr       = hba_map.value.hba_addr
      hba_src_tenant = hba_map.value.hba_src_tenant
    }
  }
}

