package merger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

type SpecMerger struct {
}

type MergeOptions struct {
	DowngradeToVersion string
	Servers            []*v3.Server
}

func NewSpecMerger() *SpecMerger {
	return &SpecMerger{}
}

func (m *SpecMerger) MergeSpecs(specAPath, specBPath, outputPath string, options *MergeOptions) error {
	file := filepath.Join(specAPath)
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de especificação A: %v", err)
	}

	docA, err := libopenapi.NewDocument(fileBytes)
	if err != nil {
		return fmt.Errorf("erro ao carregar a especificação A: %v", err)
	}

	var docB libopenapi.Document

	if specBPath != "" {
		file = filepath.Join(specBPath)
		fileBytes, err = os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo de especificação B: %v", err)
		}

		docB, err = libopenapi.NewDocument(fileBytes)
		if err != nil {
			return fmt.Errorf("erro ao carregar a especificação B: %v", err)
		}
	}

	mergedSpec, err := m.performMerge(docA, docB, options)
	if err != nil {
		return fmt.Errorf("erro ao mesclar especificações: %v", err)
	}

	err = m.saveSpec(mergedSpec, outputPath)
	if err != nil {
		return fmt.Errorf("erro ao salvar a especificação mesclada: %v", err)
	}

	return nil
}

func (m *SpecMerger) performMerge(specA, specB libopenapi.Document, options *MergeOptions) (libopenapi.Document, error) {
	mergedSpec := specA

	mergedSpecA, err := mergedSpec.BuildV3Model()
	if err != nil {
		return mergedSpec, fmt.Errorf("erro ao construir modelo V3 para spec A: %v", err)
	}

	var specModelB *libopenapi.DocumentModel[v3.Document]
	if specB != nil {
		var errs []error
		specModelB, errs = specB.BuildV3Model()
		if len(errs) > 0 {
			return mergedSpec, fmt.Errorf("erro ao construir modelo V3 para spec B: %v", errs)
		}
	}

	if options != nil && options.DowngradeToVersion != "" {
		mergedSpecA.Model.Version = options.DowngradeToVersion
	}

	if options != nil && options.Servers != nil && len(options.Servers) > 0 {
		mergedSpecA.Model.Servers = options.Servers
	}

	if err := m.mergePaths(mergedSpecA, specModelB); err != nil {
		return mergedSpec, fmt.Errorf("erro ao mesclar caminhos: %v", err)
	}

	if err := m.mergeSchemas(mergedSpecA, specModelB); err != nil {
		return mergedSpec, fmt.Errorf("erro ao mesclar esquemas: %v", err)
	}

	if err := m.mergeTags(mergedSpecA, specModelB); err != nil {
		return mergedSpec, fmt.Errorf("erro ao mesclar tags: %v", err)
	}

	return mergedSpec, nil
}

func (m *SpecMerger) mergePaths(mergedSpecA *libopenapi.DocumentModel[v3.Document], specModelB *libopenapi.DocumentModel[v3.Document]) error {
	if specModelB == nil || specModelB.Model.Paths == nil {
		return nil
	}

	for path := specModelB.Model.Paths.PathItems.Oldest(); path != nil; path = path.Next() {
		if pathItemA, isPathPresent := mergedSpecA.Model.Paths.PathItems.Get(path.Key); isPathPresent {
			operationsA := pathItemA.GetOperations()
			operationsB := path.Value.GetOperations()

			for opB := operationsB.Oldest(); opB != nil; opB = opB.Next() {
				if operationA, isOpPresent := operationsA.Get(opB.Key); isOpPresent {
					m.mergeOperation(operationA, opB.Value)
				}
			}
		} else {
			mergedSpecA.Model.Paths.PathItems.Set(path.Key, path.Value)
		}

		m.addPathTagsToMainTags(mergedSpecA, path.Key)
	}

	return nil
}

func (m *SpecMerger) mergeOperation(operationA, operationB *v3.Operation) {
	if operationB.Tags != nil {
		operationA.Tags = operationB.Tags
	}

	if operationB.Summary != "" {
		operationA.Summary = operationB.Summary
	}

	if operationB.Description != "" {
		operationA.Description = operationB.Description
	}

	if operationB.ExternalDocs != nil {
		operationA.ExternalDocs = operationB.ExternalDocs
	}

	if operationB.OperationId != "" {
		operationA.OperationId = operationB.OperationId
	}

	if operationB.RequestBody != nil {
		operationA.RequestBody = operationB.RequestBody
	}

	if operationB.Responses != nil {
		for cB := orderedmap.First(operationB.Responses.Codes); cB != nil; cB = cB.Next() {
			if cA, isCAPresent := operationA.Responses.Codes.Get(cB.Key()); isCAPresent {
				ccB := cB.Value()

				if ccB.Content != nil {
					cA.Content = ccB.Content
				}

				if ccB.Description != "" {
					cA.Description = ccB.Description
				}

				if ccB.Extensions != nil {
					cA.Extensions = ccB.Extensions
				}

				if ccB.Headers != nil {
					cA.Headers = ccB.Headers
				}

				if ccB.Links != nil {
					cA.Links = ccB.Links
				}
			} else {
				operationA.Responses.Codes.Set(cB.Key(), cB.Value())
			}
		}
	}

	if operationB.Parameters != nil {
		for _, pB := range operationB.Parameters {
			for _, pA := range operationA.Parameters {
				if pA.Name == pB.Name {
					m.mergeParameter(pA, pB)
				}
			}
		}
	}

	if operationB.Deprecated != nil {
		operationA.Deprecated = operationB.Deprecated
	}

	if operationB.Callbacks != nil {
		operationA.Callbacks = operationB.Callbacks
	}

	if operationB.Security != nil {
		operationA.Security = operationB.Security
	}

	if operationB.Servers != nil {
		operationA.Servers = operationB.Servers
	}

	if operationB.Extensions != nil {
		operationA.Extensions = operationB.Extensions
	}
}

func (m *SpecMerger) mergeParameter(paramA, paramB *v3.Parameter) {
	if paramB.In != "" {
		paramA.In = paramB.In
	}

	if paramB.Description != "" {
		paramA.Description = paramB.Description
	}

	if paramB.Style != "" {
		paramA.Style = paramB.Style
	}

	if paramB.Required != nil {
		paramA.Required = paramB.Required
	}

	if paramB.Explode != nil {
		paramA.Explode = paramB.Explode
	}

	if paramB.Schema != nil {
		paramA.Schema = paramB.Schema
	}

	if paramB.Example != nil {
		paramA.Example = paramB.Example
	}

	if paramB.Examples != nil {
		paramA.Examples = paramB.Examples
	}

	if paramB.Content != nil {
		paramA.Content = paramB.Content
	}

	paramA.Deprecated = paramB.Deprecated
	paramA.AllowEmptyValue = paramB.AllowEmptyValue
	paramA.AllowReserved = paramB.AllowReserved
}

func (m *SpecMerger) addPathTagsToMainTags(mergedSpecA *libopenapi.DocumentModel[v3.Document], pathKey string) {
	if patchedItem, isOk := mergedSpecA.Model.Paths.PathItems.Get(pathKey); isOk {
		httpMethods := []*v3.Operation{
			patchedItem.Get,
			patchedItem.Put,
			patchedItem.Post,
			patchedItem.Delete,
			patchedItem.Options,
			patchedItem.Head,
			patchedItem.Patch,
			patchedItem.Trace,
		}

		for _, method := range httpMethods {
			if method != nil {
				for _, tag := range method.Tags {
					m.addTagIfNotExists(mergedSpecA, tag)
				}
			}
		}
	}
}

func (m *SpecMerger) addTagIfNotExists(mergedSpecA *libopenapi.DocumentModel[v3.Document], tagName string) {
	for _, existingTag := range mergedSpecA.Model.Tags {
		if existingTag.Name == tagName {
			return
		}
	}

	mergedSpecA.Model.Tags = append(mergedSpecA.Model.Tags, &base.Tag{
		Name:        tagName,
		Description: tagName,
	})
}

func (m *SpecMerger) mergeSchemas(mergedSpecA *libopenapi.DocumentModel[v3.Document], specModelB *libopenapi.DocumentModel[v3.Document]) error {
	if specModelB == nil || specModelB.Model.Components == nil {
		return nil
	}

	if mergedSpecA.Model.Components == nil {
		mergedSpecA.Model.Components = &v3.Components{
			Schemas: orderedmap.New[string, *base.SchemaProxy](),
		}
	}

	for schema := specModelB.Model.Components.Schemas.Oldest(); schema != nil; schema = schema.Next() {
		if schemaItemA, isSchemaPresent := mergedSpecA.Model.Components.Schemas.Get(schema.Key); isSchemaPresent {
			*schemaItemA = *schema.Value
		} else {
			mergedSpecA.Model.Components.Schemas.Set(schema.Key, schema.Value)
		}
	}

	return nil
}

func (m *SpecMerger) mergeTags(mergedSpecA *libopenapi.DocumentModel[v3.Document], specModelB *libopenapi.DocumentModel[v3.Document]) error {
	if specModelB == nil || specModelB.Model.Tags == nil {
		return nil
	}

	for _, tags := range specModelB.Model.Tags {
		tagExistsInA := false

		for _, tagMerge := range mergedSpecA.Model.Tags {
			if tags.Name == tagMerge.Name {
				*tagMerge = *tags
				tagExistsInA = true
				break
			}
		}

		if !tagExistsInA {
			mergedSpecA.Model.Tags = append(mergedSpecA.Model.Tags, tags)
		}
	}

	return nil
}

func (m *SpecMerger) saveSpec(spec libopenapi.Document, filename string) error {
	model, errs := spec.BuildV3Model()
	if errs != nil {
		return fmt.Errorf("erro ao gerar modelo do spec mesclado: %v", errs)
	}

	byteFile, err := model.Model.RenderJSON("  ")
	if err != nil {
		return fmt.Errorf("erro ao renderizar spec mesclado: %v", err)
	}

	err = os.WriteFile(filename, byteFile, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo: %v", err)
	}

	return nil
}
