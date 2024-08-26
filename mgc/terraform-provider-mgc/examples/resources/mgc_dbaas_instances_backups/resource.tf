resource "mgc_dbaas_instances_backups" "backup" {
  instance_id = mgc_dbaas_instances.instance.id
  mode        = "FULL"
}