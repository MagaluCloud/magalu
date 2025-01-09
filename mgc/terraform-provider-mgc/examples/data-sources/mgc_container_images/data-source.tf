data "mgc_container_images" "image"{
	registry_id = mgc_container_registries.registry.id
	repository_name = "alohomora"
}