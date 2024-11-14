data "mgc_network_public_ip" "example" {
  id = mgc_network_public_ips.example.id
}

output "datasource_public_ip_id" {
  value = data.mgc_network_public_ip.example
}
