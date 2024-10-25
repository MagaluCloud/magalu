resource "mgc_network_public_ips" "example" {
  description = "example public ip"
  vpc_id      = mgc_network_vpc.example.id
}

output "public_ip_address" {
  value = mgc_network_public_ips.example.id
}

output "public_ip_id" {
  value = mgc_network_public_ips.example.public_ip
}

#TODO output "datasource_public_ip_id" {
# value = data.mgc_network_public_ips.example
# }
