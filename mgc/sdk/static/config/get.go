package config

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcConfigPkg "magalu.cloud/core/config"
)

type configGetParams struct {
	Key string `json:"key" validate:"required" jsonschema_description:"Name of the desired config"`
}

func newGet() *core.StaticExecute {
	return core.NewStaticExecute(
		"get",
		"",
		"Gets a specific config value",
		func(ctx context.Context, parameter configGetParams, _ struct{}) (result core.Value, err error) {
			config := mgcConfigPkg.FromContext(ctx)
			if config == nil {
				return nil, fmt.Errorf("unable to retrieve system configuration")
			}
			var out any
			if err := config.Get(parameter.Key, &out); err != nil {
				return nil, err
			}

			return out, nil
		},
	)
}
