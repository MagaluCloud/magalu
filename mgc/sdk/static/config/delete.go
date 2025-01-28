package config

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type configDeleteParams struct {
	Key string `json:"key" jsonschema_description:"Name of the config to be deleted" mgc:"positional"`
}

var getDelete = utils.NewLazyLoader[core.Executor](newDelete)

func newDelete() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:    "delete",
			Summary: "Delete/unset a Config value that had been previously set",
			Description: `Delete/unset a Config value that had been previously set. This does not
affect the environment variables`,
		},
		func(ctx context.Context, parameter configDeleteParams, _ struct{}) (result core.Value, err error) {
			config := mgcConfigPkg.FromContext(ctx)
			if config == nil {
				return nil, fmt.Errorf("unable to retrieve system configuration")
			}
			return nil, config.Delete(parameter.Key)
		},
	)
}
