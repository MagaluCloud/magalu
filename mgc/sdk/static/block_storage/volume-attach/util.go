package attachment

import (
	"context"
	"fmt"

	"magalu.cloud/core"
)

func retrieveExecutor(ctx context.Context, path []string) (core.Executor, error) {
	group := core.GrouperFromContext(ctx)
	if group == nil {
		return nil, fmt.Errorf("unable to retrieve command execution information")
	}

	root := group
	for i, step := range path {
		desc, err := root.GetChildByName(step)
		if err != nil {
			return nil, err
		}

		if i == len(path)-1 {
			if exec, ok := desc.(core.Executor); ok {
				return exec, nil
			}
		} else {
			grouper, ok := desc.(core.Grouper)
			if !ok {
				return nil, fmt.Errorf("unable to convert %s(%d) to command group", step, i)
			}
			root = grouper
		}
	}

	return nil, fmt.Errorf("unable to retrieve command executor")
}
