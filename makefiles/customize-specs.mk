include makefiles/variables.mk

customize-k8s:
	@echo "Customizing k8s spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(K8S_SPEC_MGC_FILE) \
	 -o $(K8S_SPEC_MGC_FILE) \
	 -p kubernetes \
	 -r \
	 --remove-params "x-tenant-id"
	@echo "Customizing k8s spec: done"

customize-container-registry:
	@echo "Customizing container registry spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(CONTAINER_REGISTRY_SPEC_MGC_FILE) \
	 -o $(CONTAINER_REGISTRY_SPEC_MGC_FILE) \
	 -p container-registry \
	 -r \
	 --remove-params "x-tenant-id"
	@echo "Customizing container registry spec: done"

customize-network:
	@echo "Customizing network spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(NETWORK_SPEC_MGC_FILE) \
	 -o $(NETWORK_SPEC_MGC_FILE) \
	 -p network \
	 -r \
	 --remove-params "x-tenant-id"
	@echo "Customizing network spec: done"

customize-block-storage:
	@echo "Customizing block storage spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(BLOCK_STORAGE_SPEC_MGC_FILE) \
	 -o $(BLOCK_STORAGE_SPEC_MGC_FILE) \
	 -p block-storage \
	 -r \
	 --configure-security \
	 --remove-params "x-tenant-id"
	@echo "Customizing block storage spec: done"

customize-database:
	@echo "Customizing database spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(DATABASE_SPEC_MGC_FILE) \
	 -o $(DATABASE_SPEC_MGC_FILE) \
	 -p database \
	 -r \
	 --remove-params "x-tenant-id"
	@echo "Customizing database spec: done"

customize-compute:
	@echo "Customizing compute spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(COMPUTE_SPEC_MGC_FILE) \
	 -o $(COMPUTE_SPEC_MGC_FILE) \
	 -p virtual-machine \
	 -r \
	 --configure-security \
	 --remove-params "x-tenant-id"
	@echo "Customizing compute spec: done"

customize-events-consult:
	@echo "Customizing events consult spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(EVENTS_CONSULT_SPEC_MGC_FILE) \
	 -o $(EVENTS_CONSULT_SPEC_MGC_FILE) \
	 -p events-consult \
	 -r \
	 -g \
	 --remove-params "x-tenant-id"
	@echo "Customizing events consult spec: done"

customize-profile:
	@echo "Customizing profile spec..."
	@$(CICD_DIR)cicd spec customize \
	 -i $(PROFILE_SPEC_MGC_FILE) \
	 -o $(PROFILE_SPEC_MGC_FILE) \
	 -p profile \
	 -r \
	 --remove-params "x-tenant-id"
	@echo "Customizing profile spec: done"
