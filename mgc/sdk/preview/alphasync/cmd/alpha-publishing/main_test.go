package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview/alphasync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fooSpec = `openapi: 3.0.3
info:
  title: Foo
  version: '1.0'
servers:
- url: https://{env}/{region}/foo
  variables:
    region:
      default: br-se1
tags:
- name: things
paths:
  /v1/things:
    get:
      tags: [things]
      operationId: list-things
      summary: List things
      description: List all things.
      responses:
        '200':
          description: OK
`

func setup(t *testing.T) (dir, keys string) {
	t.Helper()
	dir = t.TempDir()
	keys = filepath.Join(dir, "trusted_keys.json")
	require.NoError(t, os.WriteFile(keys, []byte("{}\n"), 0o644))
	return dir, keys
}

func TestKeygenThenSign(t *testing.T) {
	dir, keys := setup(t)
	private := filepath.Join(dir, "alpha.key")

	var out bytes.Buffer
	require.NoError(t, run([]string{"keygen", "--keys", keys, "--key-id", "2026-10", "--out", private}, &out))
	assert.Contains(t, out.String(), "store it in the vault")

	info, err := os.Stat(private)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(), "only the owner reads the private key")

	encoded, err := os.ReadFile(keys)
	require.NoError(t, err)
	trusted, err := alphasync.ParseKeys(encoded)
	require.NoError(t, err)
	assert.Contains(t, trusted, "2026-10")

	spec := filepath.Join(dir, "foo.openapi.yaml")
	require.NoError(t, os.WriteFile(spec, []byte(fooSpec), 0o644))

	signFoo := func() alphasync.ManifestSpec {
		out.Reset()
		require.NoError(t, run([]string{"sign", "--keys", keys, "--key", private, "--module", "foo",
			"--spec", spec, "--min-cli-version", "0.67.0"}, &out))
		var item alphasync.ManifestSpec
		require.NoError(t, json.Unmarshal(out.Bytes(), &item))
		return item
	}

	item := signFoo()
	assert.Equal(t, "foo", item.Module)
	assert.Regexp(t, `^alpha/specs/foo/[0-9a-f-]{36}\.openapi\.yaml$`, item.Path, "a path that cannot be guessed")
	assert.NotEqual(t, item.Path, signFoo().Path, "a new path for every signature")
	assert.Equal(t, "2026-10", item.KeyID)
	assert.Equal(t, "0.67.0", item.MinCLIVersion)
	assert.NoError(t, trusted.Verify(item), "the CLI accepts what sign prints")
}

func TestKeygenNeverOverwrites(t *testing.T) {
	dir, keys := setup(t)
	private := filepath.Join(dir, "alpha.key")
	require.NoError(t, run([]string{"keygen", "--keys", keys, "--key-id", "a", "--out", private}, &bytes.Buffer{}))
	before, err := os.ReadFile(private)
	require.NoError(t, err)

	err = run([]string{"keygen", "--keys", keys, "--key-id", "b", "--out", private}, &bytes.Buffer{})
	assert.Error(t, err, "an existing private key file is never overwritten")

	err = run([]string{"keygen", "--keys", keys, "--key-id", "a", "--out", filepath.Join(dir, "other.key")}, &bytes.Buffer{})
	assert.ErrorContains(t, err, "already trusted")

	after, err := os.ReadFile(private)
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestSignRefuses(t *testing.T) {
	dir, keys := setup(t)
	private := filepath.Join(dir, "alpha.key")
	require.NoError(t, run([]string{"keygen", "--keys", keys, "--key-id", "a", "--out", private}, &bytes.Buffer{}))

	write := func(name, data string) string {
		file := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(file, []byte(data), 0o644))
		return file
	}
	spec := write("foo.openapi.yaml", fooSpec)
	internal := write("internal.openapi.yaml", fooSpec+"security:\n- OAuth2: [foo-internal.read]\n")
	otherKey := write("other.key", "0000000000000000000000000000000000000000000000000000000000000000\n")
	badKey := write("bad.key", "not hex\n")

	sign := func(key, module, specFile string) error {
		return run([]string{"sign", "--keys", keys, "--key", key, "--module", module, "--spec", specFile}, &bytes.Buffer{})
	}
	assert.ErrorContains(t, sign(otherKey, "foo", spec), "does not pair", "a mistyped key is caught")
	assert.ErrorContains(t, sign(badKey, "foo", spec), "does not hold a private key")
	assert.ErrorContains(t, sign(private, "foo", internal), "the CLI would refuse this spec")
	assert.ErrorContains(t, sign(private, "../foo", spec), "not a valid name")
	assert.NoError(t, sign(private, "foo", spec))
}

func TestInvite(t *testing.T) {
	var out bytes.Buffer
	require.NoError(t, run([]string{"invite"}, &out))
	invite := strings.TrimSpace(out.String())
	_, err := alphasync.NormalizeInvite(invite)
	assert.NoError(t, err)

	out.Reset()
	require.NoError(t, run([]string{"invite"}, &out))
	assert.NotEqual(t, invite, strings.TrimSpace(out.String()))

	assert.Error(t, run([]string{"invite", "acme"}, &bytes.Buffer{}))
}
