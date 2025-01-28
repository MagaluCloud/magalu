package openapi

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

type requestBodyMultipart struct {
	extensionPrefix *string
	logger          *zap.SugaredLogger
	mt              *openapi3.MediaType
}

var _ requestBody = (*requestBodyMultipart)(nil)

func newRequestBodyMultipart(mt *openapi3.MediaType, logger *zap.SugaredLogger, extensionPrefix *string) *requestBodyMultipart {
	return &requestBodyMultipart{
		extensionPrefix: extensionPrefix,
		logger:          logger,
		mt:              mt,
	}
}

func getBodyUploadMultipartExternalName(internalName string, propSchema *openapi3.Schema) string {
	return fileUploadPrefix + internalName
}

func (o *requestBodyMultipart) forEach(cb cbForEachParameterName) (finished bool, err error) {
	return o.forEachSchemaProperty(func(externalName, internalName string, _ *openapi3.SchemaRef, _ *openapi3.Schema) (run bool, err error) {
		return cb(externalName, internalName, "body")
	})
}

func (o *requestBodyMultipart) forEachSchemaProperty(cb cbForEachSchemaProperty) (finished bool, err error) {
	finished, err = forEachMediaTypeProperty(o.mt, o.extensionPrefix, getBodyUploadMultipartExternalName, cb)
	if err != nil {
		err = &utils.ChainedError{Name: "multipart/form-data", Err: err}
	}
	return finished, err
}

func (o *requestBodyMultipart) addToSchema(schema *mgcSchemaPkg.Schema) (err error) {
	_, err = o.forEachSchemaProperty(func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with create()

		// TODO: https://spec.openapis.org/oas/latest.html#special-considerations-for-multipart-content

		if existing := schema.Properties[externalName]; existing != nil {
			return false, &utils.ChainedError{Name: externalName, Err: fmt.Errorf("already exists as schema %v", existing)}
		}

		schema.Properties[externalName] = propRef

		if slices.Contains(containerSchema.Required, internalName) {
			schema.Required = append(schema.Required, externalName)
		}

		return true, nil
	})

	slices.Sort(schema.Required)

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

func (o *requestBodyMultipart) create(pValues core.Parameters) (mimeType string, size int64, reader io.Reader, requestBody core.Value, err error) {
	size = -1 // always -1 for multipart content

	type uploadEntry struct {
		name     string
		filename string
		mimeType string
		size     int64
		file     *os.File
	}
	uploads := []*uploadEntry{}

	_, err = o.forEachSchemaProperty(func(externalName, internalName string, propRef *openapi3.SchemaRef, containerSchema *openapi3.Schema) (run bool, err error) {
		// NOTE: keep this paired with addToSchema()

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
			return false, fmt.Errorf("failed required parameter: %w", cerr)
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
				o.logger.Warnw("could not upload file", "name", e.name, "file", e.filename, "error", cerr)
			}
		}
	}()

	mimeType = mw.FormDataContentType()
	reader = r
	return
}
