data "mgc_virtual_machine_instances" "instances" {
}

output "vm_instances" {
  value = data.mgc_virtual_machine_instances.instances
}