terraform {
    required_providers {
        magalu = {
            version = "0.1"
            source = "magalucloud/mgc"
        }
    }
}

provider "magalu" {
    # This will be used later on to test the SDK loading functions
    apis = ["virtual-machine@1.60.0"]
}

resource "magalu_virtual-machine_instances" "myvm" {
  name = "my-tf-vm"
  type = "cloud-bs1.xsmall"
  desired_image = "cloud-ubuntu-22.04 LTS"
  key_name = "luizalabs-key"
  availability_zone = "br-ne-1c"
  status = "active"
  allocate_fip = false
}

resource "magalu_block-storage_volume" "myvmvolume" {
    name = "myvmvolume"
    description = "myvmvolumedescription"
    size = 20
    desired_volume_type = "cloud_nvme"
}

resource "magalu_block-storage_volume_attach" "myvm_myvmvolume_attachment" {
    id = magalu_block-storage_volume.myvmvolume.id
    virtual_machine_id = magalu_virtual-machine_instances.myvm.id
}

resource "magalu_dbaas_instances" "mydbaasinstance" {
    name = "mydbaasinstance"
    user = "user"
    password = "passwd" # This should be a variable
    flavor_id = "8bbe8e01-40c8-4d2b-80e8-189debc44b1c" # Should be name instead of ID like the VM instance resource
    datastore_id = "063f3994-b6c2-4c37-96c9-bab8d82d36f7" # Ditto
    volume = {
        size = 10
        type = "CLOUD_HDD"
    }
}
