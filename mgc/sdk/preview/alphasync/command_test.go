package alphasync

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const scopedSpec = `openapi: 3.0.3
info:
  title: Foo
  version: '1.0'
servers:
- url: https://{env}/{region}/foo
  variables:
    region:
      default: br-se1
    env:
      default: api.magalu.cloud
components:
  securitySchemes:
    OAuth2:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://id.magalu.com/login
          tokenUrl: https://id.magalu.com/oauth/token
          scopes:
            foo.read: Read foo
tags:
- name: things
paths:
  /v1/things:
    get:
      tags: [things]
      operationId: list-things
      summary: List things
      description: List all things.
      security:
      - OAuth2: [foo.read]
      responses:
        '200':
          description: OK
`

type acceptAll struct{}

func (acceptAll) Validate(preview.Level, preview.Package) error { return nil }

// exchange plays the token exchange of the login, recording the scopes asked.
type exchange struct {
	t      *testing.T
	asked  []string
	status int
}

func (e *exchange) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Tenant string `json:"tenant"`
		Scopes string `json:"scopes"`
	}
	data, _ := io.ReadAll(r.Body)
	require.NoError(e.t, json.Unmarshal(data, &body))
	e.asked = append(e.asked, body.Scopes)
	if e.status != 0 {
		w.WriteHeader(e.status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"access_token":  fakeJWT(e.t, strings.TrimPrefix(body.Tenant, "pf."), body.Scopes),
		"refresh_token": "refresh",
		"scope":         body.Scopes,
	})
}

func newLoggedIn(t *testing.T, scopes string) (*bucket, Config, *exchange, context.Context, *mgcAuthPkg.Auth) {
	t.Helper()
	b, cfg := newBucket(t)
	extensionPrefix := "x-mgc"
	cfg.Group = func(store preview.Store) core.Grouper {
		return preview.NewGroup(preview.Alpha, store, acceptAll{}, &extensionPrefix)
	}

	ex := &exchange{t: t}
	server := httptest.NewServer(ex)
	t.Cleanup(server.Close)
	configs := map[string]mgcAuthPkg.Config{"default": {TokenExchangeUrl: server.URL}}
	ctx, auth := newAuthContext(t, fakeJWT(t, tenantA, scopes), configs, nil)
	return b, cfg, ex, ctx, auth
}

func stderrOf(t *testing.T, run func()) string {
	t.Helper()
	read, write, err := os.Pipe()
	require.NoError(t, err)
	saved := os.Stderr
	os.Stderr = write
	run()
	os.Stderr = saved
	require.NoError(t, write.Close())
	out, err := io.ReadAll(read)
	require.NoError(t, err)
	return string(out)
}

func previewCommand(t *testing.T, cfg Config, name string) core.Executor {
	t.Helper()
	group, err := NewPreviewRoot(cfg).GetChildByName("preview")
	require.NoError(t, err)
	child, err := group.(core.Grouper).GetChildByName(name)
	require.NoError(t, err)
	return child.(core.Executor)
}

func TestPreviewCommands(t *testing.T) {
	_, cfg := newBucket(t)
	for _, name := range []string{"list", "install", "upgrade", "uninstall"} {
		exec := previewCommand(t, cfg, name)
		assert.NotContains(t, exec.ParametersSchema().Properties, "tenant", "the tenant is the one of the API key or of the login")
	}
}

func TestInstallAddsScopesToTheLogin(t *testing.T) {
	b, cfg, ex, ctx, _ := newLoggedIn(t, "openid")
	b.publish(t, EnvProd, tenantA, b.spec("foo", "1", scopedSpec))

	_, err := previewCommand(t, cfg, "install").Execute(ctx, core.Parameters{}, core.Configs{})
	require.NoError(t, err)
	require.Len(t, ex.asked, 1, "the login is older than the commands, so it lacks their scopes")
	assert.ElementsMatch(t, []string{"openid", "foo.read"}, strings.Fields(ex.asked[0]))

	foo := core.Parameters{"modules": []any{"foo"}}
	_, err = previewCommand(t, cfg, "install").Execute(ctx, foo, core.Configs{})
	require.NoError(t, err)
	assert.Len(t, ex.asked, 1, "a token with the scopes is not exchanged again")

	ex.status = http.StatusBadRequest
	b.publish(t, EnvProd, tenantA, b.spec("foo", "2", strings.Replace(scopedSpec, "foo.read", "foo.write", 2)))
	var result core.Result
	out := stderrOf(t, func() {
		result, err = previewCommand(t, cfg, "install").Execute(ctx, foo, core.Configs{})
	})
	require.NoError(t, err, "the commands are installed even when the scopes fail")
	assert.NotNil(t, result)
	assert.Contains(t, out, "run 'mgc auth login' again")
}

func TestInstallWithAnAPIKeyTellsTheScopes(t *testing.T) {
	b, cfg, ex, ctx, auth := newLoggedIn(t, "openid")
	require.NoError(t, auth.SetAPIKey("key"))
	var regions []string
	ctx = withTree(ctx, networkTree(t, map[string]any{
		"br-ne1": []map[string]any{{"is_default": true, "tenant_id": tenantB}},
	}, &regions))
	b.publish(t, EnvProd, tenantB, b.spec("foo", "1", scopedSpec))

	var err error
	out := stderrOf(t, func() {
		_, err = previewCommand(t, cfg, "install").Execute(ctx, core.Parameters{}, core.Configs{})
	})
	require.NoError(t, err)
	assert.Equal(t, manifestPath(t, EnvProd, tenantB), indexOf(t, cfg.Dir)["foo"].Manifest, "the tenant of the API key")
	assert.Empty(t, ex.asked, "an API key has no token to change")
	assert.Contains(t, out, "foo.read")
}

func TestInstallNeedsALoginOrAnAPIKey(t *testing.T) {
	b, cfg, _, _, _ := newLoggedIn(t, "openid")
	ctx, _ := newAuthContext(t, "", nil, nil)
	b.publish(t, EnvProd, tenantA, b.spec("foo", "1", scopedSpec))

	for _, name := range []string{"install", "list"} {
		_, err := previewCommand(t, cfg, name).Execute(ctx, core.Parameters{}, core.Configs{})
		assert.EqualError(t, err, "you are not logged in: run 'mgc auth login', or pass --api-key", name)
	}
	assert.Empty(t, b.hits)
}

func TestUpgradeWithoutLoginTellsTheScopes(t *testing.T) {
	b, cfg, ex, ctx, _ := newLoggedIn(t, "openid foo.read")
	b.publish(t, EnvProd, tenantA, b.spec("foo", "1", scopedSpec))
	_, err := previewCommand(t, cfg, "install").Execute(ctx, core.Parameters{}, core.Configs{})
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("foo", "2", scopedSpec+"# v2\n"))
	ctx, _ = newAuthContext(t, "", map[string]mgcAuthPkg.Config{"default": {TokenExchangeUrl: "http://127.0.0.1:1"}}, nil)
	out := stderrOf(t, func() {
		_, err = previewCommand(t, cfg, "upgrade").Execute(ctx, core.Parameters{}, core.Configs{})
	})
	require.NoError(t, err, "upgrade needs no tenant")
	assert.Equal(t, "2", indexOf(t, cfg.Dir)["foo"].Version)
	assert.Empty(t, ex.asked, "there is no token to change")
	assert.Contains(t, out, "need the scopes foo.read: log in, or use an API key that has them")
}

func TestInstallAndListByName(t *testing.T) {
	b, cfg, _, ctx, _ := newLoggedIn(t, "openid foo.read")
	b.publish(t, EnvProd, tenantA, b.spec("foo", "1", scopedSpec), b.spec("bar", "1", scopedSpec))

	_, err := previewCommand(t, cfg, "install").Execute(ctx, core.Parameters{"modules": []any{"foo@1"}}, core.Configs{})
	require.NoError(t, err)
	assert.True(t, indexOf(t, cfg.Dir)["foo"].Pinned)
	assert.NotContains(t, indexOf(t, cfg.Dir), "bar")

	result, err := previewCommand(t, cfg, "list").Execute(ctx, core.Parameters{}, core.Configs{})
	require.NoError(t, err)
	withValue, ok := core.ResultAs[core.ResultWithValue](result)
	require.True(t, ok)
	data, err := json.Marshal(withValue.Value())
	require.NoError(t, err)
	assert.Contains(t, string(data), `"installed":""`, "every row has every column")
	var listed []Listed
	require.NoError(t, json.Unmarshal(data, &listed))
	assert.Equal(t, []Listed{
		{Module: "bar", Status: "not installed", Latest: "1", Versions: []string{"1"}},
		{Module: "foo", Status: "pinned", Installed: "1", Latest: "1", Versions: []string{"1"}},
	}, listed)

	_, err = previewCommand(t, cfg, "uninstall").Execute(ctx, core.Parameters{"module": "foo"}, core.Configs{})
	require.NoError(t, err)
	assert.Empty(t, indexOf(t, cfg.Dir))
}
