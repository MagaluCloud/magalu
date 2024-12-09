resource "mgc_dbaas_instances_snapshots" "backup" {
  description  = "My description"
  instance_id  = mgc_dbaas_instances.instance.id
  name         = "my-name"
}