package workspace

import (
	"context"
	"errors"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	exec := core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all available workspaces",
		},
		list,
	)

	return exec
})

func list(ctx context.Context) ([]*profile_manager.Profile, error) {
	m := profile_manager.FromContext(ctx)
	if m == nil {
		return nil, WorkspaceError{Name: "", Err: errors.New("couldn't get ProfileManager from context")}
	}

	return m.List(), nil
}
