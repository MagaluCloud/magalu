resource "mgc_network_security_groups" "example" {
  name                  = "example-security-group"
  description           = "An example security group"
  disable_default_rules = false
}

output "security_group_id" {
  value = mgc_network_security_groups.example
}

resource "mgc_network_security_groups" "example2" {
  name                  = "example-security-group3"
}

output "security_group_id2" {
  value = mgc_network_security_groups.example2
}
