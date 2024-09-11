data "mgc_kubernetes_version" "cluster_version" {
}

output "cluster_version_output" {
  value = data.mgc_kubernetes_version.cluster_version
}