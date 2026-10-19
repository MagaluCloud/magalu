package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview/alphasync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tenantA = "3f2a9c1e-5b7d-4e8a-9c0f-1a2b3c4d5e6f"
	tenantB = "9d8c7b6a-5f4e-4d3c-8b2a-1f0e9d8c7b6a"
)

func manifestPath(t *testing.T, env, tenant string) string {
	t.Helper()
	path, err := alphasync.ManifestPath(env, tenant)
	require.NoError(t, err)
	return path
}

type world struct {
	t       *testing.T
	dir     string
	keys    string
	private ed25519.PrivateKey
	objects map[string]string
	status  map[string]int
	url     string
}

func newWorld(t *testing.T) *world {
	t.Helper()
	public, private, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	w := &world{t: t, dir: t.TempDir(), private: private, objects: map[string]string{}, status: map[string]int{}}
	w.keys = w.write("trusted_keys.json", `{"k1": "`+hex.EncodeToString(public)+`"}`)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if status, ok := w.status[path]; ok {
			rw.WriteHeader(status)
			return
		}
		body, ok := w.objects[path]
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = rw.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	w.url = server.URL
	return w
}

func (w *world) write(name, data string) string {
	w.t.Helper()
	file := filepath.Join(w.dir, name)
	require.NoError(w.t, os.WriteFile(file, []byte(data), 0o644))
	return file
}

func (w *world) item(module, version, spec string) alphasync.ManifestSpec {
	return w.itemFor(alphasync.EnvProd, module, version, spec)
}

// Each item is published a second after the one before.
var publications atomic.Int64

func (w *world) itemFor(env, module, version, spec string) alphasync.ManifestSpec {
	sum := sha256.Sum256([]byte(spec))
	return alphasync.Sign(w.private, "k1", alphasync.ManifestSpec{
		Module:      module,
		Version:     version,
		PublishedAt: time.Unix(1_790_000_000+publications.Add(1), 0).UTC().Format(time.RFC3339),
		Env:         env,
		Path:        "alpha/specs/" + module + "/" + version + ".openapi.yaml",
		SHA256:      hex.EncodeToString(sum[:]),
	})
}

func (w *world) itemFile(item alphasync.ManifestSpec) string {
	w.t.Helper()
	data, err := json.Marshal(item)
	require.NoError(w.t, err)
	return w.write(item.Module+"@"+item.Version+"-item.json", string(data))
}

func (w *world) publishManifest(env, tenant string, items ...alphasync.ManifestSpec) {
	w.t.Helper()
	data, err := json.Marshal(alphasync.Manifest{Specs: items})
	require.NoError(w.t, err)
	w.objects[manifestPath(w.t, env, tenant)] = string(data)
}

func (w *world) run(command string, args ...string) (map[string]string, string, error) {
	w.t.Helper()
	out := filepath.Join(w.dir, "out-"+command+"-"+strings.ReplaceAll(w.t.Name(), "/", "_"))
	os.RemoveAll(out)

	var stdout bytes.Buffer
	err := run(append([]string{command, "--keys", w.keys, "--base-url", w.url, "--out", out}, args...), &stdout)
	if err != nil {
		return nil, stdout.String(), err
	}

	plan, err := os.ReadFile(filepath.Join(out, planFile))
	require.NoError(w.t, err)
	uploads := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(plan)), "\n") {
		if line == "" {
			continue
		}
		local, remote, ok := strings.Cut(line, "\t")
		require.True(w.t, ok, line)
		data, err := os.ReadFile(local)
		require.NoError(w.t, err)
		uploads[remote] = string(data)
	}
	return uploads, stdout.String(), nil
}

func manifestLabels(t *testing.T, data string) []string {
	t.Helper()
	var manifest alphasync.Manifest
	require.NoError(t, json.Unmarshal([]byte(data), &manifest))
	labels := []string{}
	for _, spec := range manifest.Specs {
		labels = append(labels, spec.Module+"@"+spec.Version)
	}
	return labels
}

func TestGrantKeepsWhatTenantsHave(t *testing.T) {
	w := newWorld(t)
	w.publishManifest(alphasync.EnvProd, tenantA, w.item("dns", "1", "dns spec"))

	spec := w.write("iam.openapi.yaml", "iam spec")
	item := w.item("iam", "1", "iam spec")
	tenants := w.write("tenants.txt", "# invited\n"+strings.ToUpper(tenantA)+"\n\n"+tenantB+"  # second\n")

	uploads, stdout, err := w.run("grant", "--env", "prod", "--spec", spec, "--item", w.itemFile(item), "--tenants", tenants)
	require.NoError(t, err)

	assert.Equal(t, "iam spec", uploads["alpha/specs/iam/1.openapi.yaml"])
	assert.Equal(t, []string{"dns@1", "iam@1"}, manifestLabels(t, uploads[manifestPath(t, alphasync.EnvProd, tenantA)]))
	assert.Equal(t, []string{"iam@1"}, manifestLabels(t, uploads[manifestPath(t, alphasync.EnvProd, tenantB)]))
	assert.Contains(t, stdout, "[dns@1] -> [dns@1 iam@1]")
}

func TestGrantAddsAVersion(t *testing.T) {
	w := newWorld(t)
	older := w.item("iam", "2026-10-01", "iam spec v1")
	w.objects[older.Path] = "iam spec v1"
	w.publishManifest(alphasync.EnvProd, tenantA, older)

	newer := w.item("iam", "2026-10-08", "iam spec v2")
	uploads, _, err := w.run("grant", "--env", "prod", "--spec", w.write("iam.openapi.yaml", "iam spec v2"), "--item", w.itemFile(newer), "--tenant", tenantA)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam@2026-10-01", "iam@2026-10-08"}, manifestLabels(t, uploads[manifestPath(t, alphasync.EnvProd, tenantA)]),
		"the older version stays, so customers can keep it pinned")
}

func TestGrantWritesTheManifestOfTheEnvironment(t *testing.T) {
	w := newWorld(t)
	spec := w.write("iam.openapi.yaml", "iam spec")
	item := w.itemFor(alphasync.EnvPreProd, "iam", "1", "iam spec")

	uploads, _, err := w.run("grant", "--env", "pre-prod", "--spec", spec, "--item", w.itemFile(item), "--tenant", tenantA)
	require.NoError(t, err)
	assert.Contains(t, uploads, manifestPath(t, alphasync.EnvPreProd, tenantA))
	assert.NotContains(t, uploads, manifestPath(t, alphasync.EnvProd, tenantA))

	_, _, err = w.run("grant", "--env", "prod", "--spec", spec, "--item", w.itemFile(item), "--tenant", tenantA)
	assert.ErrorContains(t, err, `signed for "pre-prod", not for "prod"`)
}

func TestGrantUploadsTheSpecFirst(t *testing.T) {
	w := newWorld(t)
	spec := w.write("iam.openapi.yaml", "iam spec")
	item := w.item("iam", "1", "iam spec")

	out := filepath.Join(w.dir, "out")
	require.NoError(t, run([]string{"grant", "--keys", w.keys, "--base-url", w.url, "--out", out,
		"--env", "prod", "--spec", spec, "--item", w.itemFile(item), "--tenant", tenantA}, &bytes.Buffer{}))

	plan, err := os.ReadFile(filepath.Join(out, planFile))
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(string(plan)), "\n")
	require.Len(t, lines, 2)
	assert.True(t, strings.HasSuffix(lines[0], "\talpha/specs/iam/1.openapi.yaml"))
	assert.True(t, strings.HasSuffix(lines[1], "\t"+manifestPath(t, alphasync.EnvProd, tenantA)))
}

func TestGrantAgainUploadsNothing(t *testing.T) {
	w := newWorld(t)
	item := w.item("iam", "1", "iam spec")
	w.objects[item.Path] = "iam spec"
	w.publishManifest(alphasync.EnvProd, tenantA, item)

	uploads, stdout, err := w.run("grant", "--env", "prod", "--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item), "--tenant", tenantA)
	require.NoError(t, err)
	assert.Empty(t, uploads)
	assert.Contains(t, stdout, "already published")
	assert.Contains(t, stdout, "unchanged")
}

func TestGrantRefuses(t *testing.T) {
	cases := map[string]func(w *world) []string{
		"spec path holds another spec": func(w *world) []string {
			item := w.item("iam", "1", "iam spec")
			w.objects[item.Path] = "older iam spec"
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
		"spec is not the one signed": func(w *world) []string {
			item := w.item("iam", "1", "iam spec")
			return []string{"--spec", w.write("iam.openapi.yaml", "edited iam spec"), "--item", w.itemFile(item)}
		},
		"item is not signed": func(w *world) []string {
			item := w.item("iam", "1", "iam spec")
			item.Signature = ""
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
		"version already published with other content": func(w *world) []string {
			w.publishManifest(alphasync.EnvProd, tenantA, w.item("iam", "1", "older iam spec"))
			item := w.item("iam", "1", "iam spec")
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
		"bucket refuses to read a manifest": func(w *world) []string {
			w.status[manifestPath(w.t, alphasync.EnvProd, tenantA)] = http.StatusForbidden
			item := w.item("iam", "1", "iam spec")
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
		"a tenant has an item that does not verify": func(w *world) []string {
			bad := w.item("dns", "1", "dns spec")
			bad.MinCLIVersion = "0.1.0"
			w.publishManifest(alphasync.EnvProd, tenantA, bad)
			item := w.item("iam", "1", "iam spec")
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
	}

	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			w := newWorld(t)
			args := append(setup(w), "--env", "prod", "--tenant", tenantA)
			_, _, err := w.run("grant", args...)
			assert.Error(t, err)
		})
	}

	w := newWorld(t)
	item := w.item("iam", "1", "iam spec")
	spec := w.write("iam.openapi.yaml", "iam spec")
	for _, tenant := range []string{"acme", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"} {
		_, _, err := w.run("grant", "--env", "prod", "--spec", spec, "--item", w.itemFile(item), "--tenant", tenant)
		assert.ErrorContains(t, err, "not a tenant id (UUID v4)", tenant)
	}
	_, _, err := w.run("grant", "--env", "dev", "--spec", spec, "--item", w.itemFile(item), "--tenant", tenantA)
	assert.ErrorContains(t, err, "unknown environment")
	_, _, err = w.run("grant", "--spec", spec, "--item", w.itemFile(item), "--tenant", tenantA)
	assert.ErrorContains(t, err, "--env")
}

func TestRevoke(t *testing.T) {
	w := newWorld(t)
	w.publishManifest(alphasync.EnvProd, tenantA,
		w.item("dns", "1", "dns spec"),
		w.item("iam", "1", "iam spec"),
	)
	w.publishManifest(alphasync.EnvProd, tenantB, w.item("iam", "1", "iam spec"))

	uploads, stdout, err := w.run("revoke", "--env", "prod", "--module", "iam", "--tenant", tenantA, "--tenant", tenantB)
	require.NoError(t, err)
	assert.Equal(t, []string{"dns@1"}, manifestLabels(t, uploads[manifestPath(t, alphasync.EnvProd, tenantA)]))
	assert.Equal(t, []string{}, manifestLabels(t, uploads[manifestPath(t, alphasync.EnvProd, tenantB)]), "an empty manifest stays, so the next sync removes it")

	uploads, stdout, err = w.run("revoke", "--env", "prod", "--module", "dns", "--tenant", tenantB)
	require.NoError(t, err)
	assert.Empty(t, uploads)
	assert.Contains(t, stdout, "unchanged")
}

func TestRevokeOneVersion(t *testing.T) {
	w := newWorld(t)
	w.publishManifest(alphasync.EnvProd, tenantA, w.item("iam", "1", "iam spec v1"), w.item("iam", "2", "iam spec v2"))

	uploads, _, err := w.run("revoke", "--env", "prod", "--module", "iam", "--version", "1", "--tenant", tenantA)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam@2"}, manifestLabels(t, uploads[manifestPath(t, alphasync.EnvProd, tenantA)]))
}
