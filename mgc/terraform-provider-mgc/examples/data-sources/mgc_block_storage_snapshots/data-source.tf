data "mgc_block_storage_snapshots" "snapshot" {
  id         = mgc_block_storage_snapshots.my_snapshot.id
}

output "snapshot" {
  value = data.mgc_block_storage_snapshots.name
}
