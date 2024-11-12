resource "mgc_block_storage_volumes" "example_volume" {
  name = "example-volume"
  availability_zone = "br-se1-a"
  size = 10
  type = {
    name = "cloud_nvme"
  }
}