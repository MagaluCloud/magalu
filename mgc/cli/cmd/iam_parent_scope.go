package cmd

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

const (
	parentTypeParam = "parent_type"
	parentIDParam   = "parent_id"

	parentOrganization = "organization"
	parentProject      = "project"
)

func applyIamParentScope(sdk *mgcSdk.Sdk, exec core.Executor, parameters core.Parameters, configs core.Configs) {
	schema := exec.ParametersSchema()
	if schema == nil || schema.Properties[parentTypeParam] == nil {
		return
	}
	switch parameters[parentTypeParam] {
	case parentOrganization:
		delete(parameters, parentIDParam)
	case parentProject:
		if id, _ := parameters[parentIDParam].(string); id == "" {
			if resolved := parentProjectID(sdk, configs); resolved != "" {
				parameters[parentIDParam] = resolved
			} else {
				delete(parameters, parentIDParam)
			}
		}
	}
}

func parentProjectID(sdk *mgcSdk.Sdk, configs core.Configs) string {
	if id, _ := configs[config.ProjectKey].(string); id != "" {
		return id
	}
	id, _ := sdk.Auth().CurrentTenantID()
	return id
}
