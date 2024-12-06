data "mgc_block_storage_volume" "my-volume" {
  provider  = mgc.sudeste
  volume_id = mgc_block_storage_volume.my_volume.id
}

output "my-volume" {
  value = data.mgc_block_storage_volume.my-volume.name
}