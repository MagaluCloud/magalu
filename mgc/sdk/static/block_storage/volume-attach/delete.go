package attachment

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

type DeleteAttachVolumeParams struct {
	VolumeID         string `json:"id" jsonschema:"description=Block storage volume ID to be detached"`
	VirtualMachineID string `json:"virtual_machine_id" jsonschema:"description=ID of the virtual machine instance to detach the volume from"`
}

func newDelete() core.Executor {
	exec := core.NewStaticExecute(
		"delete",
		"",
		"Detach a volume from a virtual machine instance",
		delete,
	)

	// TODO: Make Confirmable?
	return exec
}

func delete(ctx context.Context, params DeleteAttachVolumeParams, cfg core.Configs) (core.Result, error) {
	exec, err := retrieveExecutor(ctx, []string{"block-storage", "volume", "detach"})
	if err != nil {
		return nil, err
	}

	paramsMap, err := utils.DecodeNewValue[map[string]any](params)
	if err != nil {
		return nil, err
	}

	result, err := exec.Execute(ctx, *paramsMap, cfg)
	if err != nil {
		return nil, err
	}

	if r, ok := result.(core.ResultWithValue); ok {
		resp, err := utils.DecodeNewValue[VolumeResponse](r.Value())
		if err != nil {
			return nil, err
		}

		index := slices.IndexFunc(resp.Attachments, func(attachment VolumeAttachmentResponse) bool {
			return attachment.VirtualMachineId == params.VirtualMachineID
		})
		if index != -1 {
			return nil, fmt.Errorf("unable to detach virtual machine %s to volume %s", params.VirtualMachineID, params.VolumeID)
		}

		return nil, nil
	}

	return nil, fmt.Errorf("unable to parse command output. %#v", result)
}
