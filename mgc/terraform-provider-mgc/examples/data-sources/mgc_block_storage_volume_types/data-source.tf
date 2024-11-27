data "mgc_block_storage_volume_types" "types" {
}

output "types" {
  value = data.mgc_block_storage_volume_types.iops
}
