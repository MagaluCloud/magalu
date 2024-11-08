resource "mgc_network_public_ips_attach" "example" {
  public_ip_id = mgc_network_public_ips.example.id
  interface_id = mgc_network_vpcs_interfaces.primary_interface.id
}
