resource "mgc_network_subnetpools_book_cidr" "book_subnetpool" {
  cidr = "172.0.0.5/32"
  subnet_pool_id   = "example-subnetpool-id"
}
