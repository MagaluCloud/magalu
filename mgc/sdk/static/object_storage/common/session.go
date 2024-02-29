package common

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/progress_report"
)

// note: must be in HTTP Header Canonical format (Title-Case-With-Dashes)
// as per https://pkg.go.dev/net/http#CanonicalHeaderKey as they should match https://pkg.go.dev/net/http#Headervar
var excludedHeaders = map[string]struct{}{
	http.CanonicalHeaderKey("Authorization"):         {},
	http.CanonicalHeaderKey("Accept-Encoding"):       {},
	http.CanonicalHeaderKey("Amz-Sdk-Invocation-Id"): {},
	http.CanonicalHeaderKey("Amz-Sdk-Request"):       {},
	http.CanonicalHeaderKey("User-Agent"):            {},
	http.CanonicalHeaderKey("X-Amzn-Trace-Id"):       {},
	http.CanonicalHeaderKey("Expect"):                {},
	http.CanonicalHeaderKey("Content-Length"):        {},
}

var bigFileCopierExcludedHeaders = func() map[string]struct{} {
	r := maps.Clone(excludedHeaders)
	r[http.CanonicalHeaderKey("Content-Type")] = struct{}{}
	r[http.CanonicalHeaderKey("Content-Md5")] = struct{}{}
	return r
}()

type HostString string
type BucketHostString string

func BuildHost(cfg Config) HostString {
	if cfg.ServerUrl != "" {
		return HostString(cfg.ServerUrl)
	}
	return HostString(strings.ReplaceAll(templateUrl, "{{region}}", cfg.translateRegion()))
}

func BuildHostURL(cfg Config) (*url.URL, error) {
	host := BuildHost(cfg)
	return url.Parse(string(host))
}

func BuildBucketHost(cfg Config, bucketName BucketName) (BucketHostString, error) {
	host, err := url.JoinPath(string(BuildHost(cfg)), bucketName.String())
	if err != nil {
		return BucketHostString(host), err
	}
	// Bucket URI cannot end in '/' as this makes it search for a
	// non existing directory
	host = strings.TrimSuffix(host, "/")
	return BucketHostString(host), nil
}

func BuildBucketHostWithPath(cfg Config, bucketName BucketName, path string) (BucketHostString, error) {
	bucketHost, err := BuildBucketHost(cfg, bucketName)
	if err != nil {
		return bucketHost, err
	}
	bucketHostWithPath, err := url.JoinPath(string(bucketHost), path)
	if err != nil {
		return BucketHostString(bucketHostWithPath), err
	}
	// Bucket URI cannot end in '/' as this makes it search for a
	// non existing directory
	bucketHostWithPath = strings.TrimSuffix(string(bucketHostWithPath), "/")
	return BucketHostString(bucketHostWithPath), err
}

func BuildBucketHostURL(cfg Config, bucketName BucketName) (*url.URL, error) {
	bucketHost, err := BuildBucketHost(cfg, bucketName)
	if err != nil {
		return nil, err
	}
	return url.Parse(string(bucketHost))
}

func BuildBucketHostWithPathURL(cfg Config, bucketName BucketName, path string) (*url.URL, error) {
	bucketHost, err := BuildBucketHostWithPath(cfg, bucketName, path)
	if err != nil {
		return nil, err
	}
	return url.Parse(string(bucketHost))
}

func SendRequestWithIgnoredHeaders(ctx context.Context, req *http.Request, ignoredHeaders map[string]struct{}) (res *http.Response, err error) {
	httpClient := mgcHttpPkg.ClientFromContext(ctx)
	if httpClient == nil {
		err = fmt.Errorf("couldn't get http client from context")
		return
	}

	var unsignedPayload bool
	if req.Method == http.MethodPut {
		unsignedPayload = true
	}

	accesskeyId, accessSecretKey := auth.FromContext(ctx).AccessKeyPair()
	if accesskeyId == "" || accessSecretKey == "" {
		err = fmt.Errorf("access key not set, see how to set it with \"auth object-storage set -h\"")
		return
	}

	if err = signHeaders(req, accesskeyId, accessSecretKey, unsignedPayload, ignoredHeaders); err != nil {
		return
	}

	if req.Body != nil {
		if reporter := progress_report.BytesReporterFromContext(ctx); reporter != nil {
			req.Body = progress_report.NewReporterReader(req.Body, reporter.Report)
			if getBodyRaw := req.GetBody; getBodyRaw != nil {
				req.GetBody = func() (io.ReadCloser, error) {
					body, err := getBodyRaw()
					if err != nil {
						return nil, err
					}
					return progress_report.NewReporterReader(body, reporter.Report), nil
				}
			}
		}
	}

	res, err = httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error to send HTTP request: %w", err)
		return
	}

	return
}

func SendRequest(ctx context.Context, req *http.Request) (res *http.Response, err error) {
	return SendRequestWithIgnoredHeaders(ctx, req, excludedHeaders)
}
