package common

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	corehttp "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
)

func TestCheckDeleteBatchResponse(t *testing.T) {
	for _, tc := range []struct{ name, body, wantError string }{
		{name: "confirmed", body: `<DeleteResult><Deleted><Key>a.txt</Key></Deleted><Deleted><Key>b.txt</Key></Deleted></DeleteResult>`},
		{name: "versioned", body: `<DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Deleted><Key>a.txt</Key><DeleteMarker>true</DeleteMarker><DeleteMarkerVersionId>v1</DeleteMarkerVersionId></Deleted><Deleted><Key>b.txt</Key><VersionId>null</VersionId></Deleted></DeleteResult>`},
		{name: "partial error", body: `<DeleteResult><Deleted><Key>a.txt</Key></Deleted><Error><Key>b.txt</Key><Code>AccessDenied</Code><Message>denied</Message></Error></DeleteResult>`, wantError: `failed to delete "b.txt": AccessDenied`},
		{name: "missing confirmation", body: `<DeleteResult><Deleted><Key>a.txt</Key></Deleted></DeleteResult>`, wantError: `missing deletion confirmation for "b.txt"`},
		{name: "duplicate confirmation", body: `<DeleteResult><Deleted><Key>a.txt</Key></Deleted><Deleted><Key>a.txt</Key></Deleted></DeleteResult>`, wantError: `unexpected deletion confirmation`},
		{name: "unexpected key", body: `<DeleteResult><Deleted><Key>other.txt</Key></Deleted></DeleteResult>`, wantError: `unexpected deletion confirmation`},
		{name: "empty body", wantError: "error decoding XML response body"},
		{name: "malformed XML", body: `<DeleteResult>`, wantError: "error decoding XML response body"},
		{name: "wrong XML root", body: `<Other/>`, wantError: "error decoding XML response body"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/xml"}}, Body: io.NopCloser(strings.NewReader(tc.body))}
			defer resp.Body.Close()
			req, err := http.NewRequest(http.MethodPost, "http://test.invalid/bucket?delete=", nil)
			if err != nil {
				t.Fatal(err)
			}
			result, err := UnwrapResponse[deleteBatchResponseBody](resp, req)
			if err == nil {
				err = checkDeleteBatchResponse(result, []objectIdentifier{{Key: "a.txt"}, {Key: "b.txt"}})
			}
			if tc.wantError == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantError) {
				t.Fatalf("error=%v, want %q", err, tc.wantError)
			}
		})
	}
}

type deleteTestTransport func(*http.Request) (*http.Response, error)

func (f deleteTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type deleteTestBody struct {
	io.Reader
	closed bool
}

func (b *deleteTestBody) Close() error { b.closed = true; return nil }

func TestDeletionProcessorChecksResponse(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		wantError  bool
	}{
		{name: "success", status: 200, body: `<DeleteResult><Deleted><Key>a.txt</Key></Deleted></DeleteResult>`},
		{name: "object error", status: 200, body: `<DeleteResult><Error><Key>a.txt</Key><Code>AccessDenied</Code><Message>denied</Message></Error></DeleteResult>`, wantError: true},
		{name: "HTTP error", status: 403, body: `<Error><Code>AccessDenied</Code><Message>denied</Message></Error>`, wantError: true},
		{name: "invalid response", status: 200, body: `not XML`, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &deleteTestBody{Reader: strings.NewReader(tc.body)}
			client := corehttp.NewClient(deleteTestTransport(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodPost || !r.URL.Query().Has("delete") {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				return &http.Response{StatusCode: tc.status, Status: fmt.Sprint(tc.status), Header: http.Header{"Content-Type": []string{"application/xml"}}, Body: body}, nil
			}))
			pm, _ := profile_manager.NewInMemoryProfileManager()
			conf := config.New(pm)
			conf.AddTempKeyPair("apikey", "test-key", "test-secret")
			ctx := auth.NewContext(context.Background(), auth.New(nil, &client.Client, pm, conf))
			ctx = corehttp.NewClientContext(ctx, client)
			cfg := Config{Workers: 1, Region: "br-se1", NetworkConfig: config.NetworkConfig{ServerUrl: "http://test.invalid"}}
			entries := []pipeline.WalkDirEntry{pipeline.NewSimpleWalkDirEntry("a.txt", &BucketContent{Key: "a.txt"}, nil)}
			err, status := createObjectDeletionProcessor(cfg, "bucket", nil)(ctx, entries)
			if (err != nil) != tc.wantError {
				t.Fatalf("error=%v, wantError=%t", err, tc.wantError)
			}
			if status != pipeline.ProcessOutput {
				t.Fatalf("result would be discarded: status=%v", status)
			}
			if !body.closed {
				t.Fatal("response body not closed")
			}
		})
	}
}

func TestDeleteObjectsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	entries := make(chan pipeline.WalkDirEntry)
	close(entries)
	err := DeleteObjects(ctx, DeleteObjectsParams{Destination: "s3://bucket/", ToDelete: entries, BatchSize: 1000}, Config{Workers: 2})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v, want canceled", err)
	}
}
