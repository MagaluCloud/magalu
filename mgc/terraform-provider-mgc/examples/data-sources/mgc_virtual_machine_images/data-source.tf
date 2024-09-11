data "mgc_virtual_machine_images" "images" {
}

output "vm_images" {
  value = data.mgc_virtual_machine_images.images
}