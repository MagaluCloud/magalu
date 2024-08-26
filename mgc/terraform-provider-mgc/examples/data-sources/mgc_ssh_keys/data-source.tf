data "mgc_ssh_keys" "keys" {
}

output "ssh_keys" {
  value = data.mgc_ssh_keys.keys
}