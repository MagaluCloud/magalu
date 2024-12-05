# Data sources
data "mgc_dbaas_engines" "active_engines" {
  status = "ACTIVE"
}

data "mgc_dbaas_engines" "deprecated_engines" {
  status = "DEPRECATED"
}

data "mgc_dbaas_engines" "all_engines" {}