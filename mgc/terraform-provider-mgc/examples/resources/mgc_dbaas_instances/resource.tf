resource "mgc_dbaas_instances" "dbaas_instances" {
  id    =  computed
  name           = "my-database-instance"
  user           = "db_user"
  password       = "secure_password123"
  engine_name    = "mysql"
  engine_version = "8.0"
  instance_type  = "mgc.db.tiny"
  volume_size    = 10
}
