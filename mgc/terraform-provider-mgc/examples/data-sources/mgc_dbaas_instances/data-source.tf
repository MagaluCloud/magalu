data "mgc_dbaas_instances" "active_instances" {
  status = "ACTIVE"
}

data "mgc_dbaas_instances" "all_instances" {}

data "mgc_dbaas_instances" "deleted_instances" {
  status = "DELETED"
}