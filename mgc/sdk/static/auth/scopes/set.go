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

type setParameters struct {
	Scopes auth.Scopes `json:"scopes" jsonschema:"description=The new scopes to be saved in the current access token" mgc:"positional"`
}

var setLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("set")
})

var getSet = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name: "set",
			Description: `Set the scopes for the current scopes in the access token.
Run 'auth scopes list-all' to see all available scopes`,
			Summary: "Set the scopes for the current scopes in the access token.",
		},
		set,
	)
})

func set(ctx context.Context, params setParameters, _ struct{}) (auth.Scopes, error) {
	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}
	httpClient := mgcHttpPkg.ClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: context did not contain http client")
	}

	setLogger().Debugw("will call token exchange with new scopes", "scopes", params.Scopes)
	_, err := a.SetScopes(ctx, params.Scopes, &httpClient.Client)
	if err != nil {
		addLogger().Warnw("token exchange failed", "scopes", params.Scopes, "err", err)
		return nil, err
	}

	slices.Sort(params.Scopes)

	return params.Scopes, nil
}
