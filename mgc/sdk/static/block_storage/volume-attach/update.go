package attachment

import (
	"context"

	"magalu.cloud/core"
)

type UpdateAttachVolumeParams struct {
	VolumeID         string `json:"id" jsonschema:"description=Block storage volume ID to be attached"`
	VirtualMachineID string `json:"virtual_machine_id" jsonschema:"description=ID of the virtual machine instance to attach the volume"`
}

func newUpdate() core.Executor {
	return core.NewStaticExecute(
		"update",
		"",
		"Update a block storage volume attachment",
		update,
	)
}

func update(ctx context.Context, params UpdateAttachVolumeParams, cfg core.Configs) (core.Result, error) {
	// No update available
	return nil, nil
}
