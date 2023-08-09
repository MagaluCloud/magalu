package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"strings"

	"magalu.cloud/core"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

const fileUploadPrefix = "upload-"
const fileUploadParam = fileUploadPrefix + "file"

// Source -> Module -> Resource -> Operation

// Operation

type parameterWithName struct {
	name      string
	parameter *openapi3.Parameter
}

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
	parameters      *[]*parameterWithName
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

func (o *Operation) collectParameters(byNameAndLocation map[string]map[string]*parameterWithName, parameters openapi3.Parameters) {
	for _, ref := range parameters {
		// "A unique parameter is defined by a combination of a name and location."
		parameter := ref.Value

		byLocation, exists := byNameAndLocation[parameter.Name]
		if !exists {
			byLocation = map[string]*parameterWithName{}
			byNameAndLocation[parameter.Name] = byLocation
		}

		name := getNameExtension(o.extensionPrefix, parameter.Extensions, "")
		byLocation[parameter.In] = &parameterWithName{name, parameter}
	}
}

func (o *Operation) finalizeParameters(byNameAndLocation map[string]map[string]*parameterWithName) *[]*parameterWithName {
	parameters := []*parameterWithName{}

	for name, byLocation := range byNameAndLocation {
		for location, pn := range byLocation {
			if pn.name == "" {
				if len(byLocation) == 1 {
					pn.name = name
				} else {
					pn.name = fmt.Sprintf("%s-%s", location, name)
				}
			}
			parameters = append(parameters, pn)
		}
	}

	return &parameters
}

func (o *Operation) getParameters() []*parameterWithName {
	if o.parameters == nil {
		// operation parameters take precedence over path:
		// https://spec.openapis.org/oas/latest.html#fixed-fields-7
		// "the new definition will override it but can never remove it"
		// "A unique parameter is defined by a combination of a name and location."
		m := map[string]map[string]*parameterWithName{}
		o.collectParameters(m, o.path.Parameters)
		o.collectParameters(m, o.operation.Parameters)
		o.parameters = o.finalizeParameters(m)
	}
	return *o.parameters
}

type cbForEachParameter func(externalName string, parameter *openapi3.Parameter) (run bool, err error)

func (o *Operation) forEachParameter(locations []string, cb cbForEachParameter) (finished bool, err error) {
	for _, pn := range o.getParameters() {
		name := pn.name
		parameter := pn.parameter

		if parameter.Schema == nil || parameter.Schema.Value == nil {
			continue
		}

		if !slices.Contains(locations, parameter.In) {
			continue
		}
		if parameter.In == openapi3.ParameterInHeader && strings.HasPrefix(strings.ToLower(parameter.Name), "content-") {
			continue
		}

		run, err := cb(name, parameter)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}

	return true, nil
}

type cbForEachParameterWithValue func(externalName string, parameter *openapi3.Parameter, value any) (run bool, err error)

func (o *Operation) forEachParameterWithValue(values map[string]any, locations []string, cb cbForEachParameterWithValue) (finished bool, err error) {
	return o.forEachParameter(locations, func(externalName string, parameter *openapi3.Parameter) (run bool, err error) {
		value, ok := values[externalName]
		if !ok {
			value = parameter.Schema.Value.Default
			if value == nil {
				return true, nil
			}
		}
		return cb(externalName, parameter, value)
	})
}

func (o *Operation) addParameters(schema *core.Schema, locations []string) {
	_, err := o.forEachParameter(locations, func(externalName string, parameter *openapi3.Parameter) (run bool, err error) {
		paramSchemaRef := parameter.Schema
		paramSchema := paramSchemaRef.Value

		desc := getDescriptionExtension(o.extensionPrefix, parameter.Extensions, parameter.Description)
		if desc == "" {
			desc = getDescriptionExtension(o.extensionPrefix, paramSchema.Extensions, paramSchema.Description)
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

		schema.Properties[externalName] = paramSchemaRef

		if parameter.Required && !slices.Contains(schema.Required, externalName) {
			schema.Required = append(schema.Required, externalName)
		}

		return true, nil
	})

	if err != nil {
		// Log, but otherwise ignore the issue
		log.Printf("%s\n", err.Error())
	}
}

type cbGetName func(internalName string, propSchema *openapi3.Schema) string

type cbForEachSchemaProperty func(
	// the external (public/visible) name
	externalName string,
	// the name in the containerSchema.Properties
	internalName string,
	// the property schema, is guaranteed to not be null and have a non-null Value
	propRef *openapi3.SchemaRef,
	// The container schema (object definition)
	containerSchema *openapi3.Schema,
) (run bool, err error)

// Use this function as base to keep both parameter adding and processing in sync,
// with the same getExternalName function
//
// NOTE: getExternalName is only called if no extension provides the specific name
func (o *Operation) forEachSchemaProperty(schema *openapi3.Schema, getExternalName cbGetName, cb cbForEachSchemaProperty) (finished bool, err error) {
	if schema == nil {
		return false, errors.New("missing schema")
	}

	if schema.Type != openapi3.TypeObject {
		return false, errors.New("must provide a schema with type 'object'!")
	}

	for internalName, propRef := range schema.Properties {
		if propRef == nil {
			continue
		}

		propSchema := propRef.Value
		if propSchema == nil {
			continue
		}

		externalName := getNameExtension(o.extensionPrefix, propSchema.Extensions, "")
		if externalName == "" {
			if getExternalName != nil {
				externalName = getExternalName(internalName, propSchema)
			} else {
				externalName = internalName
			}
		}

		run, err := cb(externalName, internalName, propRef, schema)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}

	return true, nil
}

func (o *Operation) forEachSchemaRefParameter(schemaRef *openapi3.SchemaRef, getExternalName cbGetName, cb cbForEachSchemaProperty) (finished bool, err error) {
	if schemaRef == nil {
		return false, errors.New("missing schemaRef")
	}
	return o.forEachSchemaProperty(schemaRef.Value, getExternalName, cb)
}

func (o *Operation) forEachMediaTypeProperty(mediaType *openapi3.MediaType, getExternalName cbGetName, cb cbForEachSchemaProperty) (finished bool, err error) {
	return o.forEachSchemaRefParameter(mediaType.Schema, getExternalName, cb)
}

func (o *Operation) forEachBodyJsonParameter(mediaType *openapi3.MediaType, cb cbForEachSchemaProperty) (finished bool, err error) {
	names := map[string]bool{}
	finished, err = o.forEachMediaTypeProperty(mediaType, nil, func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		for {
			if names[externalName] {
				externalName = "req-" + externalName
			} else {
				break
			}
		}
		names[externalName] = true

		return cb(externalName, internalName, propRef, containerSchema)
	})

	if err != nil {
		err = fmt.Errorf("application/json %w", err)
	}
	return finished, err
}

func (o *Operation) addRequestBodyJsonParameters(mediaType *openapi3.MediaType, schema *core.Schema) (err error) {
	_, err = o.forEachBodyJsonParameter(mediaType, func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with createRequestBodyJson()

		schema.Properties[externalName] = propRef

		if slices.Contains(containerSchema.Required, internalName) && !slices.Contains(schema.Required, externalName) {
			schema.Required = append(schema.Required, externalName)
		}
		return true, nil
	})
	return
}

func (o *Operation) createRequestBodyJson(mediaType *openapi3.MediaType, pValues map[string]core.Value) (mimeType string, size int64, reader io.Reader, err error) {
	size = -1

	body := map[string]core.Value{}
	_, err = o.forEachBodyJsonParameter(mediaType, func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with addRequestBodyJsonParameters()

		if val, ok := pValues[externalName]; ok {
			body[internalName] = val
		}
		return true, nil
	})

	if err != nil {
		return
	}

	bodyBuf := new(bytes.Buffer)
	err = json.NewEncoder(bodyBuf).Encode(body)
	if err != nil {
		err = fmt.Errorf("Error encoding body content for request: %w", err)
		return
	}

	mimeType = "application/json"
	size = int64(bodyBuf.Len())
	reader = bodyBuf
	return
}

func getBodyUploadMultipartExternalName(internalName string, propSchema *openapi3.Schema) string {
	return fileUploadPrefix + internalName
}

func (o *Operation) forEachBodyUploadMultipartParameter(mediaType *openapi3.MediaType, cb cbForEachSchemaProperty) (finished bool, err error) {
	finished, err = o.forEachMediaTypeProperty(mediaType, getBodyUploadMultipartExternalName, cb)
	if err != nil {
		err = fmt.Errorf("multipart/form-data %w", err)
	}
	return finished, err
}

func (o *Operation) addRequestBodyUploadMultipartParameters(mediaType *openapi3.MediaType, schema *core.Schema) (err error) {
	_, err = o.forEachBodyUploadMultipartParameter(mediaType, func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with createRequestBodyUploadMultipart()

		// TODO: https://spec.openapis.org/oas/latest.html#special-considerations-for-multipart-content

		schema.Properties[externalName] = propRef

		if slices.Contains(containerSchema.Required, internalName) {
			schema.Required = append(schema.Required, externalName)
		}

		return true, nil
	})
	return
}

func (o *Operation) createRequestBodyUploadMultipart(
	mediaType *openapi3.MediaType,
	content openapi3.Content,
	pValues map[string]core.Value,
) (mimeType string, size int64, reader io.Reader, err error) {
	size = -1 // always -1 for multipart content

	type uploadEntry struct {
		name     string
		filename string
		mimeType string
		size     int64
		file     *os.File
	}
	uploads := []*uploadEntry{}

	_, err = o.forEachBodyUploadMultipartParameter(mediaType, func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with addRequestBodyUploadMultipartParameters()

		// TODO: https://spec.openapis.org/oas/latest.html#special-considerations-for-multipart-content

		filename, mime, sz, file, cerr := getFileFromParameter(externalName, pValues)
		if cerr == nil {
			e := &uploadEntry{
				name:     internalName,
				filename: filename,
				mimeType: mime,
				size:     sz,
				file:     file,
			}
			uploads = append(uploads, e)
		} else if slices.Contains(containerSchema.Required, internalName) {
			for _, e := range uploads {
				_ = e.file.Close()
			}
			return false, fmt.Errorf("Failed required parameter: %w", cerr)
		}
		return true, nil
	})

	if err != nil {
		return
	}

	r, w := io.Pipe()
	mw := multipart.NewWriter(w)
	go func() {
		// This goroutine fills the pipe's write side using multipart.Writer, processing one file at a time
		// io.Copy() + createFormFile() will block until the pipe's read side is used by the http.Client.Do()
		// as the read side will be the body reader
		defer w.Close()
		defer mw.Close()

		for _, e := range uploads {
			defer e.file.Close()
			part, cerr := createFormFile(mw, e.name, e.filename, e.mimeType, e.size)
			if cerr != nil {
				return
			}
			_, cerr = io.Copy(part, e.file)
			if cerr != nil {
				log.Printf("Warning: could not upload file %q: %v\n", e.name, cerr)
			}
		}
	}()

	mimeType = mw.FormDataContentType()
	reader = r
	return
}

func (o *Operation) addRequestBodyUploadFormParameters(mediaType *openapi3.MediaType, schema *core.Schema) (err error) {
	// NOTE: keep this paired with createRequestBodyUploadForm()

	err = fmt.Errorf("application/x-www-form-urlencoded not implemented")
	// TODO: https://spec.openapis.org/oas/latest.html#support-for-x-www-form-urlencoded-request-bodies

	return
}

func (o *Operation) createRequestBodyUploadForm(mediaType *openapi3.MediaType, content openapi3.Content, pValues map[string]core.Value) (mimeType string, size int64, reader io.Reader, err error) {
	// NOTE: keep this paired with addRequestBodyUploadFormParameters()

	size = -1
	err = fmt.Errorf("application/x-www-form-urlencoded not implemented")
	// TODO: https://spec.openapis.org/oas/latest.html#support-for-x-www-form-urlencoded-request-bodies
	return
}

func (o *Operation) addRequestBodyUploadSimpleParameters(content openapi3.Content, schema *core.Schema) (err error) {
	// NOTE: keep this paired with createRequestBodyUploadSimple()

	mimeTypes := make([]string, 0, len(content))
	for k, mediaType := range content {
		mimeTypes = append(mimeTypes, k)
		if mediaType.Schema != nil && mediaType.Schema.Value != nil && mediaType.Schema.Value.Type != "" {
			// spec gives the following example, that we do not support:
			//	Binary content transferred with base64 encoding:
			//		content:
			//			image/png:
			//				schema:
			//					type: string
			//					contentMediaType: image/png
			//					contentEncoding: base64
			log.Printf("content-type %q with schema is not supported: %+v\n", k, mediaType.Schema.Value)
		}
	}

	fs := openapi3.NewStringSchema()
	fs.Description = "File to be uploaded. Supported mime-types: " + strings.Join(mimeTypes, ", ")

	ref := openapi3.NewSchemaRef("", fs)

	name := fileUploadParam
	schema.Properties[name] = ref
	schema.Required = append(schema.Required, name)

	return
}

func (o *Operation) createRequestBodyUploadSimple(content openapi3.Content, pValues map[string]core.Value) (mimeType string, size int64, reader io.Reader, err error) {
	// NOTE: keep in sync with addRequestBodyUploadSimpleParameters
	_, mimeType, size, reader, err = getFileFromParameter(fileUploadParam, pValues)
	return mimeType, size, reader, err
}

func (o *Operation) hasBody() bool {
	switch o.method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}

func (o *Operation) createRequestBody(pValues map[string]core.Value) (mimeType string, size int64, reader io.Reader, err error) {
	// NOTE: keep in sync with addRequestBodyParameters()

	size = -1

	if !o.hasBody() {
		return
	}

	rbr := o.operation.RequestBody
	if rbr == nil {
		return
	}

	rb := rbr.Value
	if rb == nil {
		return
	}

	content := rb.Content
	if len(content) == 0 {
		return
	}

	if mt := content.Get("application/json"); mt != nil {
		return o.createRequestBodyJson(mt, pValues)
	} else if mt := content.Get("multipart/form-data"); mt != nil {
		return o.createRequestBodyUploadMultipart(mt, content, pValues)
	} else if mt := content.Get("application/x-www-form-urlencoded"); mt != nil {
		return o.createRequestBodyUploadForm(mt, content, pValues)
	} else {
		return o.createRequestBodyUploadSimple(content, pValues)
	}
}

func (o *Operation) addRequestBodyParameters(schema *core.Schema) {
	// NOTE: keep in sync with createRequestBody()

	if !o.hasBody() {
		return
	}

	rbr := o.operation.RequestBody
	if rbr == nil {
		return
	}

	rb := rbr.Value
	if rb == nil {
		return
	}

	content := rb.Content
	if len(content) == 0 {
		return
	}

	var err error
	if mt := content.Get("application/json"); mt != nil {
		err = o.addRequestBodyJsonParameters(mt, schema)
	} else if mt := content.Get("multipart/form-data"); mt != nil {
		err = o.addRequestBodyUploadMultipartParameters(mt, schema)
	} else if mt := content.Get("application/x-www-form-urlencoded"); mt != nil {
		err = o.addRequestBodyUploadFormParameters(mt, schema)
	} else {
		err = o.addRequestBodyUploadSimpleParameters(content, schema)
	}

	if err != nil {
		// Log, but otherwise ignore the issue
		log.Printf("%s\n", err.Error())
	}
}

var (
	parametersLocations = []string{openapi3.ParameterInPath, openapi3.ParameterInQuery}
	configLocations     = []string{openapi3.ParameterInHeader, openapi3.ParameterInCookie}
)

func (o *Operation) ParametersSchema() *core.Schema {
	if o.paramsSchema == nil {
		rootSchema := core.NewObjectSchema(map[string]*core.Schema{}, []string{})

		o.addParameters(rootSchema, parametersLocations)
		o.addRequestBodyParameters(rootSchema)
		o.addSecurityParameters(rootSchema)

		o.paramsSchema = rootSchema
	}
	return o.paramsSchema
}

func (o *Operation) ConfigsSchema() *core.Schema {
	if o.configsSchema == nil {
		rootSchema := core.NewObjectSchema(map[string]*core.Schema{}, []string{})

		o.addParameters(rootSchema, configLocations)
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

			rootSchema.AnyOf = append(rootSchema.AnyOf, openapi3.NewSchemaRef(content.Schema.Ref, content.Schema.Value))
		}

		if len(rootSchema.AnyOf) == 1 {
			rootSchema = (*core.Schema)(rootSchema.AnyOf[0].Value)
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

	return url, nil
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

func getFileMimeTypeAndSize(filename string, file *os.File) (mimeType string, size int64) {
	pos, _ := file.Seek(0, io.SeekCurrent)

	size, _ = file.Seek(0, io.SeekEnd)

	buffer := make([]byte, 512)
	_, _ = file.Read(buffer)
	mimeType = http.DetectContentType(buffer)
	_, _ = file.Seek(pos, io.SeekStart)

	if mimeType == "application/octet-stream" {
		ext := path.Ext(filename)
		if ext != "" {
			mimeType = mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
		}
	}

	return mimeType, size
}

func getFileFromParameter(
	name string,
	pValues map[string]core.Value,
) (filename string, mimeType string, size int64, file *os.File, err error) {
	size = -1

	v, ok := pValues[name]
	if !ok {
		err = fmt.Errorf("Missing parameter %q", name)
		return
	}

	filename, ok = v.(string)
	if !ok {
		err = fmt.Errorf("Parameter %q: not a string", name)
		return
	}

	file, err = os.Open(filename)
	if err != nil {
		return
	}

	filename = path.Base(filename)
	mimeType, size = getFileMimeTypeAndSize(filename, file)
	return
}

// BEGIN: these are like mime/multipart/writer.go, required because they are not exported
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// Variant of multipart.Writer.createFormFile() with mime-type
func createFormFile(w *multipart.Writer, fieldname, filename, mimeType string, size int64) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", mimeType)
	if size > 0 {
		h.Set("Content-Length", fmt.Sprintf("%d", size))
	}
	return w.CreatePart(h)
}

// END: these are like mime/multipart/writer.go, required because they are not exported

func closeIfCloser(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		_ = closer.Close()
	}
}

func (o *Operation) getRequestUrl(
	paramValues map[string]core.Value,
	configs map[string]core.Value,
) (string, error) {
	serverURL, err := o.getServerURL(configs)
	if err != nil {
		return "", err
	}

	queryValues := url.Values{}
	path := o.key
	_, err = o.forEachParameterWithValue(paramValues, parametersLocations, func(externalName string, parameter *openapi3.Parameter, value any) (run bool, err error) {
		switch parameter.In {
		case openapi3.ParameterInPath:
			path, err = replaceInPath(path, parameter, value)
			if err != nil {
				return false, err
			}
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

func (o *Operation) configureRequest(
	req *http.Request,
	configs map[string]core.Value,
) (err error) {
	_, err = o.forEachParameterWithValue(configs, configLocations, func(externalName string, parameter *openapi3.Parameter, value any) (run bool, err error) {
		switch parameter.In {
		case openapi3.ParameterInHeader:
			addHeaderParam(req, parameter, value)
		case openapi3.ParameterInCookie:
			addCookieParam(req, parameter, value)
		}
		return true, nil
	})
	return
}

func (o *Operation) buildRequestFromParams(
	paramValues map[string]core.Value,
	configs map[string]core.Value,
) (*http.Request, error) {
	url, err := o.getRequestUrl(paramValues, configs)
	if err != nil {
		return nil, err
	}

	mimeType, size, body, err := o.createRequestBody(paramValues)
	if err != nil {
		return nil, err
	}
	// NOTE: from here on, error handling MUST close body!

	req, err := http.NewRequest(o.method, url, body)
	if err != nil {
		closeIfCloser(body)
		return nil, err
	}
	if mimeType != "" {
		req.Header.Set("Content-Type", mimeType)
	}
	if size > 0 {
		req.ContentLength = size
	}

	err = o.configureRequest(req, configs)
	if err != nil {
		closeIfCloser(body)
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
	contentType, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch {
	default:
		return resp.Body, nil

	case strings.HasPrefix(contentType, "multipart/"):
		// TODO: do we have multi-part downloads? or just single?
		// If multi, then we need to return a multipart reader...
		// return multipart.NewReader(resp.Body, params["boundary"]), nil
		r := multipart.NewReader(resp.Body, params["boundary"])
		return r.NextPart()

	case contentType == "application/json":
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
		return nil, fmt.Errorf("HTTP request error: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, core.NewHttpErrorFromResponse(resp)
	}

	return o.getResponseValue(resp)
}

var _ core.Executor = (*Operation)(nil)

// END: Executor interface
