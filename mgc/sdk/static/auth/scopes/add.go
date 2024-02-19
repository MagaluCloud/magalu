package scopes

import (
	"context"
	"fmt"
	"slices"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"
)

type addParameters struct {
	Scopes auth.Scopes `json:"scopes" jsonschema:"description=Scopes to be added to the current access token" mgc:"positional"`
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

func add(ctx context.Context, params addParameters, _ struct{}) (auth.Scopes, error) {
	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	httpClient := mgcHttpPkg.ClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: context did not contain http client")
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
	_, err = a.SetScopes(ctx, currentScopes, &httpClient.Client)
	if err != nil {
		addLogger().Warnw("token exchange failed", "scopes", currentScopes, "err", err)
		return nil, err
	}

	slices.Sort(currentScopes)

	return currentScopes, nil
}
