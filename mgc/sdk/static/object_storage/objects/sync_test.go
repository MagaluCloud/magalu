package objects

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	syncer "sync"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

// fakeBucket answers just enough of the S3 API for sync: list, head (always
// 404), put and batch delete.
type fakeBucket struct {
	mu      syncer.Mutex
	keys    []string
	puts    []string
	deletes []string
}

func (f *fakeBucket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		var b strings.Builder
		b.WriteString("<ListBucketResult><Name>bucket</Name>")
		for _, k := range f.keys {
			fmt.Fprintf(&b, "<Contents><Key>%s</Key><Size>1</Size></Contents>", k)
		}
		b.WriteString("</ListBucketResult>")
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(b.String()))
	case http.MethodHead:
		w.WriteHeader(http.StatusNotFound)
	case http.MethodPut:
		_, _ = io.Copy(io.Discard, r.Body)
		f.puts = append(f.puts, r.URL.Path)
	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		f.deletes = append(f.deletes, string(body))
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte("<DeleteResult></DeleteResult>"))
	}
}

func newSyncTestEnv(t *testing.T, fake *fakeBucket, workers int) (context.Context, common.Config) {
	t.Helper()

	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)

	pm, _ := profile_manager.NewInMemoryProfileManager()
	mgcCfg := config.New(pm)
	mgcCfg.AddTempKeyPair("apikey", "id", "secret")

	ctx := auth.NewContext(context.Background(), auth.New(nil, nil, pm, mgcCfg))
	ctx = mgcHttpPkg.NewClientContext(ctx, mgcHttpPkg.NewClient(http.DefaultTransport))

	cfg := common.Config{Workers: workers, Region: "br-se1"}
	cfg.ServerUrl = srv.URL
	return ctx, cfg
}

func writeLocalFiles(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestSyncDeleteOnlyRemovesMissingFiles(t *testing.T) {
	tests := []struct {
		name     string
		bucket   mgcSchemaPkg.URI
		keys     []string
		expected string
	}{
		{
			name:     "bucket root",
			bucket:   "bucket",
			keys:     []string{"a.txt", "stale.txt"},
			expected: "<Delete><Object><Key>stale.txt</Key></Object></Delete>",
		},
		{
			name:     "prefix",
			bucket:   "bucket/dir/",
			keys:     []string{"dir/a.txt", "dir/stale.txt"},
			expected: "<Delete><Object><Key>dir/stale.txt</Key></Object></Delete>",
		},
		{
			name:     "nested prefix without trailing slash",
			bucket:   "bucket/dir/sub",
			keys:     []string{"dir/sub/a.txt", "dir/sub/stale.txt"},
			expected: "<Delete><Object><Key>dir/sub/stale.txt</Key></Object></Delete>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeBucket{keys: tt.keys}
			ctx, cfg := newSyncTestEnv(t, fake, 1)

			_, err := sync(ctx, syncParams{
				Local:     mgcSchemaPkg.URI(writeLocalFiles(t, "a.txt")),
				Bucket:    tt.bucket,
				Delete:    true,
				BatchSize: common.MaxBatchSize,
			}, cfg)
			if err != nil {
				t.Fatal(err)
			}

			if len(fake.deletes) != 1 || fake.deletes[0] != tt.expected {
				t.Errorf("expected delete request %q, got %q", tt.expected, fake.deletes)
			}
		})
	}
}

func TestSyncParallelWorkers(t *testing.T) {
	var names, keys []string
	for i := 0; i < 100; i++ {
		names = append(names, fmt.Sprintf("f%d.txt", i))
		keys = append(keys, fmt.Sprintf("f%d.txt", i))
	}

	fake := &fakeBucket{keys: keys}
	ctx, cfg := newSyncTestEnv(t, fake, 8)

	result, err := sync(ctx, syncParams{
		Local:     mgcSchemaPkg.URI(writeLocalFiles(t, names...)),
		Bucket:    "bucket",
		BatchSize: common.MaxBatchSize,
	}, cfg)
	if err != nil {
		t.Fatal(err)
	}

	res := result.(syncResult)
	if res.FilesUploaded != len(names) || res.FilesDeleted != 0 {
		t.Errorf("expected %d uploaded and 0 to delete, got %d and %d", len(names), res.FilesUploaded, res.FilesDeleted)
	}
}

func TestSyncMissingLocalPath(t *testing.T) {
	ctx, cfg := newSyncTestEnv(t, &fakeBucket{}, 1)

	_, err := sync(ctx, syncParams{
		Local:     mgcSchemaPkg.URI(filepath.Join(t.TempDir(), "missing")),
		Bucket:    "bucket",
		BatchSize: common.MaxBatchSize,
	}, cfg)
	if err == nil {
		t.Fatal("expected an error for a missing local path")
	}
}
