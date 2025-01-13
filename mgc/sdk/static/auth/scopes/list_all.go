package scopes

import (
	"context"
	"fmt"
	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

const (
	mgcSdkDocumentUrl = "http://github.com/MagaluCloud/magalu/mgc/sdk" // url to access Sdk.Group() (executor's root)
)

type listAllParameters struct {
	Targets []string `json:"target,omitempty" jsonschema:"description=If specified\\, only scopes from the target operations will be listed,example=/virtual-machine/instances/create,/block-storage/volume/create" mgc:"positional"`
}

var getListAll = utils.NewLazyLoader(func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "list-all",
			Description: "List all available scopes for all commands",
		},
		listAll,
	)
})

func listAll(ctx context.Context, params listAllParameters, _ struct{}) (core.Scopes, error) {
	if len(params.Targets) > 0 {
		return listAllFromTargets(ctx, params)
	} else {
		return ListAllAvailable(ctx)
	}
}

func listAllFromTargets(ctx context.Context, params listAllParameters) (core.Scopes, error) {
	rootRefResolver := core.RefPathResolverFromContext(ctx)
	if rootRefResolver == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK RefResolver information")
	}

	refResolver := core.NewMultiRefPathResolver()
	refResolver.EmptyDocumentUrl = mgcSdkDocumentUrl
	refResolver.CurrentUrlPlaceholder = mgcSdkDocumentUrl

	err := refResolver.Add(mgcSdkDocumentUrl, rootRefResolver)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve target: %w", err)
	}

	allScopes := core.Scopes{}
	for _, targetRef := range params.Targets {
		target, err := refResolver.Resolve(targetRef, mgcSdkDocumentUrl)
		if err != nil {
			return nil, err
		}

		targetDesc, ok := target.(core.Descriptor)
		if !ok {
			return nil, fmt.Errorf("target was invalid, unable to get DescriptorSpec to fetch Scopes")
		}

		allScopes.Add(targetDesc.Scopes()...)
	}

	return allScopes, nil
}

func ListAllAvailable(ctx context.Context) (core.Scopes, error) {
	root := core.GrouperFromContext(ctx)
	if root == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Grouper information")
	}

	a := auth.FromContext(ctx)
	if a == nil {
		return nil, fmt.Errorf("programming error: context did not contain SDK Auth information")
	}

	builtInScopes := a.BuiltInScopes()
	scopeMap := make(map[core.Scope]struct{}, len(builtInScopes))

	for _, scope := range builtInScopes {
		scopeMap[scope] = struct{}{}
	}

	_, err := core.VisitAllExecutors(root, []string{}, false, func(executor core.Executor, path []string) (bool, error) {
		for _, scope := range executor.Scopes() {
			scopeMap[core.Scope(scope)] = struct{}{}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	// Account for scopes that were returned by the server and are not built-in or in the executors, if possible
	if currentScopes, err := a.CurrentScopes(); err == nil {
		for _, scope := range currentScopes {
			scopeMap[scope] = struct{}{}
		}
	}

	allScopes := make(core.Scopes, 0, len(scopeMap))
	for scope := range scopeMap {
		allScopes.Add(scope)
	}

	slices.Sort(allScopes)

	return allScopes, nil
}
