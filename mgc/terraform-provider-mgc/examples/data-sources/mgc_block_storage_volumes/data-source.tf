data "mgc_block_storage_volume" "my-volume" {
  provider = mgc.sudeste
}

output "my-volume" {
  value = data.mgc_block_storage_volume.my-volume.name
}