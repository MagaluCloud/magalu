package config

import (
	"context"
	"fmt"
	"reflect"

	"magalu.cloud/core"
)

func NewGroup() *core.StaticGroup {
	return core.NewStaticGroup(
		"config",
		"",
		"config related commands",
		[]core.Descriptor{
			newList(), // cmd: config list
		},
	)
}

func newList() *core.StaticExecute {
	return core.NewStaticExecuteSimple(
		"list",
		"",
		"list all possible configs",
		func(ctx context.Context) (result map[string]any, err error) {
			root := core.GrouperFromContext(ctx)
			if root == nil {
				return nil, fmt.Errorf("Couldn't get Group from context")
			}

			configMap := map[string]any{}
			_, err = core.VisitAllExecutors(root, []string{}, func(executor core.Executor, path []string) (bool, error) {
				for name, ref := range executor.ConfigsSchema().Properties {
					current := (*core.Schema)(ref.Value)

					if existing, ok := configMap[name]; ok {
						if !reflect.DeepEqual(existing, current) {
							fmt.Printf("WARNING: unhandled diverging config at %v %v: %v != %v\n", path, name, existing, current)
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
		},
	)
}
