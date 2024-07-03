package instances

import (
	"encoding/json"
	"os/exec"

	model "integration-tests/cli/virtual-machines"
)

func createInstance(name, image, instanceType, ssh string) (model.CreateNewInstanceResponse, error) {
	exec := exec.Command("mgc",
		"vm",
		"instances",
		"create",
		"--name", name,
		"--image.name", image,
		"--machine-type.name", instanceType,
		"--ssh-key-name", ssh,
		"-o", "json",
		"--raw",
		"--network.associate-public-ip", "false",
		"--cli.wait-termination")

	output, err := exec.CombinedOutput()
	if err != nil {
		return model.CreateNewInstanceResponse{}, err
	}

	var instanceResponse model.CreateNewInstanceResponse

	err = json.Unmarshal(output, &instanceResponse)
	if err != nil {
		return model.CreateNewInstanceResponse{}, err
	}

	return instanceResponse, nil
}

func deleteInstance(id string) error {
	exec := exec.Command("mgc", "vm", "instances", "delete", "--id", id, "--no-confirm", "--raw", "-o", "json")
	return exec.Run()
}

func getInstance(id string) (model.GetVMInstanceResponse, error) {
	exec := exec.Command("mgc", "vm", "instances", "get", "--id", id, "-o", "json", "--raw")
	output, err := exec.CombinedOutput()
	if err != nil {
		return model.GetVMInstanceResponse{}, err
	}

	var instance model.GetVMInstanceResponse

	err = json.Unmarshal(output, &instance)
	if err != nil {
		return model.GetVMInstanceResponse{}, err
	}

	return instance, nil
}

func listInstances() ([]model.GetVMInstanceResponse, error) {
	exec := exec.Command("mgc", "vm", "instances", "list", "-o", "json", "--raw")
	output, err := exec.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var instances []model.GetVMInstanceResponse

	err = json.Unmarshal(output, &instances)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func stopInstance(id string) error {
	exec := exec.Command("mgc", "vm", "instances", "stop", "--id", id, "--raw", "-o", "json")
	return exec.Run()
}

func startInstance(id string) error {
	exec := exec.Command("mgc", "vm", "instances", "start", "--id", id, "--raw", "-o", "json")
	return exec.Run()
}

func rebootInstance(id string) error {
	exec := exec.Command("mgc", "vm", "instances", "reboot", "--id", id, "--raw", "-o", "json")
	return exec.Run()
}
