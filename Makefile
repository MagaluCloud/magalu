MODULES := mgc/cli mgc/core mgc/sdk mgc/spec_manipulator

build-local:
	@goreleaser build --clean --snapshot --single-target -f goreleaser_cli.yaml

download-specs: --build-spec-manipulator
	@./mgc/spec_manipulator/specs download

refresh-specs: --build-spec-manipulator
	@./mgc/spec_manipulator/specs prepare
	@./mgc/spec_manipulator/specs downgrade
	@poetry install
	@poetry run ./scripts/add_all_specs.sh

--build-spec-manipulator:
	@./mgc/spec_manipulator/build.sh

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
