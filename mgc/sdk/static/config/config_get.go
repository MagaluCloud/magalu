package config

import (
	"context"
	"fmt"

	"magalu.cloud/core"
)

type configGetParams struct {
	Key string `json:"key" validate:"required" jsonschema_description:"Name of the desired config"`
}

func newConfigGet() *core.StaticExecute {
	return core.NewStaticExecute(
		"get",
		"",
		"Gets a specific config value",
		func(ctx context.Context, parameter configGetParams, _ struct{}) (result core.Value, err error) {
			config := core.ConfigFromContext(ctx)
			if config == nil {
				return nil, fmt.Errorf("unable to retrieve system configuration")
			}

			value := config.Get(parameter.Key)
			return value, nil
		},
	)
}
