resource "mgc_network_public_ips" "example" {
  description = "example public ip"
  vpc_id      = mgc_network_vpc.example.id
}
