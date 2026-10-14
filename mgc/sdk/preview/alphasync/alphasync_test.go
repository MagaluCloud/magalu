package alphasync

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview/fsstore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	inviteA = "3f2a9c1e-5b7d-4e8a-9c0f-1a2b3c4d5e6f"
	inviteB = "9d8c7b6a-5f4e-4d3c-8b2a-1f0e9d8c7b6a"
)

var testPublic, testPrivate, _ = ed25519.GenerateKey(nil)

const testKeyID = "test"

func checksum(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
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

func (b *bucket) publish(t *testing.T, invite string, specs ...ManifestSpec) {
	t.Helper()
	data, err := json.Marshal(Manifest{Specs: specs})
	require.NoError(t, err)
	b.files[ManifestPath(invite)] = string(data)
}

func (b *bucket) spec(module, version, body string, opts ...func(*ManifestSpec)) ManifestSpec {
	path := "alpha/specs/" + module + "/" + version + ".openapi.yaml"
	b.files[path] = body
	spec := ManifestSpec{Module: module, Path: path, SHA256: checksum(body)}
	for _, opt := range opts {
		opt(&spec)
	}
	return Sign(testPrivate, testKeyID, spec)
}

func minVersion(v string) func(*ManifestSpec) {
	return func(s *ManifestSpec) { s.MinCLIVersion = v }
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

func TestNewInvite(t *testing.T) {
	invite := NewInvite()
	id, err := uuid.Parse(invite)
	require.NoError(t, err)
	assert.Equal(t, uuid.Version(4), id.Version(), "122 random bits")
	assert.NotEqual(t, invite, NewInvite())

	normalized, err := NormalizeInvite(invite)
	require.NoError(t, err)
	assert.Equal(t, invite, normalized)
}

func TestNormalizeInvite(t *testing.T) {
	got, err := NormalizeInvite(" " + strings.ToUpper(inviteA) + " ")
	require.NoError(t, err)
	assert.Equal(t, inviteA, got)

	_, err = NormalizeInvite("my-invite")
	assert.ErrorContains(t, err, "not an alpha invite")
}

func TestManifestPathHidesTheInvite(t *testing.T) {
	path := ManifestPath(inviteA)
	assert.Regexp(t, `^alpha/invites/[0-9a-f]{64}\.json$`, path)
	assert.NotContains(t, path, inviteA)
	assert.NotEqual(t, path, ManifestPath(inviteB))
}

func TestSyncInstallsWhatTheInviteGives(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, inviteA, b.spec("iam", "1", "iam spec"), b.spec("dns", "1", "dns spec"))

	result, err := Sync(context.Background(), cfg, strings.ToUpper(inviteA))
	require.NoError(t, err)
	assert.Equal(t, []string{"dns", "iam"}, result.Installed)

	files := installed(t, cfg.Dir)
	assert.Equal(t, "iam spec", files["iam.openapi.yaml"])
	assert.Equal(t, "dns spec", files["dns.openapi.yaml"])
	assert.Contains(t, files, "index.json")

	saved, err := SavedInvite(cfg.Dir)
	require.NoError(t, err)
	assert.Equal(t, inviteA, saved)

	info, err := os.Stat(filepath.Join(cfg.Dir, "invite.json"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestSyncWithoutManifestRemovesEverything(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusForbidden} {
		b, cfg := newBucket(t)
		b.publish(t, inviteA, b.spec("iam", "1", "iam spec"))
		_, err := Sync(context.Background(), cfg, inviteA)
		require.NoError(t, err)

		delete(b.files, ManifestPath(inviteA))
		b.status[ManifestPath(inviteA)] = status

		result, err := Sync(context.Background(), cfg, inviteA)
		require.NoError(t, err)
		assert.Empty(t, result.Installed)
		assert.NotContains(t, installed(t, cfg.Dir), "iam.openapi.yaml")
	}
}

func TestSyncRemovesWhatIsNoLongerGiven(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, inviteA, b.spec("iam", "1", "iam spec"), b.spec("dns", "1", "dns spec"))
	_, err := Sync(context.Background(), cfg, inviteA)
	require.NoError(t, err)

	b.publish(t, inviteA, b.spec("iam", "2", "iam spec v2"))
	result, err := Sync(context.Background(), cfg, inviteA)
	require.NoError(t, err)

	assert.Equal(t, []string{"iam"}, result.Installed)
	files := installed(t, cfg.Dir)
	assert.Equal(t, "iam spec v2", files["iam.openapi.yaml"])
	assert.NotContains(t, files, "dns.openapi.yaml")
}

func TestSyncKeepsTheInstalledSpecsOnFailure(t *testing.T) {
	cases := map[string]func(t *testing.T, b *bucket, cfg *Config){
		"checksum mismatch": func(t *testing.T, b *bucket, cfg *Config) {
			b.publish(t, inviteA, b.spec("iam", "2", "iam spec v2", func(s *ManifestSpec) { s.SHA256 = checksum("something else") }))
		},
		"not signed": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2")
			spec.KeyID, spec.Signature = "", ""
			b.publish(t, inviteA, spec)
		},
		"signed fields changed": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2", minVersion("0.70.0"))
			spec.MinCLIVersion = ""
			b.publish(t, inviteA, spec)
		},
		"module renamed": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("dns", "2", "dns spec")
			spec.Module = "iam"
			b.publish(t, inviteA, spec)
		},
		"signed with an unknown key": func(t *testing.T, b *bucket, cfg *Config) {
			_, other, _ := ed25519.GenerateKey(nil)
			spec := b.spec("iam", "2", "iam spec v2")
			b.publish(t, inviteA, Sign(other, "other", spec))
		},
		"signed with another key under a trusted id": func(t *testing.T, b *bucket, cfg *Config) {
			_, other, _ := ed25519.GenerateKey(nil)
			spec := b.spec("iam", "2", "iam spec v2")
			b.publish(t, inviteA, Sign(other, testKeyID, spec))
		},
		"spec does not load": func(t *testing.T, b *bucket, cfg *Config) {
			b.publish(t, inviteA, b.spec("iam", "2", "iam spec v2"))
			cfg.Check = func(preview.Package) error { return errors.New("spec has no operations") }
		},
		"spec missing in the bucket": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2")
			delete(b.files, spec.Path)
			b.publish(t, inviteA, spec)
		},
		"module listed twice": func(t *testing.T, b *bucket, cfg *Config) {
			b.publish(t, inviteA, b.spec("iam", "2", "iam spec v2"), b.spec("iam", "3", "iam spec v3"))
		},
		"invalid module name": func(t *testing.T, b *bucket, cfg *Config) {
			b.publish(t, inviteA, b.spec("../iam", "2", "iam spec v2"))
		},
		"source fails": func(t *testing.T, b *bucket, cfg *Config) {
			b.status[ManifestPath(inviteA)] = http.StatusInternalServerError
		},
	}

	for name, breakIt := range cases {
		t.Run(name, func(t *testing.T) {
			b, cfg := newBucket(t)
			b.publish(t, inviteA, b.spec("iam", "1", "iam spec"))
			_, err := Sync(context.Background(), cfg, inviteA)
			require.NoError(t, err)
			before := installed(t, cfg.Dir)

			breakIt(t, b, &cfg)
			_, err = Sync(context.Background(), cfg, inviteA)
			assert.Error(t, err)
			assert.Equal(t, before, installed(t, cfg.Dir))
		})
	}
}

func TestSyncSkipsSpecsForNewerCLIs(t *testing.T) {
	b, cfg := newBucket(t)
	newer := b.spec("dns", "1", "dns spec", minVersion("0.70.0"))
	same := b.spec("iam", "1", "iam spec", minVersion("0.66.0"))
	b.publish(t, inviteA, newer, same)

	result, err := Sync(context.Background(), cfg, inviteA)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam"}, result.Installed)
	require.Len(t, result.Skipped, 1)
	assert.Equal(t, "dns", result.Skipped[0].Module)
	assert.Contains(t, result.Skipped[0].Reason, "0.70.0")

	printed, err := json.Marshal(result)
	require.NoError(t, err)
	assert.JSONEq(t, `{"installed": ["iam"]}`, string(printed))

	cfg.CLIVersion = "dev"
	result, err = Sync(context.Background(), cfg, inviteA)
	require.NoError(t, err)
	assert.Equal(t, []string{"dns", "iam"}, result.Installed)
}

func TestSyncRefusesBadInput(t *testing.T) {
	_, cfg := newBucket(t)

	_, err := Sync(context.Background(), cfg, "my-invite")
	assert.ErrorContains(t, err, "not an alpha invite")

	cfg.BaseURL = ""
	_, err = Sync(context.Background(), cfg, inviteA)
	assert.ErrorContains(t, err, "no alpha source")
}

func TestResolveInvite(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, inviteA, b.spec("iam", "1", "iam spec"))

	_, err := resolveInvite(cfg.Dir, "")
	assert.ErrorContains(t, err, "include your invite")

	_, err = Sync(context.Background(), cfg, inviteA)
	require.NoError(t, err)

	invite, err := resolveInvite(cfg.Dir, "")
	require.NoError(t, err)
	assert.Equal(t, inviteA, invite, "the invite of the last sync")

	invite, err = resolveInvite(cfg.Dir, inviteB)
	require.NoError(t, err)
	assert.Equal(t, inviteB, invite, "the given one wins")
}

func TestSyncErrorsLeaveTheInviteOut(t *testing.T) {
	b, cfg := newBucket(t)
	b.status[ManifestPath(inviteA)] = http.StatusInternalServerError
	_, err := Sync(context.Background(), cfg, inviteA)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), inviteA)

	b.status = map[string]int{}
	b.files[ManifestPath(inviteA)] = "not json"
	_, err = Sync(context.Background(), cfg, inviteA)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), inviteA)
}

func TestSignedPayloadCoversWhatMatters(t *testing.T) {
	base := ManifestSpec{Module: "iam", Path: "a", SHA256: "ABC", MinCLIVersion: "0.66.0", KeyID: "k"}

	moved := base
	moved.Path = "b"
	assert.Equal(t, SignedPayload(base), SignedPayload(moved), "the path is not signed")

	for name, change := range map[string]func(*ManifestSpec){
		"module":          func(s *ManifestSpec) { s.Module = "dns" },
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

func TestVerifiedStoreChecksEveryRead(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, inviteA, b.spec("iam", "1", "iam spec"))
	_, err := Sync(context.Background(), cfg, inviteA)
	require.NoError(t, err)

	store := VerifiedStore(fsstore.NewDir(cfg.Dir), cfg.Dir, cfg.Keys)

	pkg, err := store.Read("iam")
	require.NoError(t, err)
	assert.Equal(t, "iam spec", string(pkg.Spec))

	require.NoError(t, os.WriteFile(filepath.Join(cfg.Dir, "iam.openapi.yaml"), []byte("changed"), 0o644))
	_, err = store.Read("iam")
	assert.ErrorContains(t, err, "changed since it was synced")

	require.NoError(t, os.WriteFile(filepath.Join(cfg.Dir, "dns.openapi.yaml"), []byte("dns spec"), 0o644))
	_, err = store.Read("dns")
	assert.ErrorContains(t, err, "was not installed by 'mgc alpha sync'")

	_, err = VerifiedStore(fsstore.NewDir(cfg.Dir), cfg.Dir, Keys{}).Read("iam")
	assert.ErrorContains(t, err, "does not know")
}
