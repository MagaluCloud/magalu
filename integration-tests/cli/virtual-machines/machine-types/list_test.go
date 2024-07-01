package cli

import (
	"encoding/json"
	model "integration-tests/cli/virtual-machines"
	"os/exec"
	"slices"
	"testing"
)

func TestMachineTypesList(t *testing.T) {
	cmd := exec.Command("mgc", "vm", "machine-types", "list", "-o", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	var instanceTypesResponse model.InstanceTypesResponse

	err = json.Unmarshal(output, &instanceTypesResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal output: %v", err)
	}

	intancesTypes := []string{
		"cloud-bs1.small",
		"cloud-bs1.xsmall",
		"cloud-hm1.large",
		"cloud-gp1.small",
		"cloud-gp1.medium",
		"cloud-gp1.large",
		"cloud-hm1.medium",
	}

	for _, instanceType := range instanceTypesResponse.InstanceTypes {
		if ! slices.Contains(intancesTypes, instanceType.Name) {
			t.Errorf("Unexpected instance type: %s", instanceType.Name)
		}
	}
}

