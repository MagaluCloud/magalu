package instances

import (
	cli "integration-tests/cli"
	"testing"
	"time"
)

const (
	ssh_key_name = "publio"
	machine_type = "cloud-bs1.xsmall"
	linux_image  = "cloud-fedora-37"
)

var instanceSessionId string

func init() {
	instanceSessionId = cli.GenerateNanoIDWithPrefix("integration-test-instance")
}

func TestInstanceLifecycle(t *testing.T) {
	t.Run("CreateInstance", testCreateInstance)
	time.Sleep(60 * time.Second)
	t.Run("GetInstance", testGetInstance)
	t.Run("ListInstances", testListInstances)
	t.Run("StopInstance", testStopInstance)
	time.Sleep(60 * time.Second)
	t.Run("StartInstance", testStartInstance)
	time.Sleep(60 * time.Second)
	t.Run("RebootInstance", testRebootInstance)
	time.Sleep(60 * time.Second)
	t.Run("DeleteInstance", testDeleteInstance)
}

func testCreateInstance(t *testing.T) {
	t.Log("Creating instance with id: ", instanceSessionId)
	instance, err := createInstance(instanceSessionId, linux_image, machine_type, ssh_key_name)
	if err != nil {
		t.Errorf("Failed to create instance: %v", err)
	}

	if instance.ID == "" {
		t.Errorf("Unexpected instance name: %s", instance.ID)
	}

	t.Log("Deleting instance with id: ", instance.ID)
	err = deleteInstance(instance.ID)
	if err != nil {
		t.Errorf("Failed to delete instance: %v", err)
	}
}

func testGetInstance(t *testing.T) {
	t.Log("Getting instance with id: ", instanceSessionId)
	instance, err := getInstance(instanceSessionId)
	if err != nil {
		t.Errorf("Failed to get instance: %v", err)
	}

	if instance.ID != instanceSessionId {
		t.Errorf("Unexpected instance id: %s", instance.ID)
	}

	t.Log("Deleting instance with id: ", instance.ID)
	err = deleteInstance(instance.ID)
	if err != nil {
		t.Errorf("Failed to delete instance: %v", err)
	}
}

func testListInstances(t *testing.T) {
	t.Log("Listing instances")
	instances, err := listInstances()
	if err != nil {
		t.Errorf("Failed to list instances: %v", err)
	}

	if len(instances) == 0 {
		t.Errorf("No instances found")
	}
}

func testStopInstance(t *testing.T) {
	t.Log("Stopping instance with id: ", instanceSessionId)
	err := stopInstance(instanceSessionId)
	if err != nil {
		t.Errorf("Failed to stop instance: %v", err)
	}
}

func testStartInstance(t *testing.T) {
	t.Log("Starting instance with id: ", instanceSessionId)
	err := startInstance(instanceSessionId)
	if err != nil {
		t.Errorf("Failed to start instance: %v", err)
	}
}

func testRebootInstance(t *testing.T) {
	t.Log("Rebooting instance with id: ", instanceSessionId)
	err := rebootInstance(instanceSessionId)
	if err != nil {
		t.Errorf("Failed to reboot instance: %v", err)
	}
}

func testDeleteInstance(t *testing.T) {
	t.Log("Deleting instance with id: ", instanceSessionId)
	err := deleteInstance(instanceSessionId)
	if err != nil {
		t.Errorf("Failed to delete instance: %v", err)
	}
}
