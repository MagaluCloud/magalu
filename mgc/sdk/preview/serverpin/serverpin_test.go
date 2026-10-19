package serverpin

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview/fsstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const hostileSpec = `openapi: 3.0.3
info:
  title: Foo
  version: '1.0'
servers:
- url: https://{env}/{region}/foo
  variables:
    region:
      default: br-se1
    env:
      default: evil.example
      enum: [evil.example]
      x-mgc-transforms:
      - type: translate
        translations:
        - {from: prod, to: evil.example}
paths:
  /v1/things:
    servers:
    - url: https://evil.example/path-level
    get:
      servers:
      - url: https://{env}@evil.example/op-level
      responses:
        '200':
          description: OK
    post:
      responses:
        '201':
          description: Created
`

type spec struct {
	Servers []server `yaml:"servers"`
	Paths   map[string]struct {
		Servers []server `yaml:"servers"`
		Get     struct {
			Servers []server `yaml:"servers"`
		} `yaml:"get"`
	} `yaml:"paths"`
}

type server struct {
	URL       string         `yaml:"url"`
	Variables map[string]any `yaml:"variables"`
}

func canonicalEnv(t *testing.T) any {
	t.Helper()
	data, err := yaml.Marshal(envVariable)
	require.NoError(t, err)
	var env any
	require.NoError(t, yaml.Unmarshal(data, &env))
	return env
}

func TestPinURL(t *testing.T) {
	cases := map[string]string{
		"https://{env}/{region}/foo":          "https://{env}/{region}/foo",
		"https://evil.example/foo":            "https://{env}/foo",
		"https://evil.example":                "https://{env}",
		"https://{env}@evil.example/foo":      "https://{env}/foo",
		"http://evil.example:8080/foo?x=1":    "https://{env}/foo?x=1",
		"/":                                   "https://{env}/",
		"/v1":                                 "https://{env}/v1",
		"foo/v1":                              "https://{env}/foo/v1",
		"":                                    "https://{env}",
		"https://{env}//evil.example/foo":     "https://{env}//evil.example/foo",
		"https://evil.example/{region}/{env}": "https://{env}/{region}/{env}",
	}
	for raw, want := range cases {
		assert.Equal(t, want, PinURL(raw), raw)
	}
}

func TestPinRewritesEveryServer(t *testing.T) {
	pinned, err := Pin([]byte(hostileSpec))
	require.NoError(t, err)
	assert.NotContains(t, string(pinned), "evil.example")

	var got spec
	require.NoError(t, yaml.Unmarshal(pinned, &got))
	env := canonicalEnv(t)

	require.Len(t, got.Servers, 1)
	assert.Equal(t, "https://{env}/{region}/foo", got.Servers[0].URL)
	assert.Equal(t, env, got.Servers[0].Variables["env"])
	assert.Equal(t, map[string]any{"default": "br-se1"}, got.Servers[0].Variables["region"], "other variables stay")

	item := got.Paths["/v1/things"]
	require.Len(t, item.Servers, 1)
	assert.Equal(t, "https://{env}/path-level", item.Servers[0].URL)
	assert.Equal(t, env, item.Servers[0].Variables["env"])

	require.Len(t, item.Get.Servers, 1)
	assert.Equal(t, "https://{env}/op-level", item.Get.Servers[0].URL)
	assert.Equal(t, env, item.Get.Servers[0].Variables["env"])
}

func TestPinLeavesSpecsWithoutServersAlone(t *testing.T) {
	pinned, err := Pin([]byte("openapi: 3.0.3\npaths: {}\n"))
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, yaml.Unmarshal(pinned, &got))
	assert.NotContains(t, got, "servers")
}

func TestPinRefusesWhatDoesNotParse(t *testing.T) {
	_, err := Pin([]byte("info: [this is not\n"))
	assert.ErrorContains(t, err, "spec does not parse")
}

// The pin writes the env variable the pipeline compiles into every spec.
func TestEnvVariableMatchesTheCompiledSpecs(t *testing.T) {
	files, err := filepath.Glob("../../openapi/openapis/*.openapi.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		if filepath.Base(file) == "index.openapi.yaml" {
			continue
		}
		data, err := os.ReadFile(file)
		require.NoError(t, err)

		var got spec
		require.NoError(t, yaml.Unmarshal(data, &got))
		require.NotEmpty(t, got.Servers, file)
		assert.Equal(t, canonicalEnv(t), got.Servers[0].Variables["env"], file)
	}
}

func TestStorePinsWhatItReads(t *testing.T) {
	store := Wrap(fsstore.New(fstest.MapFS{"foo.openapi.yaml": {Data: []byte(hostileSpec)}}, "mem"))

	pkg, err := store.Read("foo")
	require.NoError(t, err)
	assert.NotContains(t, string(pkg.Spec), "evil.example")

	names, err := store.Names()
	require.NoError(t, err)
	assert.Equal(t, []string{"foo"}, names)
}
