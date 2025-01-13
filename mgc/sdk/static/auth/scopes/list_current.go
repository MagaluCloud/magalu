package scopes

import (
	"context"
	"fmt"
	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var getListCurrent = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "list-current",
			Description: "List scopes present in the current access token",
		},
		listCurrent,
	)
})

func listCurrent(ctx context.Context) (core.Scopes, error) {
	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	currentScopes, err := a.CurrentScopes()
	if err != nil {
		return currentScopes, err
	}

	slices.Sort(currentScopes)

	return currentScopes, nil
}
