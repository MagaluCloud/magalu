data "mgc_object_storage_bucket" "bucket" {
  name = "bucket-name"
}

output "bucket" {
  value = data.mgc_object_storage_bucket.bucket
}