data "mgc_container_repositories" "repository"{
	registry_id = mgc_container_registries.registry.id
}