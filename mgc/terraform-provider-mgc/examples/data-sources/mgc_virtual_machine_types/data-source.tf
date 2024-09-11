data "mgc_virtual_machine_types" "types" {
}

output "vm_types" {
  value = data.mgc_virtual_machine_types.types
}