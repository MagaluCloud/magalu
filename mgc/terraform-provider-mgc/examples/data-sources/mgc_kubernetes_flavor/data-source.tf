data "mgc_kubernetes_flavor" "flavor" {
}

output "flavor_output" {
  value = data.mgc_kubernetes_flavor.flavor
}