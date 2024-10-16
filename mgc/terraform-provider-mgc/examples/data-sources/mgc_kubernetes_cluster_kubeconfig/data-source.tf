data "mgc_kubernetes_cluster_kubeconfig" "cluster" {
  cluster_id = mgc_kubernetes_cluster.my_cluster.id
}

output "cluster" {
  value = data.mgc_kubernetes_cluster_kubeconfig.cluster
}
