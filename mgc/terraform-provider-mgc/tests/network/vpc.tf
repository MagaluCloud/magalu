resource "mgc_network_vpcs" "example" {
  name        = "example-vpc"
}

output "vpc_id" {
  value      = mgc_network_vpcs.example.id
}

data "mgc_network_vpc" "example" {
  id = mgc_network_vpcs.example.id
}

output "datasource_vpc_id" {
  value      = data.mgc_network_vpc.example
}

resource "mgc_network_vpcs_interfaces" "interface_example" {
    name = "example-interface"
    vpc_id = mgc_network_vpcs.example.id
}

output "interface_id" {
    value = mgc_network_vpcs_interfaces.interface_example.id
}

data "mgc_network_vpcs_interface" "interface_example" {
    id = mgc_network_vpcs_interfaces.interface_example.id
}

output "datasource_interface_id" {
    value = data.mgc_network_vpcs_interface.interface_example
}
