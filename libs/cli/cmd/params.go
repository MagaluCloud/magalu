package cmd

import (
	"log"
	"mgc_sdk"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"

	"github.com/profusion/magalu/libs/parser"
)

func SanitizeFlagName(name string) string {
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
		if req == SanitizeFlagName(property) {
			return true
		}
	}

	return false
}

func GetParams(parameters openapi3.Parameters) ([]*parser.Param, []*parser.Param) {
	pathParams := []*parser.Param{}
	headerParams := []*parser.Param{}

	for _, parameterRef := range parameters {
		parameter := parser.Param{
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

func GetRequestBodyParams(action *parser.OpenAPIAction) []*parser.Param {

	requestBodyParams := []*parser.Param{}

	request := action.Request
	if request == nil {
		return requestBodyParams
	}

	content := request.Value.Content.Get("application/json").Schema.Value
	requiredProperties := content.Required

	for _, propertyRef := range content.Properties {
		property := propertyRef.Value

		sanitizedName := SanitizeFlagName(property.Title)
		parameter := parser.Param{
			Type:        property.Type,
			Name:        sanitizedName,
			Required:    isRequiredProperty(requiredProperties, sanitizedName),
			Description: property.Description,
		}

		requestBodyParams = append(requestBodyParams, &parameter)
	}

	return requestBodyParams
}

func AddFlag(cmd *cobra.Command, name string, p mgc_sdk.Parameter) {
	switch p.Schema().Type {
	case "boolean":
		cmd.Flags().Bool(name, false, p.Description())
	case "integer":
		cmd.Flags().Int(name, 0, p.Description())
	case "number":
		cmd.Flags().Float64(name, 0.0, p.Description())
	case "string":
		cmd.Flags().String(name, "", p.Description())
	case "array[boolean]":
		cmd.Flags().BoolSlice(name, []bool{}, p.Description())
	case "array[integer]":
		cmd.Flags().IntSlice(name, []int{}, p.Description())
	case "array[number]":
		log.Printf("number slice not implemented for param %s", p.Name())
	case "array[string]":
		cmd.Flags().StringSlice(name, []string{}, p.Description())
	}

	if p.Required() {
		err := cmd.MarkFlagRequired(name)
		if err != nil {
			log.Printf("Error marking %s flag as required %s", name, err)
		}
	}
}
