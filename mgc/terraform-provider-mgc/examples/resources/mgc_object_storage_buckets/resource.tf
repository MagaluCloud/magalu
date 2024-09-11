resource "mgc_object_storage_buckets" "my-bucket" {
  bucket = "bucket-name"
  enable_versioning = true
}