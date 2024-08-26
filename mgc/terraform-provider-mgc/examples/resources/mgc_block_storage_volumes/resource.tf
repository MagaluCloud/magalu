resource "mgc_block_storage_volumes" "example_volume" {
  name = "example-volume"
  size = 10
  type = {
    name = "cloud_nvme"
  }
}