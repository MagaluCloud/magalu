package machine_types

import (
	"encoding/json"
	"os/exec"

	model "integration-tests/cli/virtual-machines"
)

func listMachineTypes() (model.InstanceTypesResponse, error) {
	cmd := exec.Command("mgc", "vm", "machine-types", "list", "-o", "json", "--raw")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return model.InstanceTypesResponse{}, err
	}

	var instanceTypesResponse model.InstanceTypesResponse

	err = json.Unmarshal(output, &instanceTypesResponse)
	if err != nil {
		return model.InstanceTypesResponse{}, err
	}
	return instanceTypesResponse, nil
}
