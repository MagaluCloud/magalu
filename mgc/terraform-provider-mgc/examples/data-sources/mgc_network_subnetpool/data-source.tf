data "mgc_network_subnetpool" "example" {
  id = "subnetpool-id" 
}

output "subnetpool_cidr" {
  value = data.mgc_network_subnetpool.example
}
