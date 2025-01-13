package openapi

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

type requestBodySingle struct {
	extensionPrefix *string
	logger          *zap.SugaredLogger
	content         openapi3.Content
}

var _ requestBody = (*requestBodySingle)(nil)

const fileUploadPrefix = "upload-"
const fileUploadParam = fileUploadPrefix + "file"

func newRequestBodySingle(content openapi3.Content, logger *zap.SugaredLogger, extensionPrefix *string) *requestBodySingle {
	return &requestBodySingle{
		extensionPrefix: extensionPrefix,
		logger:          logger,
		content:         content,
	}
}

func (o *requestBodySingle) forEach(cb cbForEachParameterName) (finished bool, err error) {
	return cb(fileUploadParam, fileUploadParam, "body")
}

func (o *requestBodySingle) addToSchema(schema *mgcSchemaPkg.Schema) (err error) {
	mimeTypes := make([]string, 0, len(o.content))
	for k, mediaType := range o.content {
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
			o.logger.Infow("content-type with schema is not supported", "content-type", k, "schema", mediaType.Schema.Value)
		}
	}

	fs := openapi3.NewStringSchema()
	fs.Description = "File to be uploaded. Supported mime-types: " + strings.Join(mimeTypes, ", ")

	ref := openapi3.NewSchemaRef("", fs)

	externalName := fileUploadParam

	if existing := schema.Properties[externalName]; existing != nil {
		return &utils.ChainedError{Name: externalName, Err: fmt.Errorf("already exists as schema %v", existing)}
	}

	schema.Properties[externalName] = ref
	schema.Required = append(schema.Required, externalName)
	slices.Sort(schema.Required)

	return nil
}

func (o *requestBodySingle) create(pValues core.Parameters) (mimeType string, size int64, reader io.Reader, requestBody core.Value, err error) {
	// NOTE: keep in sync with addToSchema
	_, mimeType, size, reader, err = getFileFromParameter(fileUploadParam, pValues)
	return
}
