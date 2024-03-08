package helpers

import (
	"context"
	"fmt"

	mgcCore "magalu.cloud/core"
	mgcClient "magalu.cloud/lib"
)

func ResolveExecutor(
	name string,
	path mgcCore.RefPath,
	client *mgcClient.Client,
) (
	exec mgcCore.Executor,
	err error,
) {
	refResolver := client.Sdk().RefResolver()
	exec, err = mgcCore.ResolveExecutorPath(refResolver, path)
	if err != nil {
		err = fmt.Errorf("%s: unsupported executor: %w", name, err)
		return
	}
	return
}

func PrepareExecutor(
	name string,
	path mgcCore.RefPath,
	client *mgcClient.Client,
	inCtx context.Context,
) (
	exec mgcCore.Executor,
	outCtx context.Context,
	err error,
) {
	exec, err = ResolveExecutor(name, path, client)
	if err != nil {
		return
	}

	if inCtx == nil {
		inCtx = context.Background()
	}
	outCtx = client.Sdk().WrapContext(inCtx)

	return
}
