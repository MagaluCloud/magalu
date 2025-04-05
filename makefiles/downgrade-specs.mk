include makefiles/variables.mk

downgrade-k8s:
	@echo "Downgrading k8s spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(K8S_SPEC_FILE) \
	 -o $(K8S_SPEC_MGC_FILE)
	@echo "Downgrading k8s spec: done"

downgrade-container-registry:
	@echo "Downgrading container registry spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(CONTAINER_REGISTRY_SPEC_FILE) \
	 -o $(CONTAINER_REGISTRY_SPEC_MGC_FILE)
	@echo "Downgrading container registry spec: done"

downgrade-network:
	@echo "Downgrading network spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(NETWORK_SPEC_FILE) \
	 -o $(NETWORK_SPEC_MGC_FILE)
	@echo "Downgrading network spec: done"

downgrade-block-storage:
	@echo "Downgrading block storage spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(BLOCK_STORAGE_SPEC_FILE) \
	 -o $(BLOCK_STORAGE_SPEC_MGC_FILE)
	@echo "Downgrading block storage spec: done"

downgrade-database:
	@echo "Downgrading database spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(DATABASE_SPEC_FILE) \
	 -o $(DATABASE_SPEC_MGC_FILE)
	@echo "Downgrading database spec: done"

downgrade-compute:
	@echo "Downgrading compute spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(COMPUTE_SPEC_FILE) \
	 -o $(COMPUTE_SPEC_MGC_FILE)
	@echo "Downgrading compute spec: done"

downgrade-events-consult:
	@echo "Downgrading events consult spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(EVENTS_CONSULT_SPEC_FILE) \
	 -o $(EVENTS_CONSULT_SPEC_MGC_FILE)
	@echo "Downgrading events consult spec: done"

downgrade-profile:
	@echo "Downgrading profile spec..."
	@$(CICD_DIR)cicd spec downgrade \
	 -i $(PROFILE_SPEC_FILE) \
	 -o $(PROFILE_SPEC_MGC_FILE)
	@echo "Downgrading profile spec: done"
