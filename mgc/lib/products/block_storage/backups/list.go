/*
Executor: list

# Summary

List all backups.

# Description

Retrieve a list of Backups for the currently authenticated tenant.

#### Notes
  - Use the **expand** argument to obtain additional details about the
    Volume used to create each Backup.

Version: v1

import "magalu.cloud/lib/products/block_storage/backups"
*/
package backups

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type ListParameters struct {
	Limit  *int                  `json:"_limit,omitempty"`
	Offset *int                  `json:"_offset,omitempty"`
	Sort   *string               `json:"_sort,omitempty"`
	Expand *ListParametersExpand `json:"expand,omitempty"`
}

type ListParametersExpand []string

type ListConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type ListResult struct {
	Backups ListResultBackups `json:"backups"`
}

type ListResultBackupsItem struct {
	CreatedAt    string                             `json:"created_at"`
	Description  *string                            `json:"description,omitempty"`
	Error        *ListResultBackupsItemError        `json:"error,omitempty"`
	Id           string                             `json:"id"`
	Name         string                             `json:"name"`
	Size         int                                `json:"size"`
	SourceBackup *ListResultBackupsItemSourceBackup `json:"source_backup,omitempty"`
	State        string                             `json:"state"`
	Status       string                             `json:"status"`
	Type         string                             `json:"type"`
	UpdatedAt    string                             `json:"updated_at"`
	Volume       *ListResultBackupsItemVolume       `json:"volume"`
}

type ListResultBackupsItemError struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type ListResultBackupsItemSourceBackup struct {
	Id string `json:"id"`
}

// any of: ListResultBackupsItemVolume
type ListResultBackupsItemVolume struct {
	CreatedAt  string                                `json:"created_at"`
	Id         string                                `json:"id"`
	Name       string                                `json:"name"`
	Size       int                                   `json:"size"`
	State      string                                `json:"state"`
	Status     string                                `json:"status"`
	UpdatedAt  string                                `json:"updated_at"`
	VolumeType ListResultBackupsItemVolumeVolumeType `json:"volume_type"`
}

type ListResultBackupsItemVolumeVolumeType struct {
	Id string `json:"id"`
}

type ListResultBackups []ListResultBackupsItem

func (s *service) List(
	parameters ListParameters,
	configs ListConfigs,
) (
	result ListResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("List", mgcCore.RefPath("/block-storage/backups/list"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[ListParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[ListConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[ListResult](r)
}

// TODO: links
// TODO: related
