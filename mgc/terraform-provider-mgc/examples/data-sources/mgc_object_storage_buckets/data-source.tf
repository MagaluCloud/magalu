data "mgc_object_storage_buckets" "buckets" {
}

output "buckets" {
  value = data.mgc_object_storage_buckets.buckets
}