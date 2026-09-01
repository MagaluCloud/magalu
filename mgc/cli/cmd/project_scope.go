package cmd

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

func resolveProjectScope(cfg *config.Config) string {
	if cfg == nil {
		return config.ProjectDefault
	}

	var id string
	if err := cfg.Get(config.ProjectKey, &id); err != nil || id == "" {
		return config.ProjectDefault
	}
	return id
}

func withProjectScope(ctx context.Context, sdk *mgcSdk.Sdk, exec core.Executor) context.Context {
	if !exec.DescriptorSpec().ProjectScoped {
		return ctx
	}
	return config.NewProjectScopeContext(ctx, resolveProjectScope(sdk.Config()))
}
