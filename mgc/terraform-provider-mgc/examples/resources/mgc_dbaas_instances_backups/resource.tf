# Create a full backup for a DBaaS instance
resource "mgc_dbaas_instances_backups" "example" {
  instance_id = mgc_dbaas_instances.my_instance.id
  mode       = "FULL"
}
