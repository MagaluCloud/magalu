data "mgc_network_vpc" "example" {
  id = mgc_network_vpc.example.id
}

output "datasource_vpc_id" {
  value      = data.mgc_network_vpc.example
}
