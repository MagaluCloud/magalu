package common

import (
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

type BucketName string

func (b BucketName) JSONSchemaExtend(s *jsonschema.Schema) {
	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
	// s.MinLength = 3
	// s.MaxLength = 63
}

func (b BucketName) AsURI() mgcSchemaPkg.URI {
	if !strings.HasPrefix(string(b), URIPrefix) {
		return mgcSchemaPkg.URI(URIPrefix + b)
	}
	return mgcSchemaPkg.URI(b)
}

func (b BucketName) String() string {
	return string(b)
}

func NewBucketNameFromURI(u mgcSchemaPkg.URI) BucketName {
	return BucketName(u.Hostname())
}

type ObjectError struct {
	Url mgcSchemaPkg.URI
	Err error
}

func (e *ObjectError) Error() string {
	return fmt.Sprintf("%s - %s, ", e.Url, e.Err.Error())
}

func (e *ObjectError) Unwrap() error {
	return e.Err
}
