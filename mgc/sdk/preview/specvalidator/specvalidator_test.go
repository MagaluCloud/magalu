package specvalidator

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
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
