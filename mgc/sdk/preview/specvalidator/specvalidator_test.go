package specvalidator

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/sdk/openapi"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview/fsstore"
	"github.com/stretchr/testify/assert"
)

const validSpec = `openapi: 3.0.3
info:
  title: Foo
  version: '1.0'
paths:
  /v1/things:
    parameters: []
    get:
      responses:
        '200':
          description: OK
`

func TestValidate(t *testing.T) {
	cases := map[string]struct {
		name    string
		spec    string
		wantErr string
	}{
		"valid":               {name: "foo", spec: validSpec},
		"valid with dashes":   {name: "load-balancer", spec: validSpec},
		"valid with a digit":  {name: "s3_v2", spec: validSpec},
		"empty":               {name: "foo", spec: "", wantErr: "spec is empty"},
		"does not parse":      {name: "foo", spec: "info: [this is not\n", wantErr: "spec does not parse"},
		"no paths":            {name: "foo", spec: "openapi: 3.0.3\n", wantErr: "spec has no operations"},
		"paths without verbs": {name: "foo", spec: "paths:\n  /v1/things:\n    parameters: []\n", wantErr: "spec has no operations"},
		"name with a space":   {name: "foo bar", spec: validSpec, wantErr: `"foo bar" is not a valid beta module name`},
		"name with a dash":    {name: "-foo", spec: validSpec, wantErr: "not a valid beta module name"},
		"internal scope in an operation": {name: "foo", spec: validSpec + `      security:
      - OAuth2: [dnsaas-internal.read]
`, wantErr: `internal scope "dnsaas-internal.read"`},
		"internal scope in the global security": {name: "foo", spec: validSpec + `security:
- OAuth2: [foo.read, Foo-Internal.write]
`, wantErr: `internal scope "Foo-Internal.write"`},
		"internal scope in a security scheme": {name: "foo", spec: validSpec + `components:
  securitySchemes:
    OAuth2:
      type: oauth2
      flows:
        authorizationCode:
          scopes:
            foo-internal.read: Read
`, wantErr: `internal scope "foo-internal.read"`},
		"internal as an ordinary value": {name: "foo", spec: validSpec + `components:
  schemas:
    visibility:
      type: string
      enum: [internal, external]
`},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := New().Validate(preview.Beta, preview.Package{Name: tc.name, Spec: []byte(tc.spec)})
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestEmbeddedBetaSpecsAreValid(t *testing.T) {
	store := fsstore.New(openapi.BetaSpecs(), "embedded")
	names, err := store.Names()
	assert.NoError(t, err)
	assert.NotEmpty(t, names)

	for _, name := range names {
		pkg, err := store.Read(name)
		assert.NoError(t, err)
		assert.NoError(t, New().Validate(preview.Beta, pkg), name)
	}
}
