package config

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"magalu.cloud/core"
	mgcConfigPkg "magalu.cloud/core/config"
	"magalu.cloud/core/utils"
)

type configSetParams struct {
	Key   string `jsonschema_description:"Name of the desired config"`
	Value any    `jsonschema_description:"New flag value"`
}

var getSet = utils.NewLazyLoader[core.Executor](newSet)

func newSet() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Set a specific Config value in the configuration file",
		},
		func(ctx context.Context, parameter configSetParams, _ struct{}) (core.Value, error) {
			config := mgcConfigPkg.FromContext(ctx)
			if config == nil {
				return nil, fmt.Errorf("unable to retrieve system configuration")
			}

			allConfigs, err := getAllConfigs(ctx)
			if err != nil {
				return nil, fmt.Errorf("error when getting possible configs: %w", err)
			}

			s, ok := allConfigs[parameter.Key]
			if !ok {
				return nil, fmt.Errorf("no config %s found", parameter.Key)
			}

			if err := s.VisitJSON(parameter.Value, openapi3.MultiErrors()); err != nil {
				return nil, err
			}

			if err := config.Set(parameter.Key, parameter.Value); err != nil {
				return nil, err
			}

			return parameter.Value, nil

		},
	)
}
