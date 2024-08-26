data "mgc_kubernetes_nodepool" "nodepool" {
  id         = mgc_kubernetes_nodepool.my_nodepool.id
  cluster_id = mgc_kubernetes_cluster.my_cluster.id
}

output "nodepool" {
  value = data.mgc_kubernetes_nodepool.nodepool
}
