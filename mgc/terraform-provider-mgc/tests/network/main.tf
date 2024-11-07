# Security Groups
resource "mgc_network_security_groups" "primary_sg" {
  name                  = "primary-security-group-tf"
  description           = "Primary security group for main services"
  disable_default_rules = true
}

resource "mgc_network_security_groups" "secondary_sg" {
  name                  = "secondary-security-group"
  description           = "Secondary security group for additional services"
  disable_default_rules = true
}

resource "mgc_network_security_groups" "auxiliary_sg" {
  name = "auxiliary-security-group"
}

data "mgc_network_security_group" "primary_sg_data" {
  id = mgc_network_security_groups.primary_sg.id
}

# Security Group Rules
resource "mgc_network_security_groups_rules" "ssh_ipv4_rule" {
  description       = "Allow incoming SSH traffic"
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "192.168.1.0/24"
  security_group_id = mgc_network_security_groups.primary_sg.id
}

resource "mgc_network_security_groups_rules" "ssh_ipv6_rule" {
  description       = "Allow incoming SSH traffic from IPv6"
  direction         = "ingress"
  ethertype         = "IPv6"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "::/0"
  security_group_id = mgc_network_security_groups.primary_sg.id
}

# VPC Resources
resource "mgc_network_vpcs" "main_vpc" {
  name = "main-vpc"
}

data "mgc_network_vpc" "main_vpc_data" {
  id = mgc_network_vpcs.main_vpc.id
}

# VPC Interfaces
resource "mgc_network_vpcs_interfaces" "primary_interface" {
  name   = "primary-interface"
  vpc_id = "9dd2d30e-565d-42ce-a0a3-f2de1c473fed"
}

data "mgc_network_vpcs_interface" "primary_interface_data" {
  id = mgc_network_vpcs_interfaces.primary_interface.id
}

resource "mgc_network_vpcs_interfaces" "pip_interface" {
  name   = "pip-interface"
  vpc_id = "9dd2d30e-565d-42ce-a0a3-f2de1c473fed"
}

# Security Group Attachment
resource "mgc_network_security_groups_attach" "primary_sg_attachment" {
  security_group_id = mgc_network_security_groups.primary_sg.id
  interface_id      = mgc_network_vpcs_interfaces.primary_interface.id
}

# Subnet Resources
data "mgc_network_subnetpool" "main_subnetpool" {
  id = "0290a302-77b4-4315-801c-087c7b96867b"
}

# resource "mgc_network_vpcs_subnets" "primary_subnet" {
#   cidr_block      = "10.0.0.0/16"  
#   description     = "Primary Network Subnet"
#   dns_nameservers = ["8.8.8.8", "8.8.4.4"] 
#   ip_version      = "IPv4"  
#   name            = "primary-subnet"  
#   subnetpool_id   = "subnetpool-12345" 
#   vpc_id          = mgc_network_vpcs.main_vpc.id  
# }

data "mgc_network_vpcs_subnet" "primary_subnet_data" {
  id = "4a073774-5a74-4bc8-9ef4-405058ed802a"
}

# Public IP
resource "mgc_network_public_ips" "example" {
  description = "example public ip"
  vpc_id      = mgc_network_vpcs.main_vpc.id
}

data "mgc_network_public_ip" "example" {
  id = mgc_network_public_ips.example.id
}

#Public IP Attachment
resource "mgc_network_public_ips_attach" "example" {
  public_ip_id = mgc_network_public_ips.example.id
  interface_id = mgc_network_vpcs_interfaces.pip_interface.id
}

#Subnetpools
# resource "mgc_network_subnetpools_book_cidr" "book_subnetpool" {
#   cidr = "172.18.0.199/32"
#   subnet_pool_id   = "0290a302-77b4-4315-801c-087c7b96867b"
# }

# Outputs
output "primary_security_group_data" {
  value = data.mgc_network_security_group.primary_sg_data
}

output "main_subnetpool_data" {
  value = data.mgc_network_subnetpool.main_subnetpool
}

output "main_vpc_data" {
  value = data.mgc_network_vpc.main_vpc_data
}

output "primary_interface_data" {
  value = data.mgc_network_vpcs_interface.primary_interface_data
}

output "primary_subnet_data" {
  value = data.mgc_network_vpcs_subnet.primary_subnet_data
}

output "datasource_public_ip_id" {
  value = data.mgc_network_public_ip.example
}
