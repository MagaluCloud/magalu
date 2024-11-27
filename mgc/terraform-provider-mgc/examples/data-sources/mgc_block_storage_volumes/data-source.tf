data "mgc_block_storage_volumes" "volume" {
  id         = mgc_block_storage_volumes.my_volume.id
}

output "volume" {
  value = data.mgc_block_storage_volumes.name
}
