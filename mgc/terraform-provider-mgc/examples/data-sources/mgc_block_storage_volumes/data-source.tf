data "mgc_block_storage_volumes" "my-volume" {
  provider = mgc.sudeste
  volume_id       = mgc_block_storage_volumes.my_volume.id
}

output "my-volume" {
  value = data.mgc_block_storage_volumes.my-volume.name
}