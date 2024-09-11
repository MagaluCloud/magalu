resource "mgc_dbaas_instances" "dbaas_instances" {
  name      = "my-database-instance"
  flavor_id = "your-flavor-id"
  user      = "db_user"
  password  = "secure_password123"

  volume {
    size = 50 # Size in GiB
  }
}
