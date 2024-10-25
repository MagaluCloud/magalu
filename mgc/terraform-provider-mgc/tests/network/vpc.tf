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

resource "mgc_network_vpcs_subnets" "example" {
  cidr_block      = "10.0.0.0/16"  
  description     = "Example Subnet"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"] 
  ip_version      = "IPv4"  
  name            = "example-subnet"  
  subnetpool_id   = "subnetpool-12345" 
  vpc_id          = mgc_network_vpcs.example.id  
}

output "subnet_id" {
  value = mgc_network_vpcs_subnets.example.id
}

data "mgc_network_vpcs_subnet" "example" {
  id = mgc_network_vpcs_subnets.example.id
}

output "datasource_subnet_id" {
  value = data.mgc_network_vpcs_subnet.example
}
