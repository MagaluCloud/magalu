resource "mgc_dbaas_instances" "dbaas_instances" {
  name             = "my-database-instance"
  user             = "db_user"
  password         = "secure_password123"
  engine_id        = "063f3994-b6c2-4c37-96c9-bab8d82d36f7" #mysql 8.0 - please check the available version
  instance_type_id = "c460d5c1-883d-4fea-afc3-1a208e982084" #cloud-dbaas-bs1.medium - please check the available engine
  volume = {
    size = 50 # Size in GiB
  }
}
