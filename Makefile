MODULES := mgc/cli mgc/core mgc/sdk mgc/spec_manipulator

MGCDIR ?= mgc/cli/
CICD_DIR ?= mgc/spec_manipulator/
SPECS_DIR ?= specs/
DUMP_TREE = mgc/cli/cli-dump-tree.json
OUT_DIR = mgc/cli/docs
OAPIDIR=mgc/sdk/openapi/openapis

build-local:
	@goreleaser build --clean --snapshot --single-target -f internal.yaml

# cicd
build-cicd:
	@cd $(CICD_DIR) && go build -o cicd

build-cli:
	@cd $(MGCDIR) && go build -tags \"embed\" -o mgc


dump-tree:
	@cd $(CICD_DIR) && go build -o cicd
	@cd $(MGCDIR) && go build -tags \"embed\" -o mgc
	@echo "generating $(DUMP_TREE)..."
	$(CICD_DIR)cicd pipeline dumptree -c $(MGCDIR)mgc -o "$(DUMP_TREE)"
	@echo "generating $(DUMP_TREE): done"
	@echo "ENDING $@"

generate-docs:
	@cd $(CICD_DIR) && go build -o cicd
	@cd $(MGCDIR) && go build -tags \"embed\" -o mgc
	@echo "generating $(OUT_DIR)..."
	$(CICD_DIR)cicd pipeline cligendoc -g true -c $(MGCDIR)mgc -d "$(DUMP_TREE)" -o "$(OUT_DIR)" -v "0"
	@echo "generating $(OUT_DIR): done"
	@$(CICD_DIR)cicd pipeline gen-docs-magalu $(OUT_DIR)
	@echo "ENDING $@"

oapi-index-gen:
	@cd $(CICD_DIR) && go build -o cicd
	@cd $(MGCDIR) && go build -tags \"embed\" -o mgc
	$(CICD_DIR)cicd pipeline oapi-index $(OAPIDIR)
# specs
download-specs:
	@cd $(CICD_DIR) && go build -o cicd
	@./mgc/spec_manipulator/cicd specs download -d $(SPECS_DIR)
	@echo "\nNow, run 'make prepare-specs' to validate and prettify the specs"

prepare-specs:
	@cd $(CICD_DIR) && go build -o cicd
	@./mgc/spec_manipulator/cicd specs prepare -d $(SPECS_DIR)
	@echo "\nNow, run 'make downgrade-specs' to downgrade the specs"

downgrade-specs:
	@cd $(CICD_DIR) && go build -o cicd
	@./mgc/spec_manipulator/cicd specs downgrade -d $(SPECS_DIR)
	@echo "\nNow, run 'make refresh-specs' to finally, refresh the specs"

refresh-specs:
	@cd $(CICD_DIR) && go build -o cicd
	@poetry install
	@poetry run ./scripts/add_all_specs.sh
	@$(CICD_DIR)cicd pipeline oapi-index $(OAPIDIR)


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
