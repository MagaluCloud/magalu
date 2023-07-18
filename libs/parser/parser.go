package parser

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/profusion/magalu/libs/functional"
)

func (module *OpenAPIModule) ActionsByTag() map[*openapi3.Tag][]*OpenAPIAction {
	result := make(map[*openapi3.Tag][]*OpenAPIAction)

	for _, action := range module.Actions {
		for _, tag := range action.Tags {
			actionList, isInitialized := result[tag]

			if isInitialized {
				result[tag] = append(actionList, action)
			} else {
				result[tag] = []*OpenAPIAction{action}
			}
		}
	}

	return result
}

func JoinParameters(base openapi3.Parameters, merger openapi3.Parameters) openapi3.Parameters {
	resultMap := make(map[string]*openapi3.ParameterRef)

	for _, o := range base {
		resultMap[o.Value.Name] = o
	}

	for _, o := range merger {
		resultMap[o.Value.Name] = o
	}

	return functional.TransformMap(
		resultMap,
		func(key string, obj *openapi3.ParameterRef) *openapi3.ParameterRef {
			return obj
		},
	)
}

func filterTags(tags openapi3.Tags, include []string) openapi3.Tags {
	result := make(openapi3.Tags, 0)
	for _, tag := range tags {
		if functional.Contains(include, tag.Name) {
			result = append(result, tag)
		}
	}
	return result
}

func CollapsePointer[T any](optional *T, fallback *T) *T {
	if optional != nil {
		return optional
	}

	return fallback
}

func fieldByCaseInsensitiveName(value reflect.Value, fieldName string) reflect.Value {
	lowerFieldName := strings.ToLower(fieldName)
	return value.FieldByNameFunc(func(s string) bool {
		return strings.ToLower(s) == lowerFieldName
	})
}

func getHttpMethodOperation(
	httpMethod HttpMethod,
	pathItem *openapi3.PathItem,
) *openapi3.Operation {
	value := reflect.Indirect(reflect.ValueOf(pathItem))
	field := fieldByCaseInsensitiveName(value, string(httpMethod))

	if !field.IsValid() {
		return nil
	}

	operationPtr := field.Interface().(*openapi3.Operation)
	return operationPtr
}

/* We only accept a single server URL for now, this will be the address used to
 * make all requests, it will probably change since we should only access all
 * endpoints through the gateway, so configuring in Viper would be a better
 * option */
func getServerURL(servers *openapi3.Servers) string {
	if servers == nil || len(*servers) < 1 {
		return ""
	}

	return (*servers)[0].URL
}

func getPathAction(
	pathName string,
	httpMethod HttpMethod,
	operation *openapi3.Operation,
	ctx OpenAPIActionContext,
) *OpenAPIAction {
	return &OpenAPIAction{
		Summary:     operation.Summary + ctx.Summary,
		Description: operation.Description + ctx.Description,
		ServerURL:   getServerURL(operation.Servers) + ctx.ServerURL,
		PathName:    pathName,
		HttpMethod:  httpMethod,
		Tags:        filterTags(ctx.Tags, operation.Tags),
		Deprecated:  operation.Deprecated,
		Parameters:  JoinParameters(ctx.Parameters, operation.Parameters),
		Request:     operation.RequestBody,
		Responses:   operation.Responses,
		Security:    operation.Security,
	}
}

func getPathActions(
	pathName string,
	pathItem *openapi3.PathItem,
	ctx *OpenAPIContext,
) []*OpenAPIAction {
	actionCtx := OpenAPIActionContext{
		ServerURL:            getServerURL(&pathItem.Servers) + ctx.ServerURL,
		Parameters:           pathItem.Parameters,
		Summary:              pathItem.Summary,
		Description:          pathItem.Description,
		Tags:                 ctx.Tags,
		SecurityRequirements: ctx.SecurityRequirements,
	}

	result := make([]*OpenAPIAction, 0)
	for _, method := range AllHttpMethods {
		operation := getHttpMethodOperation(method, pathItem)

		if operation != nil {
			action := getPathAction(pathName, method, operation, actionCtx)
			result = append(result, action)
		}
	}
	return result
}

func getAllActionsInPaths(
	paths openapi3.Paths,
	ctx *OpenAPIContext,
) []*OpenAPIAction {
	result := make([]*OpenAPIAction, 0)

	for key, value := range paths {
		pathActions := getPathActions(key, value, ctx)
		result = append(result, pathActions...)
	}

	return result
}

func LoadOpenAPI(fileInfo *OpenAPIFileInfo) (*OpenAPIModule, error) {
	ctx := context.Background()
	loader := openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(fileInfo.Path)

	if err != nil {
		log.Println("Unable to load OpenAPI document:", fileInfo.Path)
		return nil, err
	}

	/* Define BaseURL for module */
	serverURL := getServerURL(&doc.Servers)

	openAPICtx := OpenAPIContext{
		ServerURL:            serverURL,
		Tags:                 doc.Tags,
		SecurityRequirements: doc.Security,
	}
	actions := getAllActionsInPaths(doc.Paths, &openAPICtx)

	module := &OpenAPIModule{
		Name:                 fileInfo.Name,
		Description:          doc.Info.Description,
		Version:              doc.OpenAPI,
		ServerURL:            serverURL,
		Tags:                 doc.Tags,
		SecurityRequirements: &doc.Security,
		Actions:              actions,
	}

	return module, nil
}
