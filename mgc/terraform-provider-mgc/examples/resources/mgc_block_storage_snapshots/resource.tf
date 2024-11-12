resource "mgc_block_storage_snapshots" "snapshot_example" {
  description = "example of description"
  name        = "exemplo snapshot name"
  snapshot_source_id = mgc_block_storage_snapshots.other_snapshot.id
  type        = "instant"
  volume = {
    id = mgc_block_storage_volumes.example_volume.id
  }
}
