resource "mgc_object_storage_buckets" "basic_bucket" {
  bucket            = "smoke-test-basic-bucket"
  enable_versioning = true
  bucket_is_prefix  = false
}

resource "mgc_object_storage_buckets" "full_bucket" {
  bucket             = "smoke-test-full-bucket"
  bucket_is_prefix   = true
  enable_versioning  = true
  authenticated_read = true
  public_read        = false
  public_read_write  = false
  private            = true
}

# Output to verify creation
output "basic_bucket_name" {
  value = mgc_object_storage_buckets.basic_bucket.final_name
}

output "full_bucket_name" {
  value = mgc_object_storage_buckets.full_bucket.final_name
}
