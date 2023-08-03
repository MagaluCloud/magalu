package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

// Source -> Module -> Resource -> Operation

// Operation

type Operation struct {
	name            string
	key             string
	method          string
	path            *openapi3.PathItem
	operation       *openapi3.Operation
	doc             *openapi3.T
	paramsSchema    *core.Schema
	configsSchema   *core.Schema
	resultSchema    *core.Schema
	extensionPrefix *string
	servers         openapi3.Servers
}

// BEGIN: Descriptor interface:

func (o *Operation) Name() string {
	return o.name
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
		o.addSecurityParameters(rootSchema)

		o.paramsSchema = rootSchema
	}
	return o.paramsSchema
}

func (o *Operation) ConfigsSchema() *core.Schema {
	if o.configsSchema == nil {
		rootSchema := core.NewObjectSchema(map[string]*core.Schema{}, []string{})

		addParameters(rootSchema, o.path.Parameters, o.extensionPrefix, configLocations)
		addParameters(rootSchema, o.operation.Parameters, o.extensionPrefix, configLocations)
		o.addServerVariables(rootSchema)

		o.configsSchema = rootSchema
	}
	return o.configsSchema
}

func (o *Operation) ResultSchema() *core.Schema {
	if o.resultSchema == nil {
		rootSchema := core.NewAnyOfSchema()
		responses := o.operation.Responses

		for code, ref := range responses {
			if !(len(code) == 3 && strings.HasPrefix(code, "2")) {
				continue
			}

			response := ref.Value

			// TODO: Handle other media types
			content := response.Content.Get("application/json")
			if content == nil {
				continue
			}

			// TODO: Don't use 'AnyOf' when there's only one subschema
			rootSchema.AnyOf = append(rootSchema.AnyOf, openapi3.NewSchemaRef(content.Schema.Ref, content.Schema.Value))
		}

		o.resultSchema = rootSchema
	}
	return o.resultSchema
}

type cbForEachVariable func(externalName string, internalName string, spec *openapi3.ServerVariable, server *openapi3.Server) (run bool, err error)

func (o *Operation) forEachServerVariable(cb cbForEachVariable) (finished bool, err error) {
	var s *openapi3.Server
	if len(o.servers) > 0 {
		s = o.servers[0]
	}

	if s == nil {
		return false, fmt.Errorf("no available servers in spec")
	}

	for internalName, spec := range s.Variables {
		externalName := getNameExtension(o.extensionPrefix, spec.Extensions, internalName)
		run, err := cb(externalName, internalName, spec, s)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}

	return true, nil
}

func (o *Operation) addServerVariables(schema *core.Schema) {
	_, err := o.forEachServerVariable(func(externalName, internalName string, spec *openapi3.ServerVariable, server *openapi3.Server) (run bool, err error) {
		varSchema := openapi3.NewStringSchema()
		varSchema.Default = spec.Default

		varSchema.Description = getDescriptionExtension(o.extensionPrefix, spec.Extensions, spec.Description)
		for _, e := range spec.Enum {
			varSchema.Enum = append(varSchema.Enum, e)
		}

		schema.Properties[externalName] = &openapi3.SchemaRef{Value: varSchema}
		return true, nil
	})

	if err != nil {
		// Log, but otherwise ignore the issue
		log.Printf("%s\n", err.Error())
	}
}

func (o *Operation) getServerURL(configs map[string]core.Value) (string, error) {
	url := ""
	_, err := o.forEachServerVariable(func(externalName, internalName string, spec *openapi3.ServerVariable, server *openapi3.Server) (run bool, err error) {
		val, ok := configs[externalName]
		if !ok {
			val = spec.Default
		}
		tmpl := "{" + internalName + "}"

		if url == "" {
			url = server.URL
		}
		url = strings.ReplaceAll(url, tmpl, fmt.Sprintf("%v", val))

		return true, nil
	})

	if err != nil {
		return "", err
	}

	return url + o.key, nil
}

func replaceInPath(path string, param *openapi3.Parameter, val core.Value) (string, error) {
	// TODO: handle complex conversion using openapi style values
	// https://spec.openapis.org/oas/latest.html#style-values
	if val == nil || fmt.Sprintf("%v", val) == "" {
		return path, nil
	}
	paramTemplate := "{" + param.Name + "}"
	return strings.ReplaceAll(path, paramTemplate, fmt.Sprintf("%v", val)), nil
}

func addQueryParam(qValues *url.Values, param *openapi3.Parameter, val core.Value) {
	// TODO: handle complex conversion using openapi style values
	// https://spec.openapis.org/oas/latest.html#style-values
	if val == nil || fmt.Sprintf("%v", val) == "" {
		if param.Schema != nil && param.Schema.Value != nil && param.Schema.Value.Default != nil {
			qValues.Set(param.Name, fmt.Sprintf("%v", param.Schema.Value.Default))
		}
		return
	} else {
		qValues.Set(param.Name, fmt.Sprintf("%v", val))
	}
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

func (o *Operation) createRequestBody(pValues map[string]core.Value) (string, *bytes.Buffer, error) {
	rbr := o.operation.RequestBody
	if rbr == nil {
		return "", bytes.NewBuffer(nil), nil
	}

	rb := rbr.Value
	mt := rb.Content.Get("application/json")
	if mt == nil {
		// TODO: parse body for multipart and other content types
		return "", nil, fmt.Errorf("Can not parse body for content type other than application/json")
	}

	content := mt.Schema.Value
	if content == nil {
		return "", nil, fmt.Errorf("Request body with empty schema ref")
	}

	body := map[string]core.Value{}
	for pKey, ref := range content.Properties {
		parameter := ref.Value
		name := getNameExtension(o.extensionPrefix, parameter.Extensions, pKey)
		if val, ok := pValues[name]; ok {
			body[pKey] = val
		}
	}
	bodyBuf := new(bytes.Buffer)
	if err := json.NewEncoder(bodyBuf).Encode(body); err != nil {
		return "", nil, fmt.Errorf("Error encoding body content for request: %v", err)
	}
	return "application/json", bodyBuf, nil
}

func (o *Operation) buildRequestFromParams(
	paramValues map[string]core.Value,
	configs map[string]core.Value,
) (*http.Request, error) {
	serverURL, err := o.getServerURL(configs)
	if err != nil {
		return nil, err
	}

	httpMethod := strings.ToUpper(o.method)
	bodyBuf := bytes.NewBuffer(nil)
	var mimeType string
	if httpMethod == http.MethodPost || httpMethod == http.MethodPut || httpMethod == http.MethodPatch {
		mimeType, bodyBuf, err = o.createRequestBody(paramValues)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(httpMethod, serverURL, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mimeType)

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

func (o *Operation) forEachSecurityRequirement(cb func(scheme string, scopes []string) (run bool, err error)) (finished bool, err error) {
	if o.operation.Security != nil {
		for _, reqRef := range *o.operation.Security {
			for scheme, scopes := range reqRef {
				scheme = strings.ToLower(scheme)
				if scheme == "oauth2" {
					run, err := cb(scheme, scopes)
					if err != nil {
						return false, err
					}
					if !run {
						return false, nil
					}
				} else {
					log.Printf("ignored unsupported security scheme: %q %+v\n", scheme, scopes)
				}
			}
		}
	}
	return true, nil
}

func (o *Operation) needsAuth() bool {
	finished, _ := o.forEachSecurityRequirement(func(scheme string, scopes []string) (run bool, err error) {
		return false, nil
	})

	return !finished // aborted early == had a security requirement
}

const forceAuthParameter = "force-authentication"

func (o *Operation) addSecurityParameters(schema *core.Schema) {
	if o.needsAuth() {
		return
	}
	p := openapi3.NewBoolSchema()
	p.Description = "Force authentication by sending the header even if this API doesn't require it"
	schema.Properties[forceAuthParameter] = openapi3.NewSchemaRef("", p)
}

func isAuthForced(parameters map[string]core.Value) bool {
	v, ok := parameters[forceAuthParameter]
	if !ok {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}

func (o *Operation) setSecurityHeader(paramValues map[string]core.Value, req *http.Request, auth *core.Auth) (err error) {
	if isAuthForced(paramValues) || o.needsAuth() {
		// TODO: review needsAuth() usage if more security schemes are used. Assuming oauth2 + bearer
		// If others are to be used, loop using forEachSecurityRequirement()
		accessToken, err := auth.AccessToken()
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}

	return nil
}

// TODO: refactor this closer to the client that comes from a context
func (o *Operation) createHttpRequest(
	auth *core.Auth,
	paramValues map[string]core.Value,
	configs map[string]core.Value,
) (*http.Request, error) {
	req, err := o.buildRequestFromParams(paramValues, configs)
	if err != nil {
		return nil, err
	}

	// TODO: accept everything, but later we need to fine-grain if json, multipart, etc
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")

	if err := o.setSecurityHeader(paramValues, req, auth); err != nil {
		return nil, err
	}

	return req, nil
}

func (o *Operation) getResponseValue(resp *http.Response) (core.Value, error) {
	if resp.StatusCode == 204 {
		return nil, nil
	}

	var data core.Value
	switch contentType := core.GetContentType(resp); contentType {
	default:
		return resp.Body, nil

	case "application/json":
		err := core.DecodeJSON(resp, &data)
		return data, err
	}
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

	auth := core.AuthFromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("No Auth configured")
	}

	if err = parametersSchema.VisitJSON(parameters); err != nil {
		return nil, err
	}

	if err = configsSchema.VisitJSON(configs); err != nil {
		return nil, err
	}

	req, err := o.createHttpRequest(auth, parameters, configs)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error on request response body: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, core.NewHttpErrorFromResponse(resp)
	}

	return o.getResponseValue(resp)
}

var _ core.Executor = (*Operation)(nil)

// END: Executor interface
