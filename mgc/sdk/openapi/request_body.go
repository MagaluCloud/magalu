package openapi

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

type requestBody interface {
	forEach(cb cbForEachParameterName) (finished bool, err error)
	addToSchema(schema *mgcSchemaPkg.Schema) error
	create(pValues core.Parameters) (mimeType string, size int64, reader io.Reader, requestBody core.Value, err error)
}

func hasBody(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}

func newRequestBody(method string, op *openapi3.Operation, logger *zap.SugaredLogger, extensionPrefix *string) requestBody {
	if !hasBody(method) {
		return newRequestBodyNone()
	}

	rbr := op.RequestBody
	if rbr == nil {
		return newRequestBodyNone()
	}

	rb := rbr.Value
	if rb == nil {
		return newRequestBodyNone()
	}

	content := rb.Content
	if len(content) == 0 {
		return newRequestBodyNone()
	}

	if mt := content.Get("application/json"); mt != nil {
		return newRequestBodyJSON(mt, logger, extensionPrefix)
	} else if mt := content.Get("multipart/form-data"); mt != nil {
		return newRequestBodyMultipart(mt, logger, extensionPrefix)
	} else if mt := content.Get("application/x-www-form-urlencoded"); mt != nil {
		return newRequestBodyForm(mt, logger, extensionPrefix)
	} else {
		return newRequestBodySingle(content, logger, extensionPrefix)
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
	pValues core.Parameters,
) (filename string, mimeType string, size int64, file *os.File, err error) {
	size = -1

	v, ok := pValues[name]
	if !ok {
		err = fmt.Errorf("missing parameter %q", name)
		return
	}

	filename, ok = v.(string)
	if !ok {
		err = fmt.Errorf("parameter %q: not a string", name)
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

// Use this function as base to keep both parameter adding and processing in sync,
// with the same getExternalName function
//
// NOTE: getExternalName is only called if no extension provides the specific name
func forEachSchemaProperty(schema *openapi3.Schema, extensionPrefix *string, getExternalName cbGetName, cb cbForEachSchemaProperty) (finished bool, err error) {
	if schema == nil {
		return false, errors.New("missing schema")
	}

	if schema.Type != nil && !schema.Type.Includes("object") {
		return false, errors.New("must provide a schema with type 'object'")
	}

	for internalName, propRef := range schema.Properties {
		if propRef == nil {
			continue
		}

		propSchema := propRef.Value
		if propSchema == nil {
			continue
		}

		externalName := getNameExtension(extensionPrefix, propSchema.Extensions, "")
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

func forEachSchemaRefParameter(schemaRef *openapi3.SchemaRef, extensionPrefix *string, getExternalName cbGetName, cb cbForEachSchemaProperty) (finished bool, err error) {
	if schemaRef == nil {
		return false, errors.New("missing schemaRef")
	}
	return forEachSchemaProperty(schemaRef.Value, extensionPrefix, getExternalName, cb)
}

func forEachMediaTypeProperty(mediaType *openapi3.MediaType, extensionPrefix *string, getExternalName cbGetName, cb cbForEachSchemaProperty) (finished bool, err error) {
	return forEachSchemaRefParameter(mediaType.Schema, extensionPrefix, getExternalName, cb)
}
