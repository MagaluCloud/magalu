resource "mgc_dbaas_instances" "test_instance" {
  name                 = "test-instance"
  user                 = "dbadmin"
  password             = "examplepassword"
  engine_name          = "mysql"
  engine_version       = "8.0"
  instance_type        = "cloud-dbaas-gp1.small"
  volume_size          = 50
  backup_retention_days = 10
  backup_start_at      = "16:00:00"
}
