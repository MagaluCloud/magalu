package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"magalu.cloud/core"
)

type ConfigListResult struct {
	ConfigMap map[string]*core.Schema
}

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

func visitAllExecutors(child core.Descriptor, path []string, visitExecutor func(executor core.Executor, path []string) (bool, error)) (bool, error) {
	if executor, ok := child.(core.Executor); ok {
		return visitExecutor(executor, path)
	} else if group, ok := child.(core.Grouper); ok {
		return group.VisitChildren(func(child core.Descriptor) (run bool, err error) {
			path = append(path, child.Name())
			run, err = visitAllExecutors(child, path, visitExecutor)
			path = path[:len(path)-1]

			return run, err
		})
	} else {
		return false, fmt.Errorf("child %v not group/executor", child)
	}
}

func newList() *core.StaticExecute {
	return core.NewStaticExecuteSimple(
		"list",
		"",
		"list all possible configs",
		func(ctx context.Context) (result *ConfigListResult, err error) {
			root := core.GrouperFromContext(ctx)
			if root == nil {
				return nil, fmt.Errorf("Couldn't get Group from context")
			}

			configMap := map[string]*core.Schema{}
			_, err = visitAllExecutors(root, []string{}, func(executor core.Executor, path []string) (bool, error) {
				for name, ref := range executor.ConfigsSchema().Properties {
					current := (*core.Schema)(ref.Value)

					if existing, ok := configMap[name]; ok {
						if !reflect.DeepEqual(existing, current) {
							fmt.Println("WARNING: unhandled diverging config at " + strings.Join(path, ".") + "." + name)
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

			return &ConfigListResult{ConfigMap: configMap}, nil
		},
	)
}
