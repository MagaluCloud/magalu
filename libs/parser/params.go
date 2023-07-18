package parser

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func sanitizeFlagName(name string) string {
	sanitizedName := strings.ToLower(name)
	sanitizedName = strings.ReplaceAll(sanitizedName, " ", "-")

	return sanitizedName
}

func isPathParam(parameter *openapi3.Parameter) bool {
	return parameter.In == "path"
}

func isHeatherParam(parameter *openapi3.Parameter) bool {
	return parameter.In == "header"
}

func isRequiredProperty(requriedSet []string, property string) bool {
	for _, req := range requriedSet {
		if req == sanitizeFlagName(property) {
			return true
		}
	}

	return false
}

func joinParameters(base *openapi3.Parameters, merger *openapi3.Parameters) openapi3.Parameters {
	result := *base

	for _, o := range *merger {
		isPresent := false
		for _, p := range *base {
			if p == o {
				isPresent = true
				break
			}
		}

		if !isPresent {
			result = append(result, o)
		}
	}

	return result
}

func getParams(parameters openapi3.Parameters) ([]*OpenAPIParameter, []*OpenAPIParameter) {
	pathParams := []*OpenAPIParameter{}
	headerParams := []*OpenAPIParameter{}

	for _, parameterRef := range parameters {
		parameter := OpenAPIParameter{
			Type:        parameterRef.Value.Schema.Value.Type,
			Name:        parameterRef.Value.Name,
			Required:    parameterRef.Value.Required,
			Description: parameterRef.Value.Description,
		}

		if isPathParam(parameterRef.Value) {
			pathParams = append(pathParams, &parameter)
		}

		if isHeatherParam(parameterRef.Value) {
			headerParams = append(headerParams, &parameter)
		}
	}

	return pathParams, headerParams
}

func getRequestBodyParams(requestBody *openapi3.RequestBodyRef) []*OpenAPIParameter {
	requestBodyParams := []*OpenAPIParameter{}

	if requestBody == nil {
		return requestBodyParams
	}

	content := requestBody.Value.Content.Get("application/json").Schema.Value
	requiredProperties := content.Required

	for _, propertyRef := range content.Properties {
		property := propertyRef.Value

		sanitizedName := sanitizeFlagName(property.Title)
		parameter := OpenAPIParameter{
			Type:        property.Type,
			Name:        sanitizedName,
			Required:    isRequiredProperty(requiredProperties, sanitizedName),
			Description: property.Description,
		}

		requestBodyParams = append(requestBodyParams, &parameter)
	}

	return requestBodyParams
}
