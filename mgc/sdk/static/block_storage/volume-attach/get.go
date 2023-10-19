package attachment

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

type GetAttachVolumeParams struct {
	VolumeID         string `json:"id" jsonschema:"description=Block storage volume ID"`
	VirtualMachineID string `json:"virtual_machine_id" jsonschema:"description=Instance ID of the virtual machine to which the volume is attached"`
}

func newGet() core.Executor {
	return core.NewStaticExecute(
		"get",
		"",
		"Check if a volume is attached to a virtual machine instance",
		get,
	)
}

func get(ctx context.Context, params GetAttachVolumeParams, cfg core.Configs) (*AttachmentResult, error) {
	exec, err := retrieveExecutor(ctx, []string{"block-storage", "volume", "get"})
	if err != nil {
		return nil, err
	}

	paramsMap, err := utils.DecodeNewValue[map[string]any](params)
	if err != nil {
		return nil, err
	}

	p := map[string]any{}
	for k := range exec.ParametersSchema().Properties {
		p[k] = (*paramsMap)[k]
	}

	result, err := exec.Execute(ctx, p, cfg)
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
		if index == -1 {
			return nil, fmt.Errorf("virtual machine %s not attached to volume %s", params.VirtualMachineID, params.VolumeID)
		}

		return &AttachmentResult{resp.ID, resp.Attachments[index].VirtualMachineId}, nil
	}

	return nil, fmt.Errorf("unable to parse command output. %#v", result)
}
