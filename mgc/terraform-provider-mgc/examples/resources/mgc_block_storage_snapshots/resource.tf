resource "mgc_block_storage_snapshots" "snapshot_example" {
  description = "example of description"
  name        = "exemplo snapshot name"
  volume = {
    id = mgc_block_storage_volumes.example_volume.id
  }
}
