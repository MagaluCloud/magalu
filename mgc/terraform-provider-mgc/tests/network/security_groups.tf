resource "mgc_network_security_groups" "example" {
  name                  = "example-security-group-tf"
  description           = "An example security group"
  disable_default_rules = true
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

resource "mgc_network_security_groups_rules" "example" {
  description      = "Allow incoming SSH traffic"
  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_max   = 22
  port_range_min   = 22
  protocol         = "tcp"
  remote_ip_prefix = "192.168.1.0/24"
  security_group_id = mgc_network_security_groups.example.id
}

resource "mgc_network_security_groups_rules" "allow_ssh_ipv6" {
  description      = "Allow incoming SSH traffic from IPv6"
  direction        = "ingress"
  ethertype        = "IPv6"
  port_range_max   = 22
  port_range_min   = 22
  protocol         = "tcp"
  remote_ip_prefix = "::/0"
  security_group_id = mgc_network_security_groups.example.id
}