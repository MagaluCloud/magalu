package openapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"magalu.cloud/core"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

// Source -> Module -> Resource -> Operation

// Operation

type Operation struct {
	key           string
	method        string
	path          *openapi3.PathItem
	operation     *openapi3.Operation
	doc           *openapi3.T
	paramsSchema  *core.Schema
	configsSchema *core.Schema
	// TODO: configsMapping map[string]...
	extensionPrefix *string
	servers         openapi3.Servers
}

// BEGIN: Descriptor interface:

var openAPIPathArgRegex = regexp.MustCompile("[{](?P<name>[^}]+)[}]")

func getActionName(httpMethod string, pathName string) string {
	name := []string{string(httpMethod)}
	hasArgs := false

	for _, pathEntry := range strings.Split(pathName, "/") {
		match := openAPIPathArgRegex.FindStringSubmatch(pathEntry)
		for i, substr := range match {
			if openAPIPathArgRegex.SubexpNames()[i] == "name" {
				name = append(name, substr)
				hasArgs = true
			}
		}
		if len(match) == 0 && hasArgs {
			name = append(name, pathEntry)
		}
	}

	return strings.Join(name, "-")
}

func (o *Operation) Name() string {
	name := getNameExtension(o.extensionPrefix, o.operation.Extensions, "")
	if name == "" {
		name = getActionName(o.method, o.key)
	}
	return name
}

func (o *Operation) Version() string {
	return ""
}

func (o *Operation) Description() string {
	return getDescriptionExtension(o.extensionPrefix, o.operation.Extensions, o.operation.Description)
}

// END: Descriptor interface

// BEGIN: Executor interface:

func addParameters(schema *core.Schema, parameters openapi3.Parameters, extensionPrefix *string, locations []string) {
	for _, ref := range parameters {
		parameter := ref.Value

		if !slices.Contains(locations, parameter.In) {
			continue
		}

		paramSchemaRef := parameter.Schema
		paramSchema := paramSchemaRef.Value

		name := getNameExtension(extensionPrefix, parameter.Extensions, parameter.Name)

		desc := getDescriptionExtension(extensionPrefix, parameter.Extensions, parameter.Description)
		if desc == "" {
			desc = getDescriptionExtension(extensionPrefix, paramSchema.Extensions, paramSchema.Description)
		}

		if desc != "" && paramSchema.Description != desc {
			// copy, never modify parameter stuff
			newSchema := *paramSchema
			newSchema.Description = desc
			paramSchema = &newSchema

			newSchemaRef := *paramSchemaRef
			newSchemaRef.Value = paramSchema
			paramSchemaRef = &newSchemaRef
		}

		schema.Properties[name] = paramSchemaRef

		if parameter.Required && !slices.Contains(schema.Required, name) {
			schema.Required = append(schema.Required, name)
		}
	}
}

func addRequestBodyParameters(schema *core.Schema, rbr *openapi3.RequestBodyRef, extensionPrefix *string) {
	if rbr == nil {
		return
	}

	rb := rbr.Value
	mt := rb.Content.Get("application/json")
	if mt == nil {
		return
	}

	content := mt.Schema.Value
	if content == nil {
		return
	}

	for name, ref := range content.Properties {
		parameter := ref.Value
		name = getNameExtension(extensionPrefix, parameter.Extensions, name)

		for {
			_, exists := schema.Properties[name]
			if exists {
				name = "req-" + name
			} else {
				break
			}
		}

		schema.Properties[name] = ref

		if slices.Contains(content.Required, name) && !slices.Contains(schema.Required, name) {
			schema.Required = append(schema.Required, name)
		}
	}
}

var (
	parametersLocations = []string{openapi3.ParameterInPath, openapi3.ParameterInQuery}
	configLocations     = []string{openapi3.ParameterInHeader, openapi3.ParameterInCookie}
)

func (o *Operation) ParametersSchema() *core.Schema {
	if o.paramsSchema == nil {
		rootSchema := core.NewObjectSchema(map[string]*core.Schema{}, []string{})

		addParameters(rootSchema, o.path.Parameters, o.extensionPrefix, parametersLocations)
		addParameters(rootSchema, o.operation.Parameters, o.extensionPrefix, parametersLocations)
		addRequestBodyParameters(rootSchema, o.operation.RequestBody, o.extensionPrefix)

		o.paramsSchema = rootSchema
	}
	return o.paramsSchema
}

func (o *Operation) ConfigsSchema() *core.Schema {
	if o.configsSchema == nil {
		rootSchema := core.NewObjectSchema(map[string]*core.Schema{}, []string{})

		addParameters(rootSchema, o.path.Parameters, o.extensionPrefix, configLocations)
		addParameters(rootSchema, o.operation.Parameters, o.extensionPrefix, configLocations)

		o.configsSchema = rootSchema
	}
	return o.configsSchema
}

func (o *Operation) getServerURL(configs map[string]core.Value) (string, error) {
	// TODO: implement configs map[string]core.Value to replace variables in serverUrl
	var s *openapi3.Server
	if len(o.servers) > 0 {
		s = o.servers[0]
	}

	if s == nil {
		return "", fmt.Errorf("no available servers in spec")
	}
	return s.URL + o.key, nil
}

func replaceInPath(path string, param *openapi3.Parameter, val core.Value) (string, error) {
	var strValue string
	// TODO: handle complex conversion using openapi style values
	// https://spec.openapis.org/oas/latest.html#style-values
	switch valType := val.(type) {
	case []core.Value:
		return "", fmt.Errorf("Can not replace a slice into URL path: %v, param: %v", valType, param.Name)
	case map[core.Value]core.Value:
		return "", fmt.Errorf("Can not replace a map into URL path: %v, param: %v", valType, param.Name)
	default:
		strValue = fmt.Sprintf("%v", val)
	}
	paramTemplate := "{" + param.Name + "}"
	return strings.ReplaceAll(path, paramTemplate, strValue), nil
}

func addQueryParam(qValues *url.Values, param *openapi3.Parameter, value core.Value) {
	// TODO: handle complex conversion using openapi style values
	// https://spec.openapis.org/oas/latest.html#style-values
	if value == nil || fmt.Sprintf("%v", value) == "" {
		return
	}
	qValues.Set(param.Name, fmt.Sprintf("%v", value))
}

func addHeaderParam(req *http.Request, param *openapi3.Parameter, val core.Value) {
	// TODO: handle complex types passed on val
	req.Header.Set(param.Name, fmt.Sprintf("%v", val))
}

func addCookieParam(req *http.Request, param *openapi3.Parameter, val core.Value) {
	// TODO: handle complex types passed on val
	req.AddCookie(&http.Cookie{
		Name:  param.Name,
		Value: fmt.Sprintf("%v", val),
	})
}

func (o *Operation) buildRequestFromParams(
	paramValues map[string]core.Value,
	configs map[string]core.Value,
) (*http.Request, error) {
	serverURL, err := o.getServerURL(configs)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(strings.ToUpper(o.method), serverURL, nil)
	if err != nil {
		return nil, err
	}

	queryValues := url.Values{}
	for _, ref := range o.operation.Parameters {
		parameter := ref.Value
		name := getNameExtension(o.extensionPrefix, parameter.Extensions, parameter.Name)

		switch parameter.In {
		case openapi3.ParameterInPath:
			serverURL, err = replaceInPath(serverURL, parameter, paramValues[name])
			if err != nil {
				return nil, err
			}
		case openapi3.ParameterInQuery:
			addQueryParam(&queryValues, parameter, paramValues[name])
		case openapi3.ParameterInHeader:
			addHeaderParam(req, parameter, configs[name])
		case openapi3.ParameterInCookie:
			addCookieParam(req, parameter, configs[name])
		default:
			fmt.Printf("Unrecognizable param %s location %s", parameter.Name, parameter.In)
		}
	}
	if len(queryValues) > 0 {
		req.URL, err = url.Parse(serverURL + "?" + queryValues.Encode())
	} else {
		req.URL, err = url.Parse(serverURL)
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}

func setSecurityHeader(req *http.Request, secRequirements *openapi3.SecurityRequirements) error {
	if secRequirements == nil {
		return nil
	}
	for _, reqRef := range *secRequirements {
		for secScheme := range reqRef {
			if "oauth2" == strings.ToLower(secScheme) {
				// TODO: get token from config
				access_token := os.Getenv("MGC_SDK_ACCESS_TOKEN")
				if access_token == "" {
					return fmt.Errorf("Could not read acess token from env MGC_SDK_ACCESS_TOKEN")
				}
				req.Header.Set("Authorization", "Bearer "+access_token)
				return nil
			}
		}
	}
	return nil
}

// TODO: refactor this closer to the client that comes from a context
func (o *Operation) createHttpRequest(paramValues map[string]core.Value, configs map[string]core.Value) (*http.Request, error) {
	req, err := o.buildRequestFromParams(paramValues, configs)
	if err != nil {
		return nil, err
	}

	// TODO: accept everything, but later we need to fine-grain if json, multipart, etc
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")

	if err := setSecurityHeader(req, o.operation.Security); err != nil {
		return nil, err
	}

	return req, nil
}

func (o *Operation) Execute(
	ctx context.Context,
	parameters map[string]core.Value,
	configs map[string]core.Value,
) (result core.Value, err error) {
	// load definitions if not done yet
	parametersSchema := o.ParametersSchema()
	configsSchema := o.ConfigsSchema()

	client := core.HttpClientFromContext(ctx)
	if client == nil {
		return nil, fmt.Errorf("No HTTP client configured")
	}

	if err = parametersSchema.VisitJSON(parameters); err != nil {
		return nil, err
	}

	if err = configsSchema.VisitJSON(configs); err != nil {
		return nil, err
	}

	req, err := o.createHttpRequest(parameters, configs)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error on request response body: %s", err)
	}

	var data core.Value
	switch contentType := core.GetContentType(resp); contentType {
	default:
		// TODO: Handle other content types
		log.Fatalf("Unrecognized content-type %s in the response. Aborting", contentType)
	case "":
		// This will happen for 204 - No Content returns with empty body
		return map[string]core.Value{}, err
	case "application/json":
		data = map[string]core.Value{}
		err := core.DecodeJSON(resp, &data)
		return data, err
	}

	return data, err
}

var _ core.Executor = (*Operation)(nil)

// END: Executor interface
