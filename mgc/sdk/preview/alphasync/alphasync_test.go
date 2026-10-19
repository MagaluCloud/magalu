package alphasync

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tenantA = "3f2a9c1e-5b7d-4e8a-9c0f-1a2b3c4d5e6f"
	tenantB = "9d8c7b6a-5f4e-4d3c-8b2a-1f0e9d8c7b6a"
)

var testPublic, testPrivate, _ = ed25519.GenerateKey(nil)

const testKeyID = "test"

func checksum(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

// argon2id takes ~50 ms per key: tests compute each path once.
var paths sync.Map

func manifestPath(t *testing.T, env, tenant string) string {
	t.Helper()
	key := env + " " + tenant
	if path, ok := paths.Load(key); ok {
		return path.(string)
	}
	path, err := ManifestPath(env, tenant)
	require.NoError(t, err)
	paths.Store(key, path)
	return path
}

type bucket struct {
	files  map[string]string
	status map[string]int
	hits   []string
}

func (b *bucket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	b.hits = append(b.hits, path)
	if r.Header.Get("Authorization") != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if status, ok := b.status[path]; ok {
		w.WriteHeader(status)
		return
	}
	body, ok := b.files[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = w.Write([]byte(body))
}

func (b *bucket) publish(t *testing.T, env, tenant string, specs ...ManifestSpec) {
	t.Helper()
	data, err := json.Marshal(Manifest{Specs: specs})
	require.NoError(t, err)
	b.files[manifestPath(t, env, tenant)] = string(data)
}

// Each spec is published a second after the one before.
var publications atomic.Int64

func (b *bucket) spec(module, version, body string, opts ...func(*ManifestSpec)) ManifestSpec {
	path := "alpha/specs/" + module + "/" + version + ".openapi.yaml"
	b.files[path] = body
	spec := ManifestSpec{
		Module:      module,
		Version:     version,
		PublishedAt: time.Unix(1_790_000_000+publications.Add(1), 0).UTC().Format(time.RFC3339),
		Env:         EnvProd,
		Path:        path,
		SHA256:      checksum(body),
	}
	for _, opt := range opts {
		opt(&spec)
	}
	return Sign(testPrivate, testKeyID, spec)
}

func minVersion(v string) func(*ManifestSpec) {
	return func(s *ManifestSpec) { s.MinCLIVersion = v }
}

func forEnv(env string) func(*ManifestSpec) {
	return func(s *ManifestSpec) { s.Env = env }
}

func newBucket(t *testing.T) (*bucket, Config) {
	t.Helper()
	b := &bucket{files: map[string]string{}, status: map[string]int{}}
	server := httptest.NewServer(b)
	t.Cleanup(server.Close)

	cfg := Config{
		BaseURL:    server.URL,
		Dir:        filepath.Join(t.TempDir(), ".preview", "alpha"),
		CLIVersion: "v0.66.0",
		HTTPClient: server.Client(),
		Check:      func(preview.Package) error { return nil },
		Keys:       Keys{testKeyID: testPublic},
	}
	return b, cfg
}

func installed(t *testing.T, dir string) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	files := map[string]string{}
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		require.NoError(t, err)
		files[entry.Name()] = string(data)
	}
	return files
}

func TestTenantKey(t *testing.T) {
	// Pinned: the publishing tool and every released CLI must agree on it.
	key, err := TenantKey(EnvProd, tenantA)
	require.NoError(t, err)
	assert.Equal(t, "702cf6d1b0d68ed0bec9d98b921f802e2d069688369e780837bd7bb211a161d7", key)

	preProd, err := TenantKey(EnvPreProd, tenantA)
	require.NoError(t, err)
	assert.Equal(t, "5df62f3544a5a3c6f6a091164141144d6fb0d46da73476a22ea0fc87fc9aa119", preProd)

	same, err := TenantKey(EnvProd, " "+strings.ToUpper(tenantA)+" ")
	require.NoError(t, err)
	assert.Equal(t, key, same, "case and spaces do not matter")

	otherVariant, err := TenantKey(EnvProd, strings.Replace(tenantA, "-9c0f-", "-ac0f-", 1))
	require.NoError(t, err)
	assert.Equal(t, key, otherVariant, "the variant digit is left out")

	other, err := TenantKey(EnvProd, tenantB)
	require.NoError(t, err)
	assert.NotEqual(t, key, other)
}

func TestTenantKeyRefuses(t *testing.T) {
	for name, tenant := range map[string]string{
		"not a UUID":          "acme",
		"UUID v1":             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"not the RFC variant": "3f2a9c1e-5b7d-4e8a-cc0f-1a2b3c4d5e6f",
	} {
		_, err := TenantKey(EnvProd, tenant)
		assert.ErrorContains(t, err, "not a UUID v4", name)
	}

	_, err := TenantKey("dev", tenantA)
	assert.ErrorContains(t, err, "unknown environment")
}

func TestManifestPathHidesTheTenant(t *testing.T) {
	path := manifestPath(t, EnvProd, tenantA)
	assert.Regexp(t, `^alpha/tenants/[0-9a-f]{64}\.json$`, path)
	assert.NotContains(t, path, tenantA)
	assert.NotContains(t, path, strings.ReplaceAll(tenantA, "-", ""))
	assert.NotEqual(t, path, manifestPath(t, EnvProd, tenantB))
	assert.NotEqual(t, path, manifestPath(t, EnvPreProd, tenantA))
}

func TestCurrentEnv(t *testing.T) {
	assert.Equal(t, EnvPreProd, CurrentEnv("pre-prod"))
	assert.Equal(t, EnvProd, CurrentEnv("prod"))
	assert.Equal(t, EnvProd, CurrentEnv(""), "as the login does")
	assert.Equal(t, EnvProd, CurrentEnv("api.magalu.cloud"))
}

func TestManifestFindsVersions(t *testing.T) {
	b, _ := newBucket(t)
	older := b.spec("iam", "2026-10-01", "iam spec v1")
	newer := b.spec("iam", "2026-10-08", "iam spec v2")
	dns := b.spec("dns", "1", "dns spec")
	manifest := Manifest{Specs: []ManifestSpec{newer, dns, older}}
	require.NoError(t, manifest.Validate())

	found, ok := manifest.Find("iam", "")
	assert.True(t, ok)
	assert.Equal(t, newer, found, "no version means the newest")

	found, ok = manifest.Find("iam", "2026-10-01")
	assert.True(t, ok)
	assert.Equal(t, older, found)

	_, ok = manifest.Find("iam", "2026-10-02")
	assert.False(t, ok)
	_, ok = manifest.Find("connect", "")
	assert.False(t, ok)

	assert.Equal(t, []ManifestSpec{dns, newer}, manifest.Latest())
}

func TestManifestValidateRefuses(t *testing.T) {
	b, _ := newBucket(t)
	for name, change := range map[string]func(*ManifestSpec){
		"invalid module":       func(s *ManifestSpec) { s.Module = "../iam" },
		"no version":           func(s *ManifestSpec) { s.Version = "" },
		"version with a slash": func(s *ManifestSpec) { s.Version = "1/2" },
		"no publication date":  func(s *ManifestSpec) { s.PublishedAt = "" },
		"date outside UTC":     func(s *ManifestSpec) { s.PublishedAt = "2026-10-18T09:00:00-03:00" },
		"date with fractions":  func(s *ManifestSpec) { s.PublishedAt = "2026-10-18T12:00:00.5Z" },
	} {
		spec := b.spec("iam", "1", "iam spec")
		change(&spec)
		assert.Error(t, Manifest{Specs: []ManifestSpec{spec}}.Validate(), name)
	}

	spec := b.spec("iam", "1", "iam spec")
	assert.ErrorContains(t, Manifest{Specs: []ManifestSpec{spec, spec}}.Validate(), "iam@1 is listed twice")
}

func TestSignedPayloadCoversWhatMatters(t *testing.T) {
	base := ManifestSpec{Module: "iam", Version: "1", PublishedAt: "2026-10-18T12:00:00Z", Env: EnvProd, Path: "a", SHA256: "ABC", MinCLIVersion: "0.66.0", KeyID: "k"}

	moved := base
	moved.Path = "b"
	assert.Equal(t, SignedPayload(base), SignedPayload(moved), "the path is not signed")

	for name, change := range map[string]func(*ManifestSpec){
		"module":          func(s *ManifestSpec) { s.Module = "dns" },
		"version":         func(s *ManifestSpec) { s.Version = "2" },
		"published_at":    func(s *ManifestSpec) { s.PublishedAt = "2026-10-19T12:00:00Z" },
		"env":             func(s *ManifestSpec) { s.Env = EnvPreProd },
		"sha256":          func(s *ManifestSpec) { s.SHA256 = "abd" },
		"min_cli_version": func(s *ManifestSpec) { s.MinCLIVersion = "" },
		"key_id":          func(s *ManifestSpec) { s.KeyID = "other" },
	} {
		changed := base
		change(&changed)
		assert.NotEqual(t, SignedPayload(base), SignedPayload(changed), name)
	}
}

func TestTrustedKeysParse(t *testing.T) {
	assert.NotPanics(t, func() { TrustedKeys() })

	keys, err := ParseKeys([]byte(`{"2026-10": "` + hex.EncodeToString(testPublic) + `"}`))
	require.NoError(t, err)
	id, err := keys.KeyIDOf(testPrivate)
	require.NoError(t, err)
	assert.Equal(t, "2026-10", id)

	_, other, _ := ed25519.GenerateKey(nil)
	_, err = keys.KeyIDOf(other)
	assert.ErrorContains(t, err, "does not pair")

	_, err = ParseKeys([]byte(`{"bad": "abcd"}`))
	assert.ErrorContains(t, err, "not an Ed25519 public key")
}
