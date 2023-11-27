package openapi

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/dataloader"
	"magalu.cloud/core/utils"

	"github.com/getkin/kin-openapi/openapi3"
)

// Source -> Module -> Resource -> Operation

// Module

const operationIdsDocKey = "operationIds"

func loadResourceRefs(resource core.Grouper, byPaths map[string]map[string]core.Executor, byId map[string]core.Executor) (bool, error) {
	return resource.VisitChildren(func(child core.Descriptor) (bool, error) {
		if subresource, ok := child.(core.Grouper); ok {
			return loadResourceRefs(subresource, byPaths, byId)
		}

		if exec, ok := child.(core.Executor); ok {
			op, ok := core.ExecutorAs[*operation](exec)
			if !ok {
				return false, fmt.Errorf("expected operation, got %#v", child)
			}

			keyMethods, ok := byPaths[op.key]
			if !ok {
				keyMethods = map[string]core.Executor{}
				byPaths[op.key] = keyMethods
			}

			keyMethods[strings.ToLower(op.method)] = op
			byId[op.operation.OperationID] = op
			return true, nil
		}

		return false, fmt.Errorf(
			"found unknown Descriptor when traversing resource. Expected core.Executor or core.Grouper, got %#v",
			child,
		)
	})
}

func newRefsDocumentLoader(pRoot *core.Grouper) utils.LoadWithError[map[string]any] {
	return utils.NewLazyLoaderWithError(func() (doc map[string]any, err error) {
		if pRoot == nil {
			return nil, fmt.Errorf("missing module root")
		}
		byPaths := map[string]map[string]core.Executor{}
		byId := map[string]core.Executor{}
		_, err = (*pRoot).VisitChildren(func(child core.Descriptor) (bool, error) {
			if resource, ok := child.(core.Grouper); ok {
				return loadResourceRefs(resource, byPaths, byId)
			}

			return false, fmt.Errorf("expected resource to be grouper, got %#v", child)
		})
		if err != nil {
			return
		}
		doc = map[string]any{"paths": byPaths, operationIdsDocKey: byId}
		return doc, nil
	})
}

func newModule(
	indexModule *indexModuleSpec,
	extensionPrefix *string,
	loader dataloader.Loader,
	logger *zap.SugaredLogger,
	refResolver *core.MultiRefPathResolver,
) (m core.Grouper, err error) {
	logger = logger.Named(indexModule.Name)
	loadRefsDocument := newRefsDocumentLoader(&m)
	docResolver := core.NewDocumentRefPathResolver(func() (any, error) { return loadRefsDocument() })
	err = refResolver.Add(indexModule.Url, docResolver)
	if err != nil {
		return
	}

	m = core.NewSimpleGrouper(
		indexModule.DescriptorSpec,
		func() (resources []core.Grouper, err error) {
			ctx := context.Background()
			mData, err := loader.Load(indexModule.Path)
			if err != nil {
				return nil, err
			}

			oapiLoader := openapi3.Loader{Context: ctx, IsExternalRefsAllowed: false}
			doc, err := oapiLoader.LoadFromData(mData)
			if err != nil {
				return
			}

			boundRefResolver := core.NewBoundRefResolver(indexModule.Url, refResolver)
			resources = make([]core.Grouper, 0, len(doc.Tags))

			for _, tag := range doc.Tags {
				resource := newResource(
					tag,
					doc,
					extensionPrefix,
					logger,
					boundRefResolver,
				)

				resources = append(resources, resource)
			}

			return resources, nil
		})

	return m, nil
}
