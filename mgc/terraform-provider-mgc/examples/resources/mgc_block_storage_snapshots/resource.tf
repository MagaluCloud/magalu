resource "mgc_block_storage_snapshots" "snapshot_example" {
  name        = "snapshot-example-new"
  description = "Example snapshot description"
  type        = "instant"
  volume_id   = mgc_block_storage_volumes.example_volume.id
}