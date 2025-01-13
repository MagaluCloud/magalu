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

type removeParameters struct {
	Scopes core.Scopes `json:"scopes" jsonschema:"description=Scopes to be removed from the current access token" mgc:"positional"`
}

var removeLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("remove")
})

var getRemove = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name: "remove",
			Description: `Remove scopes from the current scopes in the access token.
Run 'auth scopes list-current' to see current scopes`,
			Summary: "Remove scopes from the current scopes in the access token.",
		},
		remove,
	)
})

func remove(ctx context.Context, params removeParameters, _ struct{}) (core.Scopes, error) {
	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	builtInScopes := a.BuiltInScopes()
	for _, scopeToRemove := range params.Scopes {
		if slices.Contains(builtInScopes, scopeToRemove) {
			return nil, core.UsageError{Err: fmt.Errorf("unable to remove built-in scope %q", scopeToRemove)}
		}
	}

	removeLogger().Debug("will get current scopes")

	currentScopes, err := a.CurrentScopes()
	if err != nil {
		removeLogger().Warnw("unable to get current scopes from auth", "err", err)
		return nil, err
	}

	removeLogger().Debugw("got current scopes, will remove scopes from it", "currentScopes", currentScopes, "toBeRemoved", params.Scopes)

	currentScopes.Remove(params.Scopes...)

	removeLogger().Debugw("will call token exchange with new scopes", "scopes", currentScopes)
	_, err = a.SetScopes(ctx, currentScopes)
	if err != nil {
		addLogger().Warnw("token exchange failed", "scopes", currentScopes, "err", err)
		return nil, err
	}

	currentScopes, err = a.CurrentScopes()
	if err != nil {
		addLogger().Warnw(
			"remove operation was successful but unable to get current scopes to check",
			"err", err,
		)
		return currentScopes, nil
	}

	leftovers := core.Scopes{}
	for _, scope := range params.Scopes {
		if slices.Contains(currentScopes, scope) {
			leftovers.Add(scope)
		}
	}

	if len(leftovers) > 1 {
		return currentScopes, fmt.Errorf("request was successful but resulting scopes are not as requested. Leftovers %v", leftovers)
	}

	slices.Sort(currentScopes)

	return currentScopes, nil
}
