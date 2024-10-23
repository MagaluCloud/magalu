data "mgc_network_vpcs_interface" "example" {
  id = mgc_network_vpcs_interfaces.example.id
}

output "datasource_vpcs_interface" {
  value = data.mgc_network_vpcs_interface.example
}
