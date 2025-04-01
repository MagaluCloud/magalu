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

# specs
download-specs: build-cicd
	@./mgc/spec_manipulator/cicd spec download

refresh-specs: build-cicd
	@./mgc/spec_manipulator/cicd spec prepare
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
