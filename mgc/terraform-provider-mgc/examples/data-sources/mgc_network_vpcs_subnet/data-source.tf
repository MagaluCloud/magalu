data "mgc_network_vpcs_subnet" "example" {
  id = mgc_network_vpcs_subnets.example.id
}

output "datasource_subnet_id" {
  value = data.mgc_network_vpcs_subnet.example
}
