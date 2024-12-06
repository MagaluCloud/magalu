data "mgc_block_storage_volumes" "my-volumes" {
  provider = mgc.sudeste
}

output "my-volumes" {
  value = data.mgc_block_storage_volumes.my-volumes.name
}