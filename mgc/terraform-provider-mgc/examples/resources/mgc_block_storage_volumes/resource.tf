resource "mgc_block_storage_volumes" "example_volume" {
  name              = "example-volume-renamed"
  availability_zone = "br-ne1-a"
  size              = 200
  type              = "cloud_nvme1k"
}