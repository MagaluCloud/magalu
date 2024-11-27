data "mgc_block_storage_volume_types" "types" {
}

output "types" {
value = {
    for id, volume_types in data.mgc_block_storage_volume_types.types.volume_types :
    id => {
      id = volume_types.id
      name = volume_types.name
    }
  }
}