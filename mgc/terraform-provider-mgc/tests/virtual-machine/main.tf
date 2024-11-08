resource "mgc_virtual_machine_instances" "basic_instance" {
  name = "basic-instance-test-smoke"

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
  user_data    = "ZWNobyAnSGVsbG8sIFdvcmxkCg=="
}

data "mgc_virtual_machine_instance" "basic_instance_data" {
  id = mgc_virtual_machine_instances.basic_instance.id
}

output "basic_instance_data" {
  value = data.mgc_virtual_machine_instance.basic_instance_data
}
