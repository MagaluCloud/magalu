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

type addParameters struct {
	Scopes core.Scopes `json:"scopes" jsonschema:"description=Scopes to be added to the current access token" mgc:"positional"`
}

var addLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("add")
})

var getAdd = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name: "add",
			Description: `Add new scopes to the current access token. Run 'auth scopes list-all'
to see all available scopes to be added`,
			Summary: "Add new scopes to the current access token",
		},
		add,
	)
})

func add(ctx context.Context, params addParameters, _ struct{}) (core.Scopes, error) {
	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	addLogger().Debug("will get all possible scopes for validation")

	allScopes, err := ListAllAvailable(ctx)
	if err != nil {
		addLogger().Warnw("unable to list all possible scopes", "err", err)
		return nil, fmt.Errorf("unable to list all available scopes: %w", err)
	}

	for _, scope := range params.Scopes {
		if !slices.Contains(allScopes, scope) {
			addLogger().Debugw("invalid scope passed as parameter", "scope", scope, "validScopes", allScopes)
			return nil, core.UsageError{Err: fmt.Errorf("invalid scope: %s", scope)}
		}
	}

	addLogger().Debug("will get current scopes")

	currentScopes, err := a.CurrentScopes()
	if err != nil {
		addLogger().Warnw("unable to get current scopes from auth", "err", err)
		return nil, err
	}

	addLogger().Debugw("got current scopes, will concatenate new scopes", "currentScopes", currentScopes, "newScopes", params.Scopes)

	currentScopes.Add(params.Scopes...)

	addLogger().Debugw("will call token exchange with new scopes", "scopes", currentScopes)
	_, err = a.SetScopes(ctx, currentScopes)
	if err != nil {
		addLogger().Warnw("token exchange failed", "scopes", currentScopes, "err", err)
		return nil, err
	}

	currentScopes, err = a.CurrentScopes()
	if err != nil {
		addLogger().Warnw(
			"add operation was successful but unable to get current scopes to check",
			"err", err,
		)
		return currentScopes, nil
	}

	missing := core.Scopes{}
	for _, scope := range params.Scopes {
		if !slices.Contains(currentScopes, scope) {
			missing.Add(scope)
		}
	}

	if len(missing) > 1 {
		return currentScopes, fmt.Errorf("request was successful but resulting scopes are not as requested. Missing %v", missing)
	}

	slices.Sort(currentScopes)

	return currentScopes, nil
}
