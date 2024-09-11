resource "mgc_block_storage_volume_attachment" "attach_example" {
  block_storage_id = mgc_block_storage_volumes.my_storage.id
  virtual_machine_id = mgc_virtual_machine_instances.my_vm.id
}