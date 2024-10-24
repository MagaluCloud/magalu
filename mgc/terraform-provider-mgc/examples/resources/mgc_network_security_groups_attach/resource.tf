resource "mgc_network_security_groups_attach" "attach_example" {
  security_group_id = mgc_network_security_groups.example.id
  interface_id = mgc_network_vpcs_interfaces.interface_example.id
}

output "security_group_attach_id" {
  value = mgc_network_security_groups_attach.attach_example
}
