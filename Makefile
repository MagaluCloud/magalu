.ONESHELL:
# TODO: Develop additional commands: release features, etc., and upon completion, enhance the README.

build-spec:
	python3 ./scripts/oapi_index_gen.py "--embed=mgc/sdk/openapi/embed_loader.go" mgc/cli/openapis
	cd mgc/cli; echo "Building...."; \
	go build -tags "embed" -o mgc

build:
	cd mgc/cli; echo "Building...."; \
	go build -tags "embed" -o mgc

test:
	find . -name go.mod -execdir go test ./... \;
