# Data source to get the latest Kubernetes version
data "mgc_kubernetes_version" "latest" {
}

# Basic Kubernetes cluster
resource "mgc_kubernetes_cluster" "basic_cluster" {
  name        = "smoke-test-basic-cluster"
  version     = data.mgc_kubernetes_version.latest.versions[0].version
  description = "Basic Kubernetes cluster for smoke test"
}

# Full-featured Kubernetes cluster
resource "mgc_kubernetes_cluster" "full_cluster" {
  name                 = "smoke-test-full-cluster"
  version              = data.mgc_kubernetes_version.latest.versions[0].version
  description          = "Full-featured Kubernetes cluster for smoke test"
  enabled_server_group = true
  async_creation       = true
  allowed_cidrs        = ["10.0.0.0/24", "192.168.1.0/24"]
  zone                 = "example-zone-1"
}

# Outputs for verification
output "basic_cluster_id" {
  value = mgc_kubernetes_cluster.basic_cluster.id
}

output "basic_cluster_created_at" {
  value = mgc_kubernetes_cluster.basic_cluster.created_at
}

output "full_cluster_id" {
  value = mgc_kubernetes_cluster.full_cluster.id
}

output "full_cluster_created_at" {
  value = mgc_kubernetes_cluster.full_cluster.created_at
}

data "mgc_kubernetes_flavor" "available_flavors" {
}

# Basic nodepool for the basic cluster
resource "mgc_kubernetes_nodepool" "basic_nodepool" {
  name        = "basic-nodepool"
  cluster_id  = mgc_kubernetes_cluster.basic_cluster.id
  flavor_name = data.mgc_kubernetes_flavor.available_flavors.nodepool[0].name
  replicas    = 1
}

# Full-featured nodepool for the full cluster
resource "mgc_kubernetes_nodepool" "full_nodepool" {
  name         = "full-nodepool"
  cluster_id   = mgc_kubernetes_cluster.full_cluster.id
  flavor_name  = data.mgc_kubernetes_flavor.available_flavors.nodepool[0].name
  replicas     = 2
  min_replicas = 1
  max_replicas = 5
  tags         = ["smoke-test", "full-featured"]

  taints = [
    {
      key    = "dedicated"
      value  = "smoke-test"
      effect = "NoSchedule"
    }
  ]
}
