package provider

import (
	"context"
	"fmt"
	"strings"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	mgcUtilsPkg "magalu.cloud/core/utils"
)

func isSubSchema(parent, sub *core.Schema) bool {
	if parent.Type != sub.Type {
		return false
	}

	if len(sub.AllOf) > len(parent.AllOf) || len(sub.AnyOf) > len(parent.AnyOf) || len(sub.OneOf) > len(parent.OneOf) {
		return false
	}
	for i, allOf := range sub.AllOf {
		if !isSubSchema((*core.Schema)(parent.AllOf[i].Value), (*core.Schema)(allOf.Value)) {
			return false
		}
	}
	for i, anyOf := range sub.AnyOf {
		if !isSubSchema((*core.Schema)(parent.AnyOf[i].Value), (*core.Schema)(anyOf.Value)) {
			return false
		}
	}
	for i, oneOf := range sub.OneOf {
		if !isSubSchema((*core.Schema)(parent.OneOf[i].Value), (*core.Schema)(oneOf.Value)) {
			return false
		}
	}

	switch parent.Type {
	case "array":
		// Should never happen
		if parent.Items == nil || sub.Items == nil {
			return false
		}
		return isSubSchema((*mgcSchemaPkg.Schema)(parent.Items.Value), (*mgcSchemaPkg.Schema)(sub.Items.Value))
	case "object":
		for propName, propRef := range sub.Properties {
			parentPropRef, ok := parent.Properties[propName]
			if !ok {
				return false
			}

			prop := (*core.Schema)(propRef.Value)
			parentProp := (*core.Schema)(parentPropRef.Value)

			if !isSubSchema(parentProp, prop) {
				return false
			}
		}
	}
	return true
}

func findResourceSchemaFromList(
	listResultSchema *mgcSchemaPkg.Schema,
	createResultSchema *core.Schema,
) (resourceSchema *core.Schema, resourceJsonPath string, err error) {
	if listResultSchema == nil {
		err = fmt.Errorf("list result schema is nil")
		return
	}

	// If list result is an array, return the Items type
	if listResultSchema.Items != nil {
		resourceSchema = (*mgcSchemaPkg.Schema)(listResultSchema.Items.Value)
		resourceJsonPath = ""
	} else {
		// Assume list result is an object with the array as a sub-property
		for propName, propRef := range listResultSchema.Properties {
			prop := propRef.Value

			if prop.Items == nil {
				continue
			}
			propItems := (*core.Schema)(prop.Items.Value)

			if !isSubSchema(propItems, createResultSchema) {
				continue
			}

			resourceSchema = propItems
			resourceJsonPath = "." + propName
			break
		}
	}

	if resourceSchema == nil {
		err = fmt.Errorf("unable to find resource schema from list result schema")
		return
	}

	return
}

func isListPaginated(listParamsSchema, deleteParamsSchema *core.Schema) bool {
	for listParamName, listParamRef := range listParamsSchema.Properties {
		deleteParamRef, ok := deleteParamsSchema.Properties[listParamName]
		if !ok {
			return true
		}
		listParam := (*mgcSchemaPkg.Schema)(listParamRef.Value)
		deleteParam := (*mgcSchemaPkg.Schema)(deleteParamRef.Value)
		if !mgcSchemaPkg.CheckSimilarJsonSchemas(listParam, deleteParam) {
			continue
		}
	}
	return false
}

func createReadFromList(listExec core.Executor, createResultSchema, deleteParamsSchema *core.Schema) (core.Executor, error) {
	if isListPaginated(listExec.ParametersSchema(), deleteParamsSchema) {
		return nil, fmt.Errorf("list operation cannot be paginated")
	}

	resourceSchema, resourceJsonPath, err := findResourceSchemaFromList(listExec.ResultSchema(), createResultSchema)
	if err != nil {
		return nil, err
	}

	paramsSchema, err := mgcSchemaPkg.SimplifySchema(mgcSchemaPkg.NewAllOfSchema(
		createResultSchema,
		listExec.ParametersSchema(),
	))
	if err != nil {
		return nil, err
	}

	var readExec core.Executor = core.NewSimpleExecutor(core.ExecutorSpec{
		DescriptorSpec: core.DescriptorSpec{
			Name:        "read",
			Version:     listExec.Version(),
			Description: listExec.Description(),
			Summary:     listExec.Summary(),
		},
		ParametersSchema: paramsSchema,
		ConfigsSchema:    listExec.ConfigsSchema(),
		ResultSchema:     resourceSchema,
		Execute: func(readExec core.Executor, ctx context.Context, parameters core.Parameters, configs core.Configs) (core.Result, error) {
			readParamsSchema := createResultSchema // For clarity
			listParamsSchema := listExec.ParametersSchema()
			listParams := make(core.Parameters, len(listParamsSchema.Properties))
			readParams := make(core.Parameters, len(readParamsSchema.Properties))

			for paramName := range listParamsSchema.Properties {
				param, ok := parameters[paramName]
				if !ok {
					return nil, fmt.Errorf("missing parameter: %s", paramName)
				}
				listParams[paramName] = param
			}
			for paramName := range readParamsSchema.Properties {
				param, ok := parameters[paramName]
				if !ok {
					return nil, fmt.Errorf("missing parameter: %s", paramName)
				}
				readParams[paramName] = param
			}

			listResult, err := listExec.Execute(ctx, listParams, configs)
			if err != nil {
				return nil, err
			}
			listResultWithValue, ok := listResult.(core.ResultWithValue)
			if !ok {
				return nil, fmt.Errorf("resource not found")
			}

			listFilters := make([]string, 0, len(readParamsSchema.Properties))
			for propName := range readParamsSchema.Properties {
				listFilters = append(listFilters, fmt.Sprintf("@.%s == $.read_parameters.%s", propName, propName))
			}

			resource, err := mgcUtilsPkg.GetJsonPath(
				fmt.Sprintf("$.list_result%s[?(%s)]", resourceJsonPath, strings.Join(listFilters, " && ")),
				map[string]core.Value{
					"list_result":     listResultWithValue.Value(),
					"list_parameters": listParams,
					"read_parameters": readParams,
				},
			)
			if err != nil {
				return nil, err
			}

			if resource == nil {
				return nil, fmt.Errorf("resource not found")
			}

			resultSource := core.ResultSource{
				Executor:   readExec,
				Context:    ctx,
				Parameters: parameters,
				Configs:    configs,
			}

			if arr, ok := resource.([]any); ok {
				if len(arr) == 0 {
					return nil, fmt.Errorf("resource not found")
				}

				return core.NewSimpleResult(resultSource, resourceSchema, arr[0]), nil
			}

			return core.NewSimpleResult(resultSource, resourceSchema, resource), nil
		},
	})
	return readExec, nil
}
