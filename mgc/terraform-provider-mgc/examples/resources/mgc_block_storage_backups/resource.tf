resource "mgc_block_storage_backups" "backup" {
  name = "backup"
  volume = {
    id = mgc_block_storage_volume.volume.id
  }
}