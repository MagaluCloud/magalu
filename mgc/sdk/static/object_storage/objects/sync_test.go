package objects

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	sy "sync"
	"testing"
	"time"

	"github.com/pterm/pterm"

	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	corehttp "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	schema "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/sdk/openapi"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type syncTestTransport func(*http.Request) (*http.Response, error)

func (f syncTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type syncTestStore struct {
	mu          sy.Mutex
	objects     map[string]string
	lists, puts int
	listError   bool
}

func (s *syncTestStore) serve(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	w.Header().Set("Content-Type", "application/xml")
	key := strings.TrimPrefix(r.URL.Path, "/bucket/")
	switch r.Method {
	case http.MethodGet:
		s.lists++
		if s.listError {
			http.Error(w, "list failed", http.StatusServiceUnavailable)
			return
		}
		q := r.URL.Query()
		keys := []string{}
		for k := range s.objects {
			if strings.HasPrefix(k, q.Get("prefix")) {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		start, _ := strconv.Atoi(q.Get("continuation-token"))
		limit, _ := strconv.Atoi(q.Get("max-keys"))
		end := min(start+limit, len(keys))
		fmt.Fprintf(w, "<ListBucketResult><IsTruncated>%t</IsTruncated><NextContinuationToken>%d</NextContinuationToken>", end < len(keys), end)
		for _, k := range keys[start:end] {
			fmt.Fprintf(w, "<Contents><Key>%s</Key><Size>%d</Size></Contents>", k, len(s.objects[k]))
		}
		fmt.Fprint(w, "</ListBucketResult>")
	case http.MethodHead:
		value, ok := s.objects[key]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(value)))
		w.Header().Set("Last-Modified", time.Now().Add(time.Hour).UTC().Format(http.TimeFormat))
	case http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		s.objects[key] = string(body)
		s.puts++
	default:
		http.Error(w, "unexpected request", http.StatusBadRequest)
	}
}

func syncTestContext(t *testing.T, store *syncTestStore) (context.Context, common.Config) {
	t.Helper()
	transport := syncTestTransport(func(r *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		store.serve(recorder, r)
		return recorder.Result(), nil
	})
	client := corehttp.NewClient(transport)
	pm, _ := profile_manager.NewInMemoryProfileManager()
	conf := config.New(pm)
	conf.AddTempKeyPair("apikey", "test-key", "test-secret")
	a := auth.New(nil, &client.Client, pm, conf)
	ctx := auth.NewContext(context.Background(), a)
	ctx = corehttp.NewClientContext(ctx, client)
	return ctx, common.Config{Workers: 2, ChunkSize: 8, Region: "br-se1", NetworkConfig: config.NetworkConfig{ServerUrl: "http://test.invalid"}}
}

func TestSyncPreservesLocalObjects(t *testing.T) {
	for _, destination := range []string{"bucket", "bucket/", "bucket/prefix", "bucket/prefix/", "s3://bucket/prefix/"} {
		for _, relative := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/relative=%t", destination, relative), func(t *testing.T) {
				dir := t.TempDir()
				if err := os.Mkdir(filepath.Join(dir, "nested"), 0700); err != nil {
					t.Fatal(err)
				}
				prefix := ""
				if strings.Contains(destination, "prefix") {
					prefix = "prefix/"
				}
				store := &syncTestStore{objects: map[string]string{}}
				for _, name := range []string{"a.txt", "b.txt", "nested/c.txt"} {
					store.objects[prefix+name] = "text"
					if name != "b.txt" {
						if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(name)), []byte("text"), 0600); err != nil {
							t.Fatal(err)
						}
					}
				}
				local := dir + string(os.PathSeparator)
				if relative {
					wd, err := os.Getwd()
					if err != nil {
						t.Fatal(err)
					}
					local, err = filepath.Rel(wd, dir)
					if err != nil {
						t.Fatal(err)
					}
				}
				ctx, cfg := syncTestContext(t, store)
				inventory, err := fillBucketFiles(ctx, syncParams{Bucket: schema.URI("s3://" + strings.TrimPrefix(destination, "s3://"))}, cfg)
				if err != nil {
					t.Fatal(err)
				}
				// Inspect the selection independently of the deletion pipeline.
				base, err := common.GetAbsSystemURI(schema.URI(local))
				if err != nil {
					t.Fatal(err)
				}
				files, err := walkDir(ctx, base.String(), false)
				if err != nil {
					t.Fatal(err)
				}
				state := &syncState{bucketFiles: inventory}
				progress := pterm.DefaultProgressbar.WithTotal(len(files)).WithRemoveWhenDone(true)
				err = processSyncFiles(ctx, cfg, schema.URI("s3://"+strings.TrimPrefix(destination, "s3://")), base.String(), files, progress, state)
				if err != nil {
					t.Fatal(err)
				}
				if len(state.bucketFiles) != 1 || !state.bucketFiles["/"+prefix+"b.txt"] {
					t.Fatalf("unexpected candidates: %v", state.bucketFiles)
				}
				if state.uploadFiles.Value() != 0 {
					t.Fatal("unchanged files were uploaded")
				}
			})
		}
	}
}

func TestSyncRepeatedInvocations(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("text"), 0600); err != nil {
		t.Fatal(err)
	}
	store := &syncTestStore{objects: map[string]string{}}
	ctx, cfg := syncTestContext(t, store)
	for _, want := range []int{1, 0} {
		value, err := sync(ctx, syncParams{Local: schema.URI(dir), Bucket: "bucket/prefix/", BatchSize: 1000}, cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got := value.(syncResult).FilesUploaded; got != want {
			t.Fatalf("uploaded=%d, want %d", got, want)
		}
	}
}

func TestSyncInventoryPagination(t *testing.T) {
	store := &syncTestStore{objects: map[string]string{}}
	for i := 0; i < 1001; i++ {
		store.objects[fmt.Sprintf("prefix/%04d.txt", i)] = "text"
	}
	ctx, cfg := syncTestContext(t, store)
	inventory, err := fillBucketFiles(ctx, syncParams{Bucket: "s3://bucket/prefix/"}, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory) != 1001 || store.lists != 2 {
		t.Fatalf("objects=%d requests=%d", len(inventory), store.lists)
	}
}

func TestSyncAbortsOnInventoryFailure(t *testing.T) {
	for _, transportError := range []bool{false, true} {
		t.Run(fmt.Sprint(transportError), func(t *testing.T) {
			store := &syncTestStore{objects: map[string]string{}, listError: true}
			ctx, cfg := syncTestContext(t, store)
			if transportError {
				ctx = corehttp.NewClientContext(ctx, corehttp.NewClient(syncTestTransport(func(*http.Request) (*http.Response, error) { return nil, fmt.Errorf("transport failed") })))
			}
			ctx = openapi.WithRawOutputFlag(ctx, false)
			before := len(pterm.ActiveProgressBarPrinters)
			t.Cleanup(func() {
				for _, bar := range pterm.ActiveProgressBarPrinters[before:] {
					_, _ = bar.Stop()
				}
				pterm.ActiveProgressBarPrinters = pterm.ActiveProgressBarPrinters[:before]
			})
			_, err := sync(ctx, syncParams{Local: schema.URI(t.TempDir()), Bucket: "bucket/prefix/", BatchSize: 1000, Delete: true}, cfg)
			if err == nil {
				t.Fatal("expected inventory failure")
			}
			if len(pterm.ActiveProgressBarPrinters) != before+1 || pterm.ActiveProgressBarPrinters[before].IsActive {
				t.Fatal("progress bar was not started and stopped on inventory failure")
			}
		})
	}
}

func TestSyncRejectsInvalidBatchBeforeRequests(t *testing.T) {
	for _, batch := range []int{-1, 0, 1001} {
		store := &syncTestStore{objects: map[string]string{}}
		ctx, cfg := syncTestContext(t, store)
		_, err := sync(ctx, syncParams{Local: schema.URI(t.TempDir()), Bucket: "bucket/", BatchSize: batch, Delete: true}, cfg)
		if err == nil || store.lists != 0 {
			t.Fatalf("batch=%d error=%v list requests=%d", batch, err, store.lists)
		}
	}
}
