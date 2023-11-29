package openapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/sdk/openapi/transform"

	"github.com/getkin/kin-openapi/openapi3"
)

const defaultResponseStatusCode = "default"

// Source -> Module -> Resource -> Operation

// Operation

type operation struct {
	core.SimpleDescriptor
	key                 string
	method              string
	path                *openapi3.PathItem
	operation           *openapi3.Operation
	paramsSchema        *core.Schema
	configsSchema       *core.Schema
	resultSchema        *core.Schema
	responseSchemas     map[string]*core.Schema
	links               core.Links
	related             map[string]core.Executor
	transformParameters func(value map[string]any) (map[string]any, error)
	transformConfigs    func(value map[string]any) (map[string]any, error)
	transformResult     func(value any) (any, error)
	extensionPrefix     *string
	outputFlag          string
	server              *server
	parameters          *parameters
	requestBody         requestBody
	logger              *zap.SugaredLogger
	refResolver         *core.BoundRefPathResolver
}

func newOperation(
	name string,
	desc *operationDesc,
	version string,
	method string,
	extensionPrefix *string,
	servers openapi3.Servers,
	logger *zap.SugaredLogger,
	outputFlag string,
	refResolver *core.BoundRefPathResolver,
) *operation {
	logger = logger.Named(name)
	return &operation{
		SimpleDescriptor: core.SimpleDescriptor{Spec: core.DescriptorSpec{
			Name:        name,
			Description: getDescriptionExtension(extensionPrefix, desc.op.Extensions, desc.op.Description),
			Version:     version,
			Summary:     desc.op.Summary,
			IsInternal:  getHiddenExtension(extensionPrefix, desc.op.Extensions),
		}},
		key:             desc.pathKey,
		method:          method,
		path:            desc.path,
		operation:       desc.op,
		extensionPrefix: extensionPrefix,
		logger:          logger,
		outputFlag:      outputFlag,
		refResolver:     refResolver,
		parameters:      newParameters(desc.pathKey, desc.path.Parameters, desc.op.Parameters, extensionPrefix),
		requestBody:     newRequestBody(method, desc.op, logger, extensionPrefix),
		server:          newServer(servers, extensionPrefix),
	}
}

var (
	parametersLocations = []string{openapi3.ParameterInPath, openapi3.ParameterInQuery}
	configLocations     = []string{openapi3.ParameterInHeader, openapi3.ParameterInCookie}
)

type cbForEachParameterName func(externalName, internalName, location string) (run bool, err error)

// Must match ParametersSchema!
func (o *operation) forEachParameterName(cb cbForEachParameterName) (finished bool, err error) {
	finished, err = o.parameters.forEach(parametersLocations, func(externalName string, parameter *openapi3.Parameter) (run bool, err error) {
		return cb(externalName, parameter.Name, parameter.In)
	})
	if !finished || err != nil {
		return
	}

	finished, err = o.requestBody.forEach(cb)
	if !finished || err != nil {
		return
	}

	// TODO: Walk through security requirement parameters, but currently they're not used in ParametersSchema
	// o.forEachSecurityRequirement(func(scheme string, scopes []string) (run bool, err error) {})
	return
}

func (o *operation) ParametersSchema() *core.Schema {
	if o.paramsSchema == nil {
		rootSchema := mgcSchemaPkg.NewObjectSchema(map[string]*core.Schema{}, []string{})
		var err error

		// Must match forEachParameterName!
		err = o.parameters.addToSchema(rootSchema, parametersLocations)
		if err != nil {
			o.logger.Warnw("error while adding parameters to schema", "error", err, "rootSchema", rootSchema)
		}

		err = o.requestBody.addToSchema(rootSchema)
		if err != nil {
			o.logger.Warnw("error while adding request body", "error", err)
		}

		o.addSecurityParameters(rootSchema)

		var transformSchema *core.Schema
		o.transformParameters, transformSchema, err = transform.New[map[string]any](o.logger, rootSchema, o.extensionPrefix)
		if err != nil {
			o.logger.Warnw("error while loading parameters schema", "error", err, "rootSchema", rootSchema)
		}

		if simplifiedParams, err := mgcSchemaPkg.SimplifySchema(transformSchema); err == nil {
			o.paramsSchema = simplifiedParams
		} else {
			o.logger.Warnw("error while simplifying params schema", "error", err, "transformSchema", transformSchema)
		}
	}
	return o.paramsSchema
}

func (o *operation) PositionalArgs() []string {
	return o.parameters.getPositionals()
}

func (o *operation) ConfigsSchema() *core.Schema {
	if o.configsSchema == nil {
		rootSchema := mgcSchemaPkg.NewObjectSchema(map[string]*core.Schema{}, []string{})
		var err error

		err = o.parameters.addToSchema(rootSchema, configLocations)
		if err != nil {
			o.logger.Warnw("error while adding parameters to configs schema", "error", err, "rootSchema", rootSchema)
		}

		err = o.server.addToSchema(rootSchema)
		if err != nil {
			o.logger.Warnw("error while adding server variables", "error", err)
		}

		var transformSchema *core.Schema
		o.transformConfigs, transformSchema, err = transform.New[map[string]any](o.logger, rootSchema, o.extensionPrefix)
		if err != nil {
			o.logger.Warnw("error while loading configs schema", "error", err, "rootSchema", rootSchema)
		}

		if simplifiedConfigs, err := mgcSchemaPkg.SimplifySchema(transformSchema); err == nil {
			o.configsSchema = simplifiedConfigs
		} else {
			o.logger.Warnw("error while simplifying configs schema", "error", err, "transformSchema", transformSchema)
		}
	}
	return o.configsSchema
}

func (o *operation) initResultSchema() {
	if o.resultSchema == nil {
		rootSchema := mgcSchemaPkg.NewAnyOfSchema()
		responses := o.operation.Responses
		o.responseSchemas = make(map[string]*core.Schema)

		for code, ref := range responses {
			if !(len(code) == 3 && strings.HasPrefix(code, "2")) && code != defaultResponseStatusCode {
				continue
			}

			response := ref.Value

			// TODO: Handle other media types
			content := response.Content.Get("application/json")
			if content == nil {
				continue
			}

			rootSchema.AnyOf = append(rootSchema.AnyOf, openapi3.NewSchemaRef(content.Schema.Ref, content.Schema.Value))
			o.responseSchemas[code] = (*core.Schema)(content.Schema.Value)
		}

		switch len(rootSchema.AnyOf) {
		default:
		case 0:
			rootSchema = mgcSchemaPkg.NewNullSchema()
		case 1:
			rootSchema = (*core.Schema)(rootSchema.AnyOf[0].Value)
		}

		var err error
		var transformSchema *core.Schema
		o.transformResult, transformSchema, err = transform.New[any](o.logger, rootSchema, o.extensionPrefix)
		if err != nil {
			o.logger.Warnw("error while initializing result schema", "error", err, "rootSchema", rootSchema)
		}
		if simplifiedResult, err := mgcSchemaPkg.SimplifySchema(transformSchema); err == nil {
			o.resultSchema = simplifiedResult
		} else {
			o.logger.Warnw("error while simplifying result schema", "error", err, "transformSchema", transformSchema)
		}
	}
}

func (o *operation) ResultSchema() *core.Schema {
	o.initResultSchema()
	return o.resultSchema
}

func (o *operation) resolveLink(link *openapi3.Link) (core.Executor, error) {
	if link.OperationID != "" {
		exec, err := core.ResolveExecutor(o.refResolver, "/"+operationIdsDocKey+"/"+link.OperationID)
		if err != nil {
			return nil, fmt.Errorf("linked operationId=%q: %w", link.OperationID, err)
		}
		return exec, nil
	} else if link.OperationRef != "" {
		exec, err := core.ResolveExecutor(o.refResolver, link.OperationRef)
		if err != nil {
			return nil, fmt.Errorf("linked operationRef=%q: %w", link.OperationRef, err)
		}
		return exec, nil
	} else {
		return nil, errors.New("link has no Operation ID or Ref")
	}
}

// This map should not be altered externally
func (o *operation) initLinksAndRelated() core.Links {
	if o.links == nil {
		o.links = core.Links{}
		o.related = map[string]core.Executor{}
		// TODO: Handle 'default' status code
		for _, respRef := range o.operation.Responses {
			resp := respRef.Value
			for key, linkRef := range resp.Links {
				link := linkRef.Value
				name := getNameExtension(o.extensionPrefix, link.Extensions, key)
				description := getDescriptionExtension(o.extensionPrefix, link.Extensions, link.Description)

				target, err := o.resolveLink(link)
				if err != nil {
					o.logger.Warnw("ignored broken link", "link", name, "error", err)
					continue
				}

				o.links[name] = &openapiLinker{
					name:        name,
					description: description,
					owner:       o,
					link:        linkRef.Value,
					target:      target,
				}
				o.related[name] = target
			}
		}
	}
	return o.links
}

// NOTE: it's possible to add new links using Links().AddLink(), but it's not possible to delete
// nor override any existing links
func (o *operation) Links() core.Links {
	o.initLinksAndRelated()
	return o.links
}

// This map should not be altered externally
func (o *operation) Related() map[string]core.Executor {
	o.initLinksAndRelated()
	return o.related
}

func (o *operation) getTransformResult() func(value any) (any, error) {
	// do this before checking o.transformResult as it will be initialized there
	o.initResultSchema()
	return o.transformResult
}

func (o *operation) getResponseSchemas() map[string]*core.Schema {
	// do this before checking o.responseSchemas as it will be initialized there
	o.initResultSchema()
	return o.responseSchemas
}

func replaceInPath(path string, param *openapi3.Parameter, val core.Value) string {
	// TODO: handle complex conversion using openapi style values
	// https://spec.openapis.org/oas/latest.html#style-values
	if val == nil || fmt.Sprintf("%v", val) == "" {
		return path
	}
	paramTemplate := "{" + param.Name + "}"
	return strings.ReplaceAll(path, paramTemplate, fmt.Sprintf("%v", val))
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

func closeIfCloser(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		_ = closer.Close()
	}
}

func (o *operation) getRequestUrl(
	paramValues core.Parameters,
	configs core.Configs,
) (string, error) {
	serverURL, err := o.server.url(configs)
	if err != nil {
		return "", err
	}

	queryValues := url.Values{}
	path := o.key
	_, err = o.parameters.forEachWithValue(paramValues, parametersLocations, func(externalName string, parameter *openapi3.Parameter, value any) (run bool, err error) {
		switch parameter.In {
		case openapi3.ParameterInPath:
			path = replaceInPath(path, parameter, value)
		case openapi3.ParameterInQuery:
			addQueryParam(&queryValues, parameter, value)
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}

	url := serverURL + path
	if len(queryValues) > 0 {
		url += "?" + queryValues.Encode()
	}

	return url, nil
}

func (o *operation) configureRequest(
	req *http.Request,
	configs core.Configs,
) {
	_, _ = o.parameters.forEachWithValue(configs, configLocations, func(externalName string, parameter *openapi3.Parameter, value any) (run bool, err error) {
		switch parameter.In {
		case openapi3.ParameterInHeader:
			addHeaderParam(req, parameter, value)
		case openapi3.ParameterInCookie:
			addCookieParam(req, parameter, value)
		}
		return true, nil
	})
}

func (o *operation) buildRequestFromParams(
	ctx context.Context,
	paramValues core.Parameters,
	configs core.Configs,
) (req *http.Request, requestBody core.Value, err error) {
	url, err := o.getRequestUrl(paramValues, configs)
	if err != nil {
		return
	}

	mimeType, size, reader, requestBody, err := o.requestBody.create(paramValues)
	if err != nil {
		return
	}
	// NOTE: from here on, error handling MUST close body!

	req, err = http.NewRequestWithContext(ctx, o.method, url, reader)
	if err != nil {
		closeIfCloser(reader)
		return
	}
	if mimeType != "" {
		req.Header.Set("Content-Type", mimeType)
	}
	if size > 0 {
		req.ContentLength = size
	}

	o.configureRequest(req, configs)
	return
}

func (o *operation) forEachSecurityRequirement(cb func(scheme string, scopes []string) (run bool, err error)) (finished bool, err error) {
	if o.operation.Security != nil {
		for _, reqRef := range *o.operation.Security {
			for scheme, scopes := range reqRef {
				scheme = strings.ToLower(scheme)
				if scheme == "oauth2" || scheme == "bearerauth" {
					run, err := cb(scheme, scopes)
					if err != nil {
						return false, err
					}
					if !run {
						return false, nil
					}
				} else {
					o.logger.Infow("ignored unsupported security scheme", "scheme", scheme, "scopes", scopes)
				}
			}
		}
	}
	return true, nil
}

func (o *operation) needsAuth() bool {
	finished, _ := o.forEachSecurityRequirement(func(scheme string, scopes []string) (run bool, err error) {
		return false, nil
	})

	return !finished // aborted early == had a security requirement
}

const forceAuthParameter = "force-authentication"

func (o *operation) addSecurityParameters(schema *core.Schema) {
	if o.needsAuth() {
		return
	}
	p := openapi3.NewBoolSchema()
	p.Description = "Force authentication by sending the header even if this API doesn't require it"
	schema.Properties[forceAuthParameter] = openapi3.NewSchemaRef("", p)
}

func isAuthForced(parameters core.Parameters) bool {
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

func (o *operation) setSecurityHeader(ctx context.Context, paramValues core.Parameters, req *http.Request, auth *mgcAuthPkg.Auth) (err error) {
	if isAuthForced(paramValues) || o.needsAuth() {
		// TODO: review needsAuth() usage if more security schemes are used. Assuming oauth2 + bearer
		// If others are to be used, loop using forEachSecurityRequirement()
		accessToken, err := auth.AccessToken(ctx)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}

	return nil
}

// TODO: refactor this closer to the client that comes from a context
func (o *operation) createHttpRequest(
	ctx context.Context,
	auth *mgcAuthPkg.Auth,
	paramValues core.Parameters,
	configs core.Configs,
) (req *http.Request, requestBody core.Value, err error) {
	req, requestBody, err = o.buildRequestFromParams(ctx, paramValues, configs)
	if err != nil {
		return
	}

	// TODO: accept everything, but later we need to fine-grain if json, multipart, etc
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")

	if err = o.setSecurityHeader(ctx, paramValues, req, auth); err != nil {
		return
	}

	return
}

func (o *operation) getValueFromResponseBody(value core.Value) (core.Value, error) {
	if transform := o.getTransformResult(); transform != nil {
		return transform(value)
	}
	return value, nil
}

func (o *operation) getResponseSchema(resp *http.Response) *core.Schema {
	responseSchemas := o.getResponseSchemas()
	code := fmt.Sprint(resp.StatusCode)
	if schema, ok := responseSchemas[code]; ok {
		return schema
	}
	if schema, ok := responseSchemas[defaultResponseStatusCode]; ok {
		return schema
	}
	return o.ResultSchema()
}

func isVisitJSONMultiErrFatal(multi openapi3.MultiError) bool {
	for _, err := range multi {
		if isVisitJSONErrFatal(err) {
			return true
		}
	}

	return false
}

func isVisitJSONErrFatal(err error) bool {
	var multiErr openapi3.MultiError
	ok := errors.As(err, &multiErr)
	if ok {
		return isVisitJSONMultiErrFatal(multiErr)
	}

	var schemaErr *openapi3.SchemaError
	ok = errors.As(err, &schemaErr)
	if ok {
		// Extra parameters are not fatal
		return !strings.Contains(schemaErr.Reason, "is unsupported")
	}

	return true
}

func (o *operation) Execute(
	ctx context.Context,
	parameters core.Parameters,
	configs core.Configs,
) (result core.Result, err error) {
	// keep the original parameters, configs -- do not use the transformed versions!
	// transformed versions will be new instances, so no worries changing the map pointer we reference here
	source := core.ResultSource{
		Executor:   o,
		Context:    ctx,
		Parameters: parameters,
		Configs:    configs,
	}

	// load definitions if not done yet
	parametersSchema := o.ParametersSchema()
	configsSchema := o.ConfigsSchema()

	client := mgcHttpPkg.ClientFromContext(ctx)
	if client == nil {
		return nil, fmt.Errorf("no HTTP client configured")
	}

	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("no Auth configured")
	}

	if err = parametersSchema.VisitJSON(parameters, openapi3.MultiErrors()); err != nil {
		if isVisitJSONErrFatal(err) {
			return nil, core.UsageError{Err: err}
		} else {
			o.logger.Warn(err)
		}
	}
	if o.transformParameters != nil {
		o.logger.Debugw("Starting parameter transforms", "parameters", parameters)
		// Safe because transformParameters doesn't modify the input map
		parameters, err = o.transformParameters(parameters)
		if err != nil {
			return nil, err
		}
		o.logger.Debugw("Finished parameter transforms", "transformed parameters", parameters)
	}

	if err = configsSchema.VisitJSON(configs, openapi3.MultiErrors()); err != nil {
		if isVisitJSONErrFatal(err) {
			return nil, core.UsageError{Err: err}
		} else {
			o.logger.Warn(err)
		}
	}
	if o.transformConfigs != nil {
		o.logger.Debugw("Starting config transforms", "configs", configs)
		// Safe because transformConfigs doesn't modify the input map
		configs, err = o.transformConfigs(configs)
		if err != nil {
			return nil, err
		}
		o.logger.Debugw("Finished config transforms", "transformed configs", configs)
	}

	req, requestBody, err := o.createHttpRequest(ctx, auth, parameters, configs)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp)
	}

	schema := o.getResponseSchema(resp)
	result, err = mgcHttpPkg.NewHttpResult(source, schema, req, requestBody, resp, o.getValueFromResponseBody)
	if err != nil {
		return nil, err
	}
	if o.outputFlag != "" {
		if resultWithValue, ok := core.ResultAs[core.ResultWithValue](result); ok {
			result = core.NewResultWithOriginalSource(result.Source(), core.NewResultWithDefaultOutputOptions(resultWithValue, o.outputFlag))
		}
	}
	return
}

// implemented by embedded SimpleDescriptor
var _ core.Executor = (*operation)(nil)
