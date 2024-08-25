resource "mgc_kubernetes_cluster" "cluster" {
  name                 = "my_cluster"
  version              = mgc_kubernetes_version.versions[0].version
  enabled_server_group = false
  description          = "Cluster Example"
}