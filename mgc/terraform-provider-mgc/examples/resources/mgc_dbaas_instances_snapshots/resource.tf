# Create a snapshot for a DBaaS instance
resource "mgc_dbaas_instances_snapshots" "example" {
  instance_id  = mgc_dbaas_instances.my_instance.id
  name        = "example-snapshot"
  description = "Snapshot created via Terraform"
}