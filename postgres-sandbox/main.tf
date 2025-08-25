module "module_pg_db_demo11" {
  source     = "database/module_pg_db/database"
  token      = var.token
  name       = "yourdbjob11"
  tenant     = "core"
  dbname     = "core_iso1"
  dbsize     = "28"
  dbconn     = "551"
  db_release = "15"
  db_win     = "2"

  role_map = [
    {
      rol_type    = "user"
      rol_name    = "mydatabase_admin"
      rol_group   = "mydatabase_rw" # rw or ro or adm
      rol_conn    = 20
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
    }
  ]

  hba_map = [
    {
      hba_role       = "mydatabase_access"
      hba_addr       = "192.168.1.0/24"
      hba_src_tenant = "my_network"
      hba_auth       = "xxxx"
    }
  ]
}
