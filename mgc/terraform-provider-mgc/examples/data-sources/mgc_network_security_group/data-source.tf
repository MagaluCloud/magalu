data "mgc_network_security_groups" "example" {
  id = mgc_network_security_groups.example.id
}

output "datasource_security_group_id" {
  value = data.mgc_network_security_groups.example
}
