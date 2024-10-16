data "mgc_kubernetes_cluster" "cluster" {
  id = mgc_kubernetes_cluster.my_cluster.id
}

output "cluster" {
  value = data.mgc_kubernetes_cluster.cluster
}
