include makefiles/variables.mk

merge-kubernetes:
	@echo "Merging kubernetes spec..."
	@$(CICD_DIR)cicd spec merge \
	 -a $(K8S_SPEC_MGC_FILE) \
	 -b $(K8S_SPEC_CUSTOMIZED_FILE) \
	 -o $(K8S_SPEC_MGC_FILE) 
	@echo "Merging kubernetes spec: done"

merge-compute:
	@echo "Merging compute spec..."
	@$(CICD_DIR)cicd spec merge \
	 -a $(COMPUTE_SPEC_MGC_FILE) \
	 -b $(COMPUTE_SPEC_CUSTOMIZED_FILE) \
	 -o $(COMPUTE_SPEC_MGC_FILE) 
	@echo "Merging kubernetes spec: done"

merge-events-consult:
	@echo "Merging events consult spec..."
	@$(CICD_DIR)cicd spec merge \
	 -a $(EVENTS_CONSULT_SPEC_MGC_FILE) \
	 -b $(EVENTS_CONSULT_SPEC_CUSTOMIZED_FILE) \
	 -o $(EVENTS_CONSULT_SPEC_MGC_FILE) \
	@echo "Merging kubernetes spec: done"
