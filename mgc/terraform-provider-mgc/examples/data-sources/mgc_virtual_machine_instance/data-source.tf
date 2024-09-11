data "mgc_virtual_machine_instances" "instances" {
  id = mgc_virtual_machine_instances.my_vm.id
}

output "vm_instances" {
  value = data.mgc_virtual_machine_instances.instances
}