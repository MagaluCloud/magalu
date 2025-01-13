package common

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	mgcUtilsPkg "github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

var listAllSchemasLogger = mgcUtilsPkg.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("list-all-schemas")
})

func ListAllConfigSchemas(ctx context.Context) (map[string]*core.Schema, error) {
	root := core.GrouperFromContext(ctx)
	if root == nil {
		return nil, fmt.Errorf("programming error: couldn't get Group from context")
	}

	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: couldn't get Config from context")
	}

	configMap, err := config.BuiltInConfigs()
	if err != nil {
		return nil, fmt.Errorf("unable to get built-in configs: %w", err)
	}

	_, err = core.VisitAllExecutors(root, []string{}, false, func(executor core.Executor, path []string) (bool, error) {
		for name, ref := range executor.ConfigsSchema().Properties {
			current := (*core.Schema)(ref.Value)

			if name == "region" {
				var copyOfCurrent []interface{}
				for _, enum := range current.Enum {
					if enum != "global" {
						copyOfCurrent = append(copyOfCurrent, enum)
					}
				}
				current.Enum = copyOfCurrent
			}

			if existing, ok := configMap[name]; ok {
				if err := mgcSchemaPkg.CompareJsonSchemas(existing, current); err != nil {
					listAllSchemasLogger().Warnw("unhandled diverging config", "config", name, "path", path, "current", current, "existing", existing, "error", err)
				}

				continue
			}
			configMap[name] = current
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return configMap, nil
}
