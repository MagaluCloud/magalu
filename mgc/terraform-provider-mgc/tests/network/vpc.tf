resource "mgc_network_vpc" "example" {
  name        = "example-vpc"
}

output "vpc_id" {
  value      = mgc_network_vpc.example.id
}

data "mgc_network_vpc" "example" {
  id = mgc_network_vpc.example.id
}

output "datasource_vpc_id" {
  value      = data.mgc_network_vpc.example
}
