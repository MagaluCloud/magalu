resource "mgc_virtual_machine_instances" "tc1_basic_instance" {
  name         = "basic-instance-name"
  machine_type = "BV1-1-40"
  image        = "cloud-ubuntu-24.04 LTS"
  ssh_key_name = "your-ssh-key-name"
}

resource "mgc_virtual_machine_instances" "tc2_instance_with_az" {
  name              = "tc2-instance-with-az"
  availability_zone = "br-ne1-a"
  machine_type      = "BV4-8-100"
  image             = "cloud-ubuntu-24.04 LTS"
  ssh_key_name      = "your-ssh-key-name"
}

resource "mgc_virtual_machine_instances" "tc3_instance_with_usardata" {
  name              = "tc3-instance-with-userdata"
  machine_type      = "BV4-8-100"
  image             = "cloud-ubuntu-24.04 LTS"
  ssh_key_name      = "your-ssh-key-name"
  user_data         = base64encode("#!/bin/bash\necho 'Hello, World!'")
}

resource "mgc_virtual_machine_instances" "tc4_instance_with_windows" {
  name              = "tc4-instance-with-windows"
  machine_type      = "BV4-8-100"
  image             = "windows-server-2022"
}
