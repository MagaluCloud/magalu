package machine_types

import (
	"slices"
	"strings"
	"testing"
)

func TestMachineTypesList(t *testing.T) {
	instanceTypesResponse, err := listMachineTypes()
	if err != nil {
		t.Errorf("Failed to list machine types: %v", err)
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
		if strings.HasPrefix(instanceType.Name, "cloud") && !slices.Contains(intancesTypes, instanceType.Name) {
			t.Errorf("Unexpected instance type: %s", instanceType.Name)
		}
	}
}
