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