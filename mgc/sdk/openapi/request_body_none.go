package openapi

import (
	"io"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

type requestBodyNone struct{}

var _ requestBody = (*requestBodyNone)(nil)

var requestBodyNoneSingleton = &requestBodyNone{}

func newRequestBodyNone() *requestBodyNone {
	return requestBodyNoneSingleton
}

func (o *requestBodyNone) forEach(cb cbForEachParameterName) (finished bool, err error) {
	return true, nil
}

func (o *requestBodyNone) addToSchema(schema *mgcSchemaPkg.Schema) (err error) {
	return nil
}

func (o *requestBodyNone) create(pValues core.Parameters) (mimeType string, size int64, reader io.Reader, requestBody core.Value, err error) {
	size = -1
	return
}
