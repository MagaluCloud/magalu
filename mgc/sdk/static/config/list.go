package config

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/config"
)

var listLoggerInstance *zap.SugaredLogger

func listLogger() *zap.SugaredLogger {
	if listLoggerInstance == nil {
		listLoggerInstance = logger().Named("list")
	}
	return listLoggerInstance
}

func configListFormatter(exec core.Executor, result core.Result) string {
	// it must be this, no need to check
	resultWithValue, _ := core.ResultAs[core.ResultWithValue](result)
	configMap := resultWithValue.Value().(map[string]any)

	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"Config", "Type", "Description"})

	for name, schema := range configMap {
		schema := schema.(map[string]any)
		writer.AppendRow(table.Row{name, schema["type"], schema["description"]})
	}

	return writer.Render()
}

func newList() core.Executor {
	executor := core.NewStaticExecuteSimple(
		"list",
		"",
		"list all possible configs",
		getAllConfigs,
	)
	return core.NewExecuteFormat(executor, configListFormatter)
}

func getAllConfigs(ctx context.Context) (map[string]*core.Schema, error) {
	root := core.GrouperFromContext(ctx)
	if root == nil {
		return nil, fmt.Errorf("Couldn't get Group from context")
	}

	config := config.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("Couldn't get Config from context")
	}

	configMap, err := config.BuiltInConfigs()
	if err != nil {
		return nil, fmt.Errorf("unable to get built in configs: %w", err)
	}

	_, err = core.VisitAllExecutors(root, []string{}, func(executor core.Executor, path []string) (bool, error) {
		for name, ref := range executor.ConfigsSchema().Properties {
			current := (*core.Schema)(ref.Value)

			if existing, ok := configMap[name]; ok {
				if !reflect.DeepEqual(existing, current) {
					listLogger().Warnw("unhandled diverging config", "config", name, "path", path, "current", current, "existing", existing)
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
