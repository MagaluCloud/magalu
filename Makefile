MODULES := mgc/cli mgc/core mgc/sdk mgc/spec_manipulator

MGCDIR ?= mgc/cli/
SPECS_DIR ?= mgc/spec_manipulator/
DUMP_TREE = mgc/cli/cli-dump-tree.json
OUT_DIR = mgc/cli/docs


build-local:
	@goreleaser build --clean --snapshot --single-target -f internal.yaml

# cicd
build-cicd:
	@echo "RUNNING $@"
	cd $(SPECS_DIR) && go build -o cicd
	cd $(MGCDIR) && go build -tags \"embed\" -o mgc

dump-tree: build-cicd
	@echo "generating $(DUMP_TREE)..."
	$(SPECS_DIR)cicd pipeline dumptree -c $(MGCDIR)mgc -o "$(DUMP_TREE)"
	@echo "generating $(DUMP_TREE): done"
	@echo "ENDING $@"

generate-docs: build-cicd
	@echo "generating $(OUT_DIR)..."
	$(SPECS_DIR)cicd pipeline cligendoc -g true -c $(MGCDIR)mgc -d "$(DUMP_TREE)" -o "$(OUT_DIR)" -v "0"
	@echo "generating $(OUT_DIR): done"
	@echo "ENDING $@"


spec-add-vm: build-cicd
	$(SPECS_DIR)cicd spec merge \
	 -p compute \
	 -a /home/gfz/git/magaluCloud/magalu/specs/virtual-machine.jaxyendy.openapi.json \
	 -b /home/gfz/git/magaluCloud/magalu/openapi-customizations/virtual-machine.openapi.yaml\
	 -o /home/gfz/git/magaluCloud/magalu/mgc/sdk/openapi/openapis/virtual-machine.openapi.yaml 

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s /home/gfz/git/magaluCloud/magalu/mgc/sdk/openapi/openapis/virtual-machine.openapi.yaml

spec-add-k8s: build-cicd
	$(SPECS_DIR)cicd spec merge \
	 -p kubernetes \
	 -a /home/gfz/git/magaluCloud/magalu/specs/kubernetes.jaxyendy.openapi.json \
	 -b /home/gfz/git/magaluCloud/magalu/openapi-customizations/kubernetes.openapi.yaml\
	 -o /home/gfz/git/magaluCloud/magalu/mgc/sdk/openapi/openapis/kubernetes.openapi.yaml 

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s /home/gfz/git/magaluCloud/magalu/mgc/sdk/openapi/openapis/kubernetes.openapi.yaml

spec-add-container: build-cicd
	$(SPECS_DIR)cicd spec merge \
	 -p container \
	 -a /home/gfz/git/magaluCloud/magalu/specs/container-registry.openapi.yaml \
	 -o /home/gfz/git/magaluCloud/magalu/mgc/sdk/openapi/openapis/container-registry.openapi.yaml
	 
	$(SPECS_DIR)cicd spec downgrade-unique \
	-s /home/gfz/git/magaluCloud/magalu/mgc/sdk/openapi/openapis/container-registry.openapi.yaml



add-all-specs: spec-add-vm spec-add-k8s spec-add-container 

# specs
download-specs: build-cicd
	@./mgc/spec_manipulator/cicd spec download
	@echo "Now, run 'make prepare-specs' validate and pretify the specs"

prepare-specs: build-cicd
	@./mgc/spec_manipulator/cicd spec prepare

refresh-specs: build-cicd
	@./mgc/spec_manipulator/cicd spec downgrade
	@poetry install
	@poetry run ./scripts/add_all_specs.sh


# Testing targets
test:
	@echo "Running tests for all modules..."
	@for module in $(MODULES); do \
		echo "Testing $$module"; \
		(cd $$module && go test ./...); \
	done

# Code quality targets
vet:
	@echo "Vetting all modules..."
	@for module in $(MODULES); do \
		echo "Vetting $$module"; \
		(cd $$module && go vet ./...); \
	done

lint:
	@echo "Linting all modules..."
	@for module in $(MODULES); do \
		echo "Linting $$module"; \
		(cd $$module && go vet ./...); \
	done

format:
	@echo "Formatting all modules..."
	@for module in $(MODULES); do \
		echo "Formatting $$module"; \
		(cd $$module && gofmt -s -w .); \
	done

# Combined check
check: format vet lint test
