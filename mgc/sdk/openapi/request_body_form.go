package openapi

import (
	"errors"
	"io"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

type requestBodyForm struct {
	extensionPrefix *string
	logger          *zap.SugaredLogger
	mt              *openapi3.MediaType
}

var _ requestBody = (*requestBodyForm)(nil)

var errFormNotImplemented = errors.New("application/x-www-form-urlencoded not implemented")

func newRequestBodyForm(mt *openapi3.MediaType, logger *zap.SugaredLogger, extensionPrefix *string) *requestBodyForm {
	return &requestBodyForm{
		extensionPrefix: extensionPrefix,
		logger:          logger,
		mt:              mt,
	}
}

func (o *requestBodyForm) forEach(cb cbForEachParameterName) (finished bool, err error) {
	return false, errFormNotImplemented
}

func (o *requestBodyForm) addToSchema(schema *mgcSchemaPkg.Schema) (err error) {
	return errFormNotImplemented
}

func (o *requestBodyForm) create(pValues core.Parameters) (mimeType string, size int64, reader io.Reader, requestBody core.Value, err error) {
	err = errFormNotImplemented
	return
}
