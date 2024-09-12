# Basic instance
resource "mgc_virtual_machine_instances" "basic_instance" {
  name = "basic-instance-TEST-SMOKE"

  machine_type = {
    name = "cloud-bs1.xsmall"
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

# Instance with Security Group
resource "mgc_virtual_machine_instances" "instance_with_sg" {
  name = "instance-with-sgTEST-SMOKE"

  machine_type = {
    name = "cloud-bs1.small"
  }

  image = {
    name = "cloud-centos-7"
  }

  network = {
    associate_public_ip = true
    delete_public_ip    = true
    interface = {
      security_groups = [
        {
          id = "sg-123456"
        }
      ]
    }
  }

  ssh_key_name = "publio"

}

# Instance with custom VPC
resource "mgc_virtual_machine_instances" "instance_with_vpc" {
  name = "instance-with-vpc-TEST-SMOKE"

  machine_type = {
    name = "cloud-bs1.medium"
  }

  image = {
    name = "cloud-windows-2019"
  }

  network = {
    associate_public_ip = true
    delete_public_ip    = false
    vpc = {
      id   = "vpc-987654"
      name = "my-custom-vpc"
    }
  }

  # Note: SSH key is not used for Windows images
}

# Instance with name as prefix
resource "mgc_virtual_machine_instances" "instance_with_prefix" {
  name           = "prefix-TEST-SMOKE-"
  name_is_prefix = true

  machine_type = {
    name = "cloud-bs1.large"
  }

  image = {
    name = "cloud-debian-11"
  }

  network = {
    associate_public_ip = true
    delete_public_ip    = false
  }

  ssh_key_name = "publio"

}
