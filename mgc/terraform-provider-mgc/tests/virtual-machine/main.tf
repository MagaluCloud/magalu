# Test Case 1: Basic VM Instance
resource "mgc_virtual_machine_instances" "tc1_basic_instance" {
  name              = "tc1-basic-instance-name"
  availability_zone = "br-ne1-a"

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
  user_data    = base64encode("#!/bin/bash\necho 'Test Case 1: Basic Instance'")

  lifecycle {
    create_before_destroy = true
  }
}

# Test Case 2: VM Instance with Availability Zone
resource "mgc_virtual_machine_instances" "tc2_instance_with_az" {
  name              = "tc2-instance-with-az"
  availability_zone = "br-ne1-a"

  machine_type = {
    name = "BV4-8-100"
  }

  image = {
    name = "cloud-ubuntu-22.04 LTS"
  }

  network = {
    associate_public_ip = false
    delete_public_ip    = false
  }

  ssh_key_name = "publio"
  user_data    = base64encode("#!/bin/bash\necho 'Test Case 2: AZ Instance'")

  depends_on = [mgc_virtual_machine_instances.tc1_basic_instance]
}

# Data Sources for Validation
data "mgc_virtual_machine_instance" "tc1_validation" {
  id = mgc_virtual_machine_instances.tc1_basic_instance.id
}

data "mgc_virtual_machine_instance" "tc2_validation" {
  id = mgc_virtual_machine_instances.tc2_instance_with_az.id
}

# List Resources for Testing
data "mgc_virtual_machine_instances" "all_instances" {}
data "mgc_virtual_machine_images" "available_images" {}
data "mgc_virtual_machine_types" "available_types" {}

# Test Outputs
output "test_case_1_validation" {
  description = "Validation output for basic instance test case"
  value = {
    instance_id   = data.mgc_virtual_machine_instance.tc1_validation.id
    instance_name = data.mgc_virtual_machine_instance.tc1_validation.name
    status        = data.mgc_virtual_machine_instance.tc1_validation.status
    az            = data.mgc_virtual_machine_instance.tc1_validation.availability_zone
    instace_type  = data.mgc_virtual_machine_instance.tc1_validation.machine_type_id
  }
}

output "test_case_2_validation" {
  description = "Validation output for AZ instance test case"
  value = {
    instance_id   = data.mgc_virtual_machine_instance.tc2_validation.id
    instance_name = data.mgc_virtual_machine_instance.tc2_validation.name
    status        = data.mgc_virtual_machine_instance.tc2_validation.status
    az            = data.mgc_virtual_machine_instance.tc2_validation.availability_zone
    instance_type = data.mgc_virtual_machine_instance.tc2_validation.machine_type_id
  }
}
