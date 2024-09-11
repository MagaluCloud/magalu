resource "mgc_kubernetes_nodepool" "nodepool" {
  name         = "Gandalf"
  cluster_id   = mgc_kubernetes_cluster.cluster_with_nodepool.id
  flavor_name  = kubernetes_flavor.flavors[0].name
  replicas     = 1
  min_replicas = 1
  max_replicas = 5
}
