package attachment

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

type CreateAttachVolumeParams struct {
	VolumeID         string `json:"id" jsonschema:"description=Block storage volume ID to be attached"`
	VirtualMachineID string `json:"virtual_machine_id" jsonschema:"description=ID of the virtual machine instance to attach the volume"`
}

func newCreate() core.Executor {
	return core.NewStaticExecute(
		"create",
		"",
		"Attach a volume to a virtual machine instance",
		Create,
	)
}

func Create(ctx context.Context, params CreateAttachVolumeParams, cfg core.Configs) (*AttachmentResult, error) {
	exec, err := retrieveExecutor(ctx, []string{"block-storage", "volume", "attach"})
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
		if index == -1 {
			return nil, fmt.Errorf("unable to attach virtual machine %s to volume %s", params.VirtualMachineID, params.VolumeID)
		}

		return &AttachmentResult{resp.ID, resp.Attachments[index].VirtualMachineId}, nil
	}

	return nil, fmt.Errorf("unable to parse command output. %#v", result)
}
