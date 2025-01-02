package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

type requestBodyJSON struct {
	extensionPrefix *string
	logger          *zap.SugaredLogger
	mt              *openapi3.MediaType
}

var _ requestBody = (*requestBodyJSON)(nil)

func newRequestBodyJSON(mt *openapi3.MediaType, logger *zap.SugaredLogger, extensionPrefix *string) *requestBodyJSON {
	return &requestBodyJSON{
		extensionPrefix: extensionPrefix,
		logger:          logger,
		mt:              mt,
	}
}

func (o *requestBodyJSON) forEach(cb cbForEachParameterName) (finished bool, err error) {
	return o.forEachSchemaProperty(func(externalName, internalName string, _ *openapi3.SchemaRef, _ *openapi3.Schema) (run bool, err error) {
		return cb(externalName, internalName, "body")
	})
}

func (o *requestBodyJSON) forEachSchemaProperty(cb cbForEachSchemaProperty) (finished bool, err error) {
	names := map[string]bool{}
	finished, err = forEachMediaTypeProperty(o.mt, o.extensionPrefix, nil, func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
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
		err = &utils.ChainedError{Name: "application/json", Err: err}
	}
	return finished, err
}

func (o *requestBodyJSON) addToSchema(schema *core.Schema) (err error) {
	_, err = o.forEachSchemaProperty(func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with create()

		if existing := schema.Properties[externalName]; existing != nil {
			return false, &utils.ChainedError{Name: externalName, Err: fmt.Errorf("already exists as schema %v", existing)}
		}

		schema.Properties[externalName] = propRef

		if slices.Contains(containerSchema.Required, internalName) && !slices.Contains(schema.Required, externalName) {
			schema.Required = append(schema.Required, externalName)
		}
		return true, nil
	})

	slices.Sort(schema.Required)

	return
}

func (o *requestBodyJSON) create(pValues core.Parameters) (mimeType string, size int64, reader io.Reader, requestBody core.Value, err error) {
	size = -1

	body := map[string]core.Value{}
	_, err = o.forEachSchemaProperty(func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with addToSchema()

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
		err = fmt.Errorf("error encoding body content for request: %w", err)
		return
	}

	mimeType = "application/json"
	size = int64(bodyBuf.Len())
	reader = bodyBuf
	requestBody = body
	return
}
