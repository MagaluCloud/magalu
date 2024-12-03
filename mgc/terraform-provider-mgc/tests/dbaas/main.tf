# Data sources
data "mgc_dbaas_engines" "active_engines" {
  status = "ACTIVE"
}

data "mgc_dbaas_engines" "deprecated_engines" {
  status = "DEPRECATED"
}

data "mgc_dbaas_engines" "all_engines" {}

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