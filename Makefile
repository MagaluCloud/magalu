include makefiles/cicd.mk
include makefiles/devqa.mk
include makefiles/download-specs.mk
include makefiles/downgrade-specs.mk
include makefiles/merge-specs.mk
include makefiles/customize-specs.mk

build-local:
	@goreleaser build --clean --snapshot --single-target -f internal.yaml

build-go:
	@cd mgc/cli && go build -tags "embed" -o mgc 

download-all:
	@make download-k8s
	@make download-container-registry
	@make download-network
	@make download-block-storage
	@make download-database
	@make download-compute
	@make download-events-consult
	@make download-profile

downgrade-all:
	@make downgrade-k8s
	@make downgrade-container-registry
	@make downgrade-network
	@make downgrade-block-storage
	@make downgrade-database
	@make downgrade-compute
	@make downgrade-events-consult
	@make downgrade-profile

customize-all:
	@make customize-k8s
	@make customize-container-registry
	@make customize-network
	@make customize-block-storage
	@make customize-database
	@make customize-compute
	@make customize-events-consult
	@make customize-profile

merge-all:
	@make merge-kubernetes
	@make merge-compute
	@make merge-events-consult

refresh-all: download-all downgrade-all customize-all merge-all
add-all-specs: downgrade-all customize-all merge-all