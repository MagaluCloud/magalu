terraform {
  required_providers {
    mgc = {
      version = "0.1"
      source  = "magalucloud/mgc"
    }
  }
}

provider "mgc" {}

resource "mgc_virtual-machine_instances" "myvm" {
  name = "my-tf-vm"
  machine_type = {
    name = "cloud-bs1.xsmall"
  }
  image = {
    name = "cloud-ubuntu-22.04 LTS"
  }
  ssh_key_name      = "luizalabs-key"
  availability_zone = "br-ne-1c"
}

resource "mgc_block-storage_volumes" "myvmvolume" {
  name = "myvmvolume"
  size = 20
  type = {
    name = "cloud_nvme"
  }
}

resource "mgc_block-storage_volume-attachment" "myvmvolumeattachment" {
  block_storage_id   = mgc_block-storage_volumes.myvmvolume.id
  virtual_machine_id = mgc_virtual-machine_instances.myvm.id
}
