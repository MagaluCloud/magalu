resource "mgc_virtual_machine_instances" "basic_instance" {
  name     = "basic-instance"
  machine_type = {
    name = "cloud-bs1.xsmall"
  }
  image = {
    name = "cloud-ubuntu-22.04 LTS"
  }
  network = {
    associate_public_ip = false # If true, will create a public IP
    delete_public_ip    = false
  }

  ssh_key_name = "ssh_key"
}


resource "mgc_virtual_machine_instances" "basic_instance_with_SG" {
  name = "basic-instance"
  machine_type = {
    name = "cloud-bs1.xsmall"
  }
  image = {
    name = "cloud-ubuntu-22.04 LTS"
  }
  network = {
    associate_public_ip = false # If true, will create a public IP
    delete_public_ip = false
    interface = {
      security_groups = [{ id = "aa622bcb-6861-4251-9cdb-aaadf3" }]
    }
  }
  ssh_key_name = "ssh_key"
}