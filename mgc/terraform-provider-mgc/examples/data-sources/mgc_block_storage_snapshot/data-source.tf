data "mgc_block_storage_snapshot" "snapshot" {
  id         = mgc_block_storage_snapshot.my_snapshot.id
}

output "snapshot" {
  value = data.mgc_block_storage_snapshot.name
}
