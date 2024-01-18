package openapi

import (
	"fmt"
	"strings"

	"slices"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"

	"github.com/getkin/kin-openapi/openapi3"
)

type confirmation struct {
	Message string `json:"message"`
}

// Source -> Module -> Resource -> Operation

// Resource

func getServers(p *openapi3.PathItem, op *openapi3.Operation) openapi3.Servers {
	var servers openapi3.Servers
	if op.Servers != nil && len(*op.Servers) > 0 {
		servers = *op.Servers
	}
	if servers == nil && len(p.Servers) > 0 {
		servers = p.Servers
	}

	return servers
}

func collectOperationsv2(
	name string,
	doc *openapi3.T,
	extensionPrefix *string,
	logger *zap.SugaredLogger,
) *operationTable {
	descs := make([]*operationDesc, 0)
	for key, path := range doc.Paths {
		pathOps := map[string]*openapi3.Operation{
			"get":    path.Get,
			"post":   path.Post,
			"put":    path.Put,
			"patch":  path.Patch,
			"delete": path.Delete,
		}

		for method, op := range pathOps {
			if op == nil {
				continue
			}

			descs = append(descs, &operationDesc{path, op, method, key})
		}
	}

	return newOperationTable(name, descs)
}

func collectOperations(
	tag *openapi3.Tag,
	doc *openapi3.T,
	extensionPrefix *string,
	logger *zap.SugaredLogger,
) *operationTable {
	descs := []*operationDesc{}

	for key, path := range doc.Paths {
		pathOps := map[string]*openapi3.Operation{
			"get":    path.Get,
			"post":   path.Post,
			"put":    path.Put,
			"patch":  path.Patch,
			"delete": path.Delete,
		}

		for method, op := range pathOps {
			if op == nil {
				continue
			}

			if !slices.Contains(op.Tags, tag.Name) {
				continue
			}

			descs = append(descs, &operationDesc{path, op, method, key})
		}
	}

	return newOperationTable(tag.Name, descs)
}

func collectResourceChildren(
	descriptionPrefix string,
	table *operationTable,
	doc *openapi3.T,
	extensionPrefix *string,
	logger *zap.SugaredLogger,
	refResolver *core.BoundRefPathResolver,
) (children []core.Descriptor, err error) {
	children = []core.Descriptor{}
	childrenByName := map[string]core.Descriptor{}

	for _, opTableEntry := range table.childOperations {
		desc := opTableEntry.desc

		opName := getNameExtension(extensionPrefix, desc.op.Extensions, opTableEntry.key)
		servers := getServers(desc.path, desc.op)
		if servers == nil {
			servers = doc.Servers
		}

		outputFlag, _ := getExtensionString(extensionPrefix, "output-flag", desc.op.Extensions, "")
		method := strings.ToUpper(desc.method)

		var operation core.Executor = newOperation(
			opName,
			desc,
			doc.Info.Version,
			method,
			extensionPrefix,
			servers,
			logger,
			outputFlag,
			refResolver,
		)

		isDelete := method == "DELETE"
		cExt, ok := getExtensionObject(extensionPrefix, "confirmable", desc.op.Extensions, nil)

		if (ok && cExt != nil) || isDelete {
			cExec, err := wrapInConfirmableExecutor(cExt, isDelete, operation)
			if err != nil {
				return children, err
			}
			operation = cExec
		}

		if wtExt, ok := getExtensionObject(extensionPrefix, "wait-termination", desc.op.Extensions, nil); ok && wtExt != nil {
			if tExec, err := wrapInTerminatorExecutor(operation, wtExt); err == nil {
				operation = tExec
			} else {
				return children, err
			}
		}

		children = append(children, operation)
		childrenByName[opName] = operation
	}

	for _, childTable := range table.childTables {
		subResource := newSubResource(
			descriptionPrefix,
			childTable,
			doc,
			extensionPrefix,
			logger,
			refResolver,
		)
		children = append(children, subResource)
		childrenByName[childTable.name] = subResource
	}

	return children, nil
}

func newResourcev2(
	table *operationTable,
	doc *openapi3.T,
	extensionPrefix *string,
	logger *zap.SugaredLogger,
	refResolver *core.BoundRefPathResolver,
) *core.SimpleGrouper[core.Descriptor] {
	logger = logger.Named(table.name)
	description := table.name
	for _, tag := range doc.Tags {
		if tag.Name == table.name {
			description = getDescriptionExtension(extensionPrefix, tag.Extensions, tag.Description)
			break
		}
	}
	return core.NewSimpleGrouper[core.Descriptor](
		core.DescriptorSpec{
			Name:        table.name,
			Description: description,
			Version:     doc.Info.Version,
			IsInternal:  false,
		},
		func() ([]core.Descriptor, error) {
			return collectResourceChildren(description, table, doc, extensionPrefix, logger, refResolver)
		},
	)
}

func newResource(
	tag *openapi3.Tag,
	doc *openapi3.T,
	extensionPrefix *string,
	logger *zap.SugaredLogger,
	refResolver *core.BoundRefPathResolver,
) *core.SimpleGrouper[core.Descriptor] {
	logger = logger.Named(tag.Name)
	name := getNameExtension(extensionPrefix, tag.Extensions, tag.Name)
	description := getDescriptionExtension(extensionPrefix, tag.Extensions, tag.Description)
	return core.NewSimpleGrouper[core.Descriptor](
		core.DescriptorSpec{
			Name:        name,
			Description: description,
			Version:     doc.Info.Version,
			IsInternal:  getHiddenExtension(extensionPrefix, tag.Extensions),
		},
		func() ([]core.Descriptor, error) {
			opTable := collectOperations(tag, doc, extensionPrefix, logger)

			return collectResourceChildren(description, opTable, doc, extensionPrefix, logger, refResolver)
		},
	)
}

func newSubResource(
	descriptionPrefix string,
	table *operationTable,
	doc *openapi3.T,
	extensionPrefix *string,
	logger *zap.SugaredLogger,
	refResolver *core.BoundRefPathResolver,
) *core.SimpleGrouper[core.Descriptor] {
	logger = logger.Named(table.name)
	return core.NewSimpleGrouper(
		core.DescriptorSpec{
			Name:        table.name,
			Version:     doc.Info.Version,
			Description: fmt.Sprintf("%s | %s", descriptionPrefix, table.name),
			Summary:     table.name,
		},
		func() ([]core.Descriptor, error) {
			return collectResourceChildren(descriptionPrefix, table, doc, extensionPrefix, logger, refResolver)
		},
	)
}

func wrapInConfirmableExecutor(cExt map[string]any, isDelete bool, exec core.Executor) (core.ConfirmableExecutor, error) {
	c := &confirmation{}

	if cExt != nil {
		if err := utils.DecodeValue(cExt, c); err != nil {
			return nil, fmt.Errorf("error decoding confirmable extension: %w", err)
		}
	}

	return core.NewConfirmableExecutor(exec, core.ConfirmPromptWithTemplate(c.Message)), nil
}
