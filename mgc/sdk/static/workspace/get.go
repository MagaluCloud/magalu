package workspace

import (
	"context"
	"errors"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	exec := core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get current workspace.",
		},
		getProfile,
	)

	return core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "template={{.name}}\n"
	})
})

func getProfile(ctx context.Context) (*profile_manager.Profile, error) {
	m := profile_manager.FromContext(ctx)
	if m == nil {
		return nil, WorkspaceError{Name: "", Err: errors.New("couldn't get ProfileManager from context")}
	}

	return m.Current(), nil
}
