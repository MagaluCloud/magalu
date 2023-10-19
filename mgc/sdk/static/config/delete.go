package config

import (
	"context"
	"fmt"

	"magalu.cloud/core"
	mgcConfigPkg "magalu.cloud/core/config"
)

type configDeleteParams struct {
	Key string `jsonschema_description:"Name of the config to be deleted"`
}

func newDelete() *core.StaticExecute {
	return core.NewStaticExecute(
		"delete",
		"",
		"Deletes a key from config file",
		func(ctx context.Context, parameter configDeleteParams, _ struct{}) (result core.Value, err error) {
			config := mgcConfigPkg.FromContext(ctx)
			if config == nil {
				return nil, fmt.Errorf("unable to retrieve system configuration")
			}
			return nil, config.Delete(parameter.Key)
		},
	)
}
