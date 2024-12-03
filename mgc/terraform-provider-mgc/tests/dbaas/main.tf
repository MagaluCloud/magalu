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
