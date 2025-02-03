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

type setParameters struct {
	Scopes core.Scopes `json:"scopes" jsonschema:"description=The new scopes to be saved in the current access token" mgc:"positional"`
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

func set(ctx context.Context, params setParameters, _ struct{}) (core.Scopes, error) {
	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	setLogger().Debugw("will call token exchange with new scopes", "scopes", params.Scopes)
	_, err := a.SetScopes(ctx, params.Scopes)
	if err != nil {
		addLogger().Warnw("token exchange failed", "scopes", params.Scopes, "err", err)
		return nil, err
	}

	currentScopes, err := a.CurrentScopes()
	if err != nil {
		addLogger().Warnw(
			"set operation was successful but unable to get current scopes to check",
			"err", err,
		)
		return params.Scopes, nil
	}

	missing := core.Scopes{}
	for _, scope := range params.Scopes {
		if !slices.Contains(currentScopes, scope) {
			missing.Add(scope)
		}
	}

	if len(missing) > 1 {
		return params.Scopes, fmt.Errorf("request was successful but resulting scopes are not as requested. Missing %v", missing)
	}

	slices.Sort(params.Scopes)

	return params.Scopes, nil
}
