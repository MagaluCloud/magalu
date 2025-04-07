MODULES := mgc/cli mgc/core mgc/sdk mgc/spec_manipulator

MGCDIR ?= mgc/cli/
CICD_DIR ?= mgc/spec_manipulator/
DUMP_TREE = mgc/cli/cli-dump-tree.json
OUT_DIR = mgc/cli/docs

OAPIDIR=mgc/sdk/openapi/openapis
SPECS_DIR=specs

# SPECS PRODUCTS
K8S_ORIGIN_URL = https://mke.br-ne-1.jaxyendy.com/docs/openapi-with-snippets.json
K8S_SPEC_FILE = specs/kubernetes.openapi.json
K8S_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/kubernetes.openapi.yaml
K8S_SPEC_CUSTOMIZED_FILE = openapi-customizations/kubernetes.openapi.yaml

CONTAINER_REGISTRY_ORIGIN_URL = https://mcr.br-ne-1.jaxyendy.com/docs/openapi.yaml
CONTAINER_REGISTRY_SPEC_FILE = specs/container-registry.openapi.yaml
CONTAINER_REGISTRY_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/container-registry.openapi.yaml

NETWORK_ORIGIN_URL = https://vpc.prod.se1.br.jaxyendy.com/openapi/openapi-cli.json
NETWORK_SPEC_FILE = specs/network.jaxyendy.openapi.json
NETWORK_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/network.openapi.yaml

BLOCK_STORAGE_ORIGIN_URL = https://block-storage.br-ne-1.jaxyendy.com/v1/openapi.json
BLOCK_STORAGE_SPEC_FILE = specs/block-storage.jaxyendy.openapi.json
BLOCK_STORAGE_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/block-storage.openapi.yaml


DATABASE_ORIGIN_URL = https://dbaas.br-ne-1.jaxyendy.com/v1/cli/openapi.yaml
DATABASE_SPEC_FILE = specs/dbaas.jaxyendy.openapi.yaml
DATABASE_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/dbaas.openapi.yaml


COMPUTE_ORIGIN_URL = https://virtual-machine.br-ne-1.jaxyendy.com/v1/openapi.json
COMPUTE_SPEC_FILE = specs/virtual-machine.jaxyendy.openapi.json
COMPUTE_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/virtual-machine.openapi.yaml
COMPUTE_SPEC_CUSTOMIZED_FILE = openapi-customizations/virtual-machine.openapi.yaml

EVENTS_CONSULT_ORIGIN_URL = https://events-consult-api.br-ne-1.jaxyendy.com/openapi-cli.json
EVENTS_CONSULT_SPEC_FILE = specs/audit.openapi.yaml
EVENTS_CONSULT_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/audit.openapi.yaml
EVENTS_CONSULT_SPEC_CUSTOMIZED_FILE = openapi-customizations/audit.openapi.yaml

# PROFILE_ORIGIN_URL = https://globaldb-public-api.prod.gcp-sa-east1.br.jaxyendy.com/openapi.json
PROFILE_ORIGIN_URL = @specs/globaldb.openapi.yaml
PROFILE_SPEC_FILE = specs/globaldb.openapi.yaml
PROFILE_SPEC_MGC_FILE = mgc/sdk/openapi/openapis/profile.openapi.yaml
