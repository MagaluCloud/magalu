/*
Executor: versions

# Description

# Retrieve all versions of an object

import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type VersionsParameters struct {
	Dst string `json:"dst"`
}

type VersionsConfigs struct {
	ChunkSize int    `json:"chunkSize,omitempty"`
	Env       string `json:"env,omitempty"`
	Region    string `json:"region,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

type VersionsResultItem struct {
	ETag         string                    `json:"ETag"`
	IsLatest     bool                      `json:"IsLatest"`
	Key          string                    `json:"Key"`
	LastModified string                    `json:"LastModified"`
	Owner        VersionsResultItemOwner   `json:"Owner"`
	Size         int                       `json:"Size"`
	StorageClass string                    `json:"StorageClass"`
	VersionId    string                    `json:"VersionID"`
	XmlName      VersionsResultItemXmlName `json:"XMLName"`
}

type VersionsResultItemOwner struct {
	DisplayName string `json:"DisplayName"`
	Id          string `json:"ID"`
}

type VersionsResultItemXmlName struct {
	Local string `json:"Local"`
	Space string `json:"Space"`
}

type VersionsResult []VersionsResultItem

func Versions(
	client *mgcClient.Client,
	ctx context.Context,
	parameters VersionsParameters,
	configs VersionsConfigs,
) (
	result VersionsResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("Versions", mgcCore.RefPath("/object-storage/objects/versions"), client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[VersionsParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[VersionsConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[VersionsResult](r)
}

// TODO: links
// TODO: related
