package provider

import (
	"context"

	"magalu.cloud/core"
	"magalu.cloud/core/config"
)

func getConfigs(ctx context.Context, schema *core.Schema) core.Configs {
	config := config.FromContext(ctx)
	result := core.Configs{}
	for propName, propRef := range schema.Properties {
		prop := (*core.Schema)(propRef.Value)

		if config != nil {
			var value any
			if err := config.Get(propName, &value); err == nil {
				if value != nil || prop.Nullable {
					result[propName] = value
					continue
				}
			}
		}

		if prop.Default != nil {
			result[propName] = prop.Default
		}
	}
	return result
}
