data "mgc_network_security_group" "example" {
  id = mgc_network_security_groups.example.id
}

output "datasource_security_group" {
  value = data.mgc_network_security_group.example
}
