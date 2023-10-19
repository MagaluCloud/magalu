terraform {
    required_providers {
        mgc = {
            version = "0.1"
            source = "magalucloud/mgc"
        }
    }
}

provider "mgc" {
    # This will be used later on to test the SDK loading functions
    apis = ["virtual-machine@1.60.0", "block-storage@1.52.0"]
}

resource "mgc_virtual-machine_instances" "myvm" {
  name = "my-tf-vm-f37-st-rs"
  type = "cloud-bs1.xsmall"
  desired_image = "cloud-fedora-37"
  key_name = "luizalabs-key"
  availability_zone = "br-ne-1c"
  status = "active"
  allocate_fip = false
}

resource "mgc_block-storage_volume" "myvmvolume" {
    name = "myvmvolume"
    description = "myvmvolumedescription"
    size = 1
    desired_volume_type = "cloud_nvme"
}

resource "mgc_block-storage_volume-attach" "myvmvolumeattachment" {
    id = mgc_block-storage_volume.myvmvolume.id
    virtual_machine_id = mgc_virtual-machine_instances.myvm.id
}
