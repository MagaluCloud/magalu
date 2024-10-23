resource "mgc_network_vpcs_interfaces" "interface_example" {
    name = "example-interface"
    vpc_id = mgc_network_vpcs.example.id
}