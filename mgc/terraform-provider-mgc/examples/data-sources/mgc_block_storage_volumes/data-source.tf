data "mgc_block_storage_volumes" "my-volume" {
  provider = mgc.sudeste
}

output "my-volume" {
  value = data.mgc_block_storage_volumes.my-volume.name
}