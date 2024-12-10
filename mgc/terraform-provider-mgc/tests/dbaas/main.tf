# Data sources
data "mgc_dbaas_engines" "active_engines" {
  status = "ACTIVE"
}

data "mgc_dbaas_engines" "deprecated_engines" {
  status = "DEPRECATED"
}

data "mgc_dbaas_engines" "all_engines" {}

# Instance Types data sources
data "mgc_dbaas_instance_types" "active_instance_types" {
  status = "ACTIVE"
}

data "mgc_dbaas_instance_types" "deprecated_instance_types" {
  status = "DEPRECATED"
}

data "mgc_dbaas_instance_types" "default_instance_types" {}

# DBaaS Instances data sources
data "mgc_dbaas_instances" "active_instances" {
  status = "ACTIVE"
}

data "mgc_dbaas_instances" "all_instances" {}

data "mgc_dbaas_instances" "deleted_instances" {
  status = "DELETED"
}

# Get specific instance test
data "mgc_dbaas_instance" "test_instance" {
  id = data.mgc_dbaas_instances.all_instances.instances[0].id
}

# DBaaS Instance resource
resource "mgc_dbaas_instances" "test_instance" {
  name                 = "test-instance-2"
  user                 = "dbadmin"
  password             = "aaaaaaaaaa"
  engine_name          = "mysql"
  engine_version       = "8.0"
  instance_type        = "cloud-dbaas-gp1.small"
  volume_size          = 50
  backup_retention_days = 10
  backup_start_at      = "16:00:00"
  
  # parameters = {
  #   character_set_server = "utf8mb4"
  #   max_connections     = "100"
  # }
}

# Outputs for debugging
output "active_engines" {
  value = data.mgc_dbaas_engines.active_engines.engines
}

output "deprecated_engines" {
  value = data.mgc_dbaas_engines.deprecated_engines.engines
}

output "all_engines" {
  value = data.mgc_dbaas_engines.all_engines.engines
}

# Additional outputs for debugging
output "active_instance_types" {
  value = data.mgc_dbaas_instance_types.active_instance_types.instance_types
}

output "deprecated_instance_types" {
  value = data.mgc_dbaas_instance_types.deprecated_instance_types.instance_types
}

output "default_instance_types" {
  value = data.mgc_dbaas_instance_types.default_instance_types.instance_types
}

output "active_instances" {
  value = data.mgc_dbaas_instances.active_instances.instances
}

output "all_instances" {
  value = data.mgc_dbaas_instances.all_instances.instances
}

output "deleted_instances" {
  value = data.mgc_dbaas_instances.deleted_instances.instances
}

# Output for the test instance
output "test_instance" {
  value = data.mgc_dbaas_instance.test_instance
}

# Optional: Output the instance details
output "dbaas_instance" {
  value = {
    id           = mgc_dbaas_instances.test_instance.id
    name         = mgc_dbaas_instances.test_instance.name
    engine_name  = mgc_dbaas_instances.test_instance.engine_name
    instance_type = mgc_dbaas_instances.test_instance.instance_type
    volume_size  = mgc_dbaas_instances.test_instance.volume_size
  }

  sensitive = true # Because it contains instance information
}
