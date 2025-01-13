resource "mgc_virtual_machine_interface_attach" "attach_vm" {
  instance_id  = mgc_virtual_machine_instances.your_instance.id
  interface_id = mgc_network_vpcs_interfaces.your_interface.id
}
