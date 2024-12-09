data "mgc_block_storage_snapshots" "snapshot" {
  provider = mgc.nordeste
}

output "snapshot" {
  value = data.mgc_block_storage_snapshots.name
}
