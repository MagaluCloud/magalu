/*
Executor: head

# Description

# Get object metadata

import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "github.com/MagaluCloud/magalu/mgc/core"
	mgcHelpers "github.com/MagaluCloud/magalu/mgc/lib/helpers"
)

type HeadParameters struct {
	Dst        string  `json:"dst"`
	ObjVersion *string `json:"objVersion,omitempty"`
}

type HeadConfigs struct {
	ChunkSize *int    `json:"chunkSize,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
	Workers   *int    `json:"workers,omitempty"`
}

type HeadResult struct {
	AcceptRanges  string `json:"AcceptRanges"`
	ContentLength int    `json:"ContentLength"`
	ContentType   string `json:"ContentType"`
	ETag          string `json:"ETag"`
	LastModified  string `json:"LastModified"`
	StorageClass  string `json:"StorageClass"`
}

func (s *service) Head(
	parameters HeadParameters,
	configs HeadConfigs,
) (
	result HeadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Head", mgcCore.RefPath("/object-storage/objects/head"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[HeadParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[HeadConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[HeadResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) HeadContext(
	ctx context.Context,
	parameters HeadParameters,
	configs HeadConfigs,
) (
	result HeadResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Head", mgcCore.RefPath("/object-storage/objects/head"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[HeadParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[HeadConfigs](configs); err != nil {
		return
	}

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[HeadResult](r)
}

// TODO: links
// TODO: related
