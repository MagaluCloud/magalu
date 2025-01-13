package scopes

import (
	"context"
	"fmt"
	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

var addAllLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("add-all")
})

var getAddAll = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "add-all",
			Description: "Add all scopes from all operations to the current access token.",
			Summary:     "Add all scopes to the current access token",
		},
		addAll,
	)
})

func addAll(ctx context.Context) (core.Scopes, error) {
	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	allScopes, err := ListAllAvailable(ctx)
	if err != nil {
		addAllLogger().Warnw("unable to list all possible scopes", "err", err)
		return nil, fmt.Errorf("unable to list all available scopes: %w", err)
	}

	currentScopes, err := a.CurrentScopes()
	if err != nil {
		addAllLogger().Warnw("unable to get current scopes from auth", "err", err)
		return nil, err
	}

	currentScopes.Add(allScopes...)
	_, err = a.SetScopes(ctx, currentScopes)
	if err != nil {
		addAllLogger().Warnw("token exchange failed", "scopes", currentScopes, "err", err)
		return nil, err
	}

	currentScopes, err = a.CurrentScopes()
	if err != nil {
		addAllLogger().Warnw(
			"add-all operation was successful but unable to get current scopes to check",
			"err", err,
		)
		return currentScopes, nil
	}

	missing := core.Scopes{}
	for _, scope := range allScopes {
		if !slices.Contains(currentScopes, scope) {
			missing.Add(scope)
		}
	}

	if len(missing) > 0 {
		return currentScopes, fmt.Errorf("request was successful but resulting scopes are not as expected. Missing: %v", missing)
	}

	slices.Sort(currentScopes)
	return currentScopes, nil
}
