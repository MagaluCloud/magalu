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
	"testing"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview/alphasync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	inviteA = "3f2a9c1e-5b7d-4e8a-9c0f-1a2b3c4d5e6f"
	inviteB = "9d8c7b6a-5f4e-4d3c-8b2a-1f0e9d8c7b6a"
)

type world struct {
	t       *testing.T
	dir     string
	keys    string
	private ed25519.PrivateKey
	objects map[string]string
	url     string
}

func newWorld(t *testing.T) *world {
	t.Helper()
	public, private, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	w := &world{t: t, dir: t.TempDir(), private: private, objects: map[string]string{}}
	w.keys = w.write("trusted_keys.json", `{"k1": "`+hex.EncodeToString(public)+`"}`)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, ok := w.objects[strings.TrimPrefix(r.URL.Path, "/")]
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
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

func (w *world) item(module, path, spec string) alphasync.ManifestSpec {
	sum := sha256.Sum256([]byte(spec))
	return alphasync.Sign(w.private, "k1", alphasync.ManifestSpec{Module: module, Path: path, SHA256: hex.EncodeToString(sum[:])})
}

func (w *world) itemFile(item alphasync.ManifestSpec) string {
	w.t.Helper()
	data, err := json.Marshal(item)
	require.NoError(w.t, err)
	return w.write(item.Module+"-item.json", string(data))
}

func (w *world) publishManifest(invite string, items ...alphasync.ManifestSpec) {
	w.t.Helper()
	data, err := json.Marshal(alphasync.Manifest{Specs: items})
	require.NoError(w.t, err)
	w.objects[alphasync.ManifestPath(invite)] = string(data)
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

func manifestModules(t *testing.T, data string) []string {
	t.Helper()
	var manifest alphasync.Manifest
	require.NoError(t, json.Unmarshal([]byte(data), &manifest))
	names := []string{}
	for _, spec := range manifest.Specs {
		names = append(names, spec.Module)
	}
	return names
}

func TestGrantKeepsWhatInvitesHave(t *testing.T) {
	w := newWorld(t)
	w.publishManifest(inviteA, w.item("dns", "alpha/specs/dns/1.openapi.yaml", "dns spec"))

	spec := w.write("iam.openapi.yaml", "iam spec")
	item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")
	invites := w.write("invites.txt", "# invited\n"+strings.ToUpper(inviteA)+"\n\n"+inviteB+"  # second\n")

	uploads, stdout, err := w.run("grant", "--spec", spec, "--item", w.itemFile(item), "--invites", invites)
	require.NoError(t, err)

	assert.Equal(t, "iam spec", uploads["alpha/specs/iam/1.openapi.yaml"])
	assert.Equal(t, []string{"dns", "iam"}, manifestModules(t, uploads[alphasync.ManifestPath(inviteA)]))
	assert.Equal(t, []string{"iam"}, manifestModules(t, uploads[alphasync.ManifestPath(inviteB)]))
	assert.Contains(t, stdout, "[dns] -> [dns iam]")
}

func TestGrantUploadsTheSpecFirst(t *testing.T) {
	w := newWorld(t)
	spec := w.write("iam.openapi.yaml", "iam spec")
	item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")

	out := filepath.Join(w.dir, "out")
	require.NoError(t, run([]string{"grant", "--keys", w.keys, "--base-url", w.url, "--out", out,
		"--spec", spec, "--item", w.itemFile(item), "--invite", inviteA}, &bytes.Buffer{}))

	plan, err := os.ReadFile(filepath.Join(out, planFile))
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(string(plan)), "\n")
	require.Len(t, lines, 2)
	assert.True(t, strings.HasSuffix(lines[0], "\talpha/specs/iam/1.openapi.yaml"))
	assert.True(t, strings.HasSuffix(lines[1], "\t"+alphasync.ManifestPath(inviteA)))
}

func TestGrantAgainUploadsNothing(t *testing.T) {
	w := newWorld(t)
	item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")
	w.objects[item.Path] = "iam spec"
	w.publishManifest(inviteA, item)

	uploads, stdout, err := w.run("grant", "--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item), "--invite", inviteA)
	require.NoError(t, err)
	assert.Empty(t, uploads)
	assert.Contains(t, stdout, "already published")
	assert.Contains(t, stdout, "unchanged")
}

func TestGrantRefuses(t *testing.T) {
	cases := map[string]func(w *world) []string{
		"spec path holds another spec": func(w *world) []string {
			item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")
			w.objects[item.Path] = "older iam spec"
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
		"spec is not the one signed": func(w *world) []string {
			item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")
			return []string{"--spec", w.write("iam.openapi.yaml", "edited iam spec"), "--item", w.itemFile(item)}
		},
		"item is not signed": func(w *world) []string {
			item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")
			item.Signature = ""
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
		"an invite has an item that does not verify": func(w *world) []string {
			bad := w.item("dns", "alpha/specs/dns/1.openapi.yaml", "dns spec")
			bad.MinCLIVersion = "0.1.0"
			w.publishManifest(inviteA, bad)
			item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")
			return []string{"--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item)}
		},
	}

	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			w := newWorld(t)
			args := append(setup(w), "--invite", inviteA)
			_, _, err := w.run("grant", args...)
			assert.Error(t, err)
		})
	}

	w := newWorld(t)
	item := w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec")
	_, _, err := w.run("grant", "--spec", w.write("iam.openapi.yaml", "iam spec"), "--item", w.itemFile(item), "--invite", "my-invite")
	assert.ErrorContains(t, err, "not an alpha invite")
}

func TestRevoke(t *testing.T) {
	w := newWorld(t)
	w.publishManifest(inviteA,
		w.item("dns", "alpha/specs/dns/1.openapi.yaml", "dns spec"),
		w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec"),
	)
	w.publishManifest(inviteB, w.item("iam", "alpha/specs/iam/1.openapi.yaml", "iam spec"))

	uploads, stdout, err := w.run("revoke", "--module", "iam", "--invite", inviteA, "--invite", inviteB)
	require.NoError(t, err)
	assert.Equal(t, []string{"dns"}, manifestModules(t, uploads[alphasync.ManifestPath(inviteA)]))
	assert.Equal(t, []string{}, manifestModules(t, uploads[alphasync.ManifestPath(inviteB)]), "an empty manifest stays, so the next sync removes it")

	uploads, stdout, err = w.run("revoke", "--module", "dns", "--invite", inviteB)
	require.NoError(t, err)
	assert.Empty(t, uploads)
	assert.Contains(t, stdout, "unchanged")
}
