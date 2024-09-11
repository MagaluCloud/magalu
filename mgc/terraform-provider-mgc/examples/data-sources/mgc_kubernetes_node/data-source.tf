data "mgc_kubernetes_node" "node" {
  nodepool_id = mgc_kubernetes_nodepool.my_nodepool.id
  cluster_id  = mgc_kubernetes_cluster.my_cluster.id
}

output "node" {
  value = data.mgc_kubernetes_node.node
}