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


add-spec-vm:
	@echo "Adding virtual-machine spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p compute \
	 -a specs/virtual-machine.jaxyendy.openapi.json \
	 -b openapi-customizations/virtual-machine.openapi.yaml \
	 -o mgc/sdk/openapi/openapis/virtual-machine.openapi.yaml 

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/virtual-machine.openapi.yaml
	@echo "Virtual-machine spec added successfully.\n\n"

add-spec-vm:
	@echo "Adding virtual-machine spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p compute \
	 -a specs/virtual-machine.jaxyendy.openapi.json \
	 -b openapi-customizations/virtual-machine.openapi.yaml \
	 -o mgc/sdk/openapi/openapis/virtual-machine.openapi.yaml 

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/virtual-machine.openapi.yaml
	@echo "Virtual-machine spec added successfully.\n\n"

add-spec-k8s:
	@echo "Adding kubernetes spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p kubernetes \
	 -a specs/kubernetes.jaxyendy.openapi.json \
	 -b openapi-customizations/kubernetes.openapi.yaml \
	 -o mgc/sdk/openapi/openapis/kubernetes.openapi.yaml 

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/kubernetes.openapi.yaml
	@echo "Kubernetes spec added successfully.\n\n"

add-spec-container:
	@echo "Adding container-registry spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p container-registry \
	 -a specs/container-registry.openapi.yaml \
	 -o mgc/sdk/openapi/openapis/container-registry.openapi.yaml

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/container-registry.openapi.yaml

	@echo "Container-registry spec added successfully.\n\n"

add-spec-block-storage:
	@echo "Adding block-storage spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p block-storage \
	 -a specs/block-storage.jaxyendy.openapi.json \
	 -o mgc/sdk/openapi/openapis/block-storage.openapi.yaml

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/block-storage.openapi.yaml

	@echo "Block-storage spec added successfully.\n\n"

add-spec-database:
	@echo "Adding database(dbaas) spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p dbaas \
	 -a specs/database.jaxyendy.openapi.json \
	 -o mgc/sdk/openapi/openapis/database.openapi.yaml

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/database.openapi.yaml
	@echo "Database spec added successfully.\n\n"


add-spec-events:
	@echo "Adding events(audit) spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p audit \
	 -g \
	 -a specs/events-consult.openapi.yaml \
	 -b openapi-customizations/audit.openapi.yaml \
	 -o mgc/sdk/openapi/openapis/audit.openapi.yaml

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/audit.openapi.yaml

	@echo "Audit spec added successfully.\n\n"

add-spec-globaldb:
	@echo "Adding globaldb(profile) spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p profile \
	 -a specs/globaldb.openapi.yaml \
	 -b openapi-customizations/profile.openapi.yaml \
	 -o mgc/sdk/openapi/openapis/profile.openapi.yaml

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/profile.openapi.yaml
	@echo "Globaldb spec added successfully.\n\n"


add-spec-network:
	@echo "Adding network spec..."
	$(SPECS_DIR)cicd spec merge \
	 -p network \
	 -a specs/network.jaxyendy.openapi.json \
	 -o mgc/sdk/openapi/openapis/network.openapi.yaml

	@echo "\nDowngrading to 3.0.3...\n"

	$(SPECS_DIR)cicd spec downgrade-unique \
	-s mgc/sdk/openapi/openapis/network.openapi.yaml

	@echo "Network spec added successfully.\n\n"

add-all-specs: build-cicd add-spec-vm add-spec-k8s add-spec-container add-spec-block-storage add-spec-database add-spec-events add-spec-globaldb add-spec-network

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
