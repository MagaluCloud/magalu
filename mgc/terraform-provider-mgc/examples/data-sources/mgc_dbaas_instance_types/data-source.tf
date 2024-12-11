data "mgc_dbaas_instance_types" "active_instance_types" {
  status = "ACTIVE"
}

data "mgc_dbaas_instance_types" "deprecated_instance_types" {
  status = "DEPRECATED"
}

data "mgc_dbaas_instance_types" "all_instance_types" {}
