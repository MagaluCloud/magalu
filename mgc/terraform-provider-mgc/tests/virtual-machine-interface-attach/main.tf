resource "mgc_network_vpcs" "main_vpc" {
  name = "main-vpc-test-tf-attach-vm"
}

resource "mgc_network_subnetpools" "main_subnetpool" {
  name        = "main-subnetpool"
  description = "Main Subnet Pool"
  type        = "pip"
  cidr        = "172.29.0.0/16"
}

resource "mgc_network_vpcs_subnets" "primary_subnet" {
  cidr_block      = "172.29.1.0/24"
  description     = "Primary Network Subnet"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
  ip_version      = "IPv4"
  name            = "primary-subnet"
  subnetpool_id   = mgc_network_subnetpools.main_subnetpool.id
  vpc_id          = mgc_network_vpcs.main_vpc.id

  depends_on = [ mgc_network_subnetpools.main_subnetpool ]
}

resource "mgc_network_vpcs_interfaces" "primary_interface" {
  name   = "interface-attach-vm"
  vpc_id = mgc_network_vpcs.main_vpc.id

  depends_on = [ mgc_network_vpcs_subnets.primary_subnet ]
}

# Test Case 1: Basic VM Instance
resource "mgc_virtual_machine_instances" "tc1_basic_instance" {
  name = "tc1-basic-instance-attach-vm"
  machine_type = {
    name = "BV1-1-40"
  }

  image = {
    name = "cloud-ubuntu-22.04 LTS"
  }

  network = {
    associate_public_ip = false
    delete_public_ip    = false
  }

  ssh_key_name = "publio"
}

resource "mgc_virtual_machine_interface_attach" "attach_vm" {
  instance_id  = mgc_virtual_machine_instances.tc1_basic_instance.id
  interface_id = mgc_network_vpcs_interfaces.primary_interface.id
}
