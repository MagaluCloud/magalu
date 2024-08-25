# Variables
BINARY_NAME=mgc
SCRIPT_PATH=scripts/helper_gen.go
DOC_OUTPUT_DIR=docs
CLI_DUMP_TREE_JSON=cli-dump-tree.json

# Default target
all: build generate_docs

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -tags "embed" -o $(BINARY_NAME)

# Generate documentation
generate_docs: build
	@echo "Generating documentation..."
	@go run $(SCRIPT_PATH) \
		--cli ./$(BINARY_NAME) \
		--cli-dump-tree-json $(CLI_DUMP_TREE_JSON) \
		--output-directory $(DOC_OUTPUT_DIR)

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f $(CLI_DUMP_TREE_JSON)
	@rm -rf $(DOC_OUTPUT_DIR)

# Run the script (equivalent to the bash script)
run:
	@echo "Running script..."
	@cd "$$(dirname "$$0")"
	@$(MAKE) build
	@$(MAKE) generate_docs

.PHONY: all build generate_docs clean run