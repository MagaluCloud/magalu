include makefiles/variables.mk

download-k8s:
	@echo "Downloading k8s spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(K8S_ORIGIN_URL) \
	 -d $(K8S_SPEC_FILE)
	@echo "Downloading k8s spec: done"

download-container-registry:
	@echo "Downloading container registry spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(CONTAINER_REGISTRY_ORIGIN_URL) \
	 -d $(CONTAINER_REGISTRY_SPEC_FILE)
	@echo "Downloading container registry spec: done"


download-network:
	@echo "Downloading network spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(NETWORK_ORIGIN_URL) \
	 -d $(NETWORK_SPEC_FILE)
	@echo "Downloading network spec: done"


download-block-storage:
	@echo "Downloading block storage spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(BLOCK_STORAGE_ORIGIN_URL) \
	 -d $(BLOCK_STORAGE_SPEC_FILE)
	@echo "Downloading block storage spec: done"

download-database:
	@echo "Downloading database spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(DATABASE_ORIGIN_URL) \
	 -d $(DATABASE_SPEC_FILE)
	@echo "Downloading database spec: done"

download-compute:
	@echo "Downloading compute spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(COMPUTE_ORIGIN_URL) \
	 -d $(COMPUTE_SPEC_FILE)
	@echo "Downloading compute spec: done"

download-events-consult:
	@echo "Downloading events consult spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(EVENTS_CONSULT_ORIGIN_URL) \
	 -d $(EVENTS_CONSULT_SPEC_FILE)
	@echo "Downloading events consult spec: done"


download-profile:
	@echo "Downloading profile spec..."
	@$(CICD_DIR)cicd spec download \
	 -s $(PROFILE_ORIGIN_URL) \
	 -d $(PROFILE_SPEC_FILE)
	@echo "Downloading profile spec: done"

