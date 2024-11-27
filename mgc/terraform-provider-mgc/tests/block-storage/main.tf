resource "mgc_block_storage_volumes" "example_volume" {
  name               = "example-volume"
  availability_zones = ["br-ne1-a"]
  size               = 110
  type = {
    name = "nvme"
  }
}

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
}

resource "mgc_block_storage_volume_attachment" "example_attachment" {
  block_storage_id   = mgc_block_storage_volumes.example_volume.id
  virtual_machine_id = mgc_virtual_machine_instances.basic_instance.id
}

resource "mgc_block_storage_snapshots" "snapshot_example" {
  description = "snapshot-example"
  name        = "snapshot-example"
  type        = "instant"
  volume = {
    id = mgc_block_storage_volumes.example_volume.id
  }
  depends_on = [mgc_block_storage_volume_attachment.example_attachment]
}
