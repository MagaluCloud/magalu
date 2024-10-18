resource "mgc_network_security_groups" "example" {
  name                  = "example-security-group"
  description           = "An example security group"
  disable_default_rules = false
}

output "security_group_id" {
  value = mgc_network_security_groups.example
}

resource "mgc_network_security_groups" "example2" {
  name                  = "example-security-group2"
  description           = "An example security group"
  disable_default_rules = true
}

output "security_group_id2" {
  value = mgc_network_security_groups.example2
}

resource "mgc_network_security_groups" "example3" {
  name = "example-security-group3"
}

output "security_group_id3" {
  value = mgc_network_security_groups.example3
}

data "mgc_network_security_group" "example" {
  id = mgc_network_security_groups.example.id
}

output "datasource_security_group_id" {
  value = data.mgc_network_security_group.example
}
