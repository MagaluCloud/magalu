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
