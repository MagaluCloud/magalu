package config

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/config/common"
	"go.uber.org/zap"
)

var getSchemaLogger = utils.NewLazyLoader[*zap.SugaredLogger](func() *zap.SugaredLogger {
	return logger().Named("get-schema")
})

type getSchemaParams struct {
	Key string `json:"key" mgc:"positional"`
}

var getGetSchema = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:    "get-schema",
			Summary: "Get the JSON Schema for the specified Config",
			Description: `Get the JSON Schema for the specified Config. The Schema has
information about the accepted values for the Config, constraints, type, description, etc.`,
		},
		getSchema,
	)
})

func getSchema(ctx context.Context, params getSchemaParams, _ struct{}) (*core.Schema, error) {
	allSchemas, err := common.ListAllConfigSchemas(ctx)
	if err != nil {
		return nil, err
	}

	schema, ok := allSchemas[params.Key]
	if !ok {
		getSchemaLogger().Debugw("no schema found for key, returning nil", "key", params.Key)
		return nil, nil
	}

	return schema, nil
}
