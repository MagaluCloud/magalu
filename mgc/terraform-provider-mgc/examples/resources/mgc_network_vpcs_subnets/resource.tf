resource "mgc_network_vpcs_subnets" "example" {
  cidr_block      = "10.0.0.0/16"  
  description     = "Example Subnet"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"] 
  ip_version      = "IPv4"  
  name            = "example-subnet"  
  subnetpool_id   = "subnetpool-12345" 
  vpc_id          = mgc_network_vpcs.example.id  
}
