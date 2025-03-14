build-local:
	goreleaser build --clean --snapshot --single-target -f goreleaser_cli.yaml

generate:
	./mgc/spec_manipulator/build.sh
	./mgc/spec_manipulator/specs downgrade
	poetry install
	poetry run ./scripts/add_all_specs.sh
