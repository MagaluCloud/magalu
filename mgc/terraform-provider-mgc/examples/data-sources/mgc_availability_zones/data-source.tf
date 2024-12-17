data "mgc_availability_zones" "availability_zones" {
}

output "availability_zones" {
  value = data.mgc_availability_zones.availability_zones.regions
}
