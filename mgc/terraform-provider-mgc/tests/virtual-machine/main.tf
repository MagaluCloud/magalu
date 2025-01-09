# Variables for testing updates
variable "instance_name" {
  description = "Name for the test instance"
  type        = string
  default     = "tc1-basic-instance"
}

variable "machine_type" {
  description = "Machine type for the test instance"
  type        = string
  default     = "BV1-1-40"
}

# Test Case 1: Basic VM Instance
resource "mgc_virtual_machine_instances" "tc1_basic_instance" {
  name              = var.instance_name
  availability_zone = "br-ne1-a"
  machine_type      = var.machine_type
  image             = "cloud-ubuntu-24.04 LTS"
  ssh_key_name      = "publio"
  user_data         = base64encode("#!/bin/bash\necho 'Test Case 1: Basic Instance'")

  lifecycle {
    create_before_destroy = true
  }
}

# Test Case 2: VM Instance with Availability Zone
resource "mgc_virtual_machine_instances" "tc2_instance_with_az" {
  name              = "tc2-instance-with-az"
  availability_zone = "br-ne1-a"
  machine_type      = "BV4-8-100"
  image             = "cloud-ubuntu-24.04 LTS"
  ssh_key_name      = "publio"
  user_data         = base64encode("#!/bin/bash\necho 'Test Case 2: AZ Instance'")

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

output "all_vm_instances" {
  description = "All VM Instances for testing"
  value       = data.mgc_virtual_machine_instances.all_instances
}

output "tc1_basic_instance_details" {
  value = {
    created_at         = mgc_virtual_machine_instances.tc1_basic_instance.created_at
    vpc_id             = mgc_virtual_machine_instances.tc1_basic_instance.vpc_id
    network_interfaces = mgc_virtual_machine_instances.tc1_basic_instance.network_interfaces
  }
}

output "tc2_instance_with_az_details" {
  value = {
    created_at         = mgc_virtual_machine_instances.tc2_instance_with_az.created_at
    vpc_id             = mgc_virtual_machine_instances.tc2_instance_with_az.vpc_id
    network_interfaces = mgc_virtual_machine_instances.tc2_instance_with_az.network_interfaces
  }
}

# Additional outputs for test validation
output "vm_instance_details" {
  description = "Virtual Machine Instance Details"
  value = {
    machine_type = {
      id   = data.mgc_virtual_machine_instance.tc1_validation.machine_type_id
      name = data.mgc_virtual_machine_instance.tc1_validation.machine_type_name
      ram  = data.mgc_virtual_machine_instance.tc1_validation.machine_type_ram
      vcpu = data.mgc_virtual_machine_instance.tc1_validation.machine_type_vcpus
      disk = data.mgc_virtual_machine_instance.tc1_validation.machine_type_disk
    }
    image = {
      id       = data.mgc_virtual_machine_instance.tc1_validation.image_id
      name     = data.mgc_virtual_machine_instance.tc1_validation.image_name
      platform = data.mgc_virtual_machine_instance.tc1_validation.image_platform
    }
    network = {
      vpc_id   = data.mgc_virtual_machine_instance.tc1_validation.vpc_id
      vpc_name = data.mgc_virtual_machine_instance.tc1_validation.vpc_name
    }
  }
}

output "tc1_details" {
  value = {
    id                 = mgc_virtual_machine_instances.tc1_basic_instance.id
    name               = mgc_virtual_machine_instances.tc1_basic_instance.name
    created_at         = mgc_virtual_machine_instances.tc1_basic_instance.created_at
    vpc_id             = mgc_virtual_machine_instances.tc1_basic_instance.vpc_id
    machine_type       = mgc_virtual_machine_instances.tc1_basic_instance.machine_type
    network_interfaces = mgc_virtual_machine_instances.tc1_basic_instance.network_interfaces
  }
}

output "tc2_details" {
  value = {
    id                 = mgc_virtual_machine_instances.tc2_instance_with_az.id
    name               = mgc_virtual_machine_instances.tc2_instance_with_az.name
    created_at         = mgc_virtual_machine_instances.tc2_instance_with_az.created_at
    vpc_id             = mgc_virtual_machine_instances.tc2_instance_with_az.vpc_id
    machine_type       = mgc_virtual_machine_instances.tc2_instance_with_az.machine_type
    network_interfaces = mgc_virtual_machine_instances.tc2_instance_with_az.network_interfaces
  }
}

output "vm_instance_validation" {
  description = "Complete validation of VM instance fields"
  value = {
    basic_instance = {
      id                 = data.mgc_virtual_machine_instance.tc1_validation.id
      name               = data.mgc_virtual_machine_instance.tc1_validation.name
      created_at         = data.mgc_virtual_machine_instance.tc1_validation.created_at
      updated_at         = data.mgc_virtual_machine_instance.tc1_validation.updated_at
      image_id          = data.mgc_virtual_machine_instance.tc1_validation.image_id
      image_name        = data.mgc_virtual_machine_instance.tc1_validation.image_name
      image_platform    = data.mgc_virtual_machine_instance.tc1_validation.image_platform
      machine_type_id   = data.mgc_virtual_machine_instance.tc1_validation.machine_type_id
      machine_type_name = data.mgc_virtual_machine_instance.tc1_validation.machine_type_name
      machine_type_disk = data.mgc_virtual_machine_instance.tc1_validation.machine_type_disk
      machine_type_ram  = data.mgc_virtual_machine_instance.tc1_validation.machine_type_ram
      machine_type_vcpus = data.mgc_virtual_machine_instance.tc1_validation.machine_type_vcpus
      vpc_id            = data.mgc_virtual_machine_instance.tc1_validation.vpc_id
      vpc_name          = data.mgc_virtual_machine_instance.tc1_validation.vpc_name
      ssh_key_name      = data.mgc_virtual_machine_instance.tc1_validation.ssh_key_name
      status            = data.mgc_virtual_machine_instance.tc1_validation.status
      state             = data.mgc_virtual_machine_instance.tc1_validation.state
      user_data         = data.mgc_virtual_machine_instance.tc1_validation.user_data
      availability_zone = data.mgc_virtual_machine_instance.tc1_validation.availability_zone
      labels            = data.mgc_virtual_machine_instance.tc1_validation.labels
      error_message     = data.mgc_virtual_machine_instance.tc1_validation.error_message
      error_slug        = data.mgc_virtual_machine_instance.tc1_validation.error_slug
      interfaces        = data.mgc_virtual_machine_instance.tc1_validation.interfaces
    }
  }
}
