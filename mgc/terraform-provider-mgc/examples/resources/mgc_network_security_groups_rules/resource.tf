resource "mgc_network_security_groups_rules" "allow_ssh" {
  description      = "Allow incoming SSH traffic"
  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_max   = 22
  port_range_min   = 22
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
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
