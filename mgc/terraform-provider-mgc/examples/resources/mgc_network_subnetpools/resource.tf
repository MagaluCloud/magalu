resource "mgc_network_subnetpools" "main_subnetpool" {
  name        = "main-subnetpool"
  description = "Main Subnet Pool"
  type        = "pip"
  cidr        = "172.26.0.0/16"
}
