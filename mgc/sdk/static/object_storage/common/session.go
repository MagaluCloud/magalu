package common

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
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
	r[http.CanonicalHeaderKey(contentMD5Header)] = struct{}{}
	return r
}()

type HostString string
type BucketHostString string

func BuildHost(cfg Config) HostString {
	var hostStr string
	if cfg.ServerUrl != "" {
		hostStr = cfg.ServerUrl
	} else {
		hostStr = strings.ReplaceAll(templateUrl, "{{region}}", cfg.translateRegion())
	}

	if hostStr[len(hostStr)-1] != '/' {
		hostStr += "/"
	}
	return HostString(hostStr)
}

func BuildHostURL(cfg Config) (*url.URL, error) {
	host := BuildHost(cfg)
	return url.Parse(string(host))
}

func BuildBucketHost(cfg Config, bucketName BucketName) (BucketHostString, error) {
	simpleHost := BuildHost(cfg)
	escapedBucketName := url.PathEscape(bucketName.String())
	host, err := url.JoinPath(string(simpleHost), escapedBucketName)
	if err != nil {
		return "", err
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
		err = fmt.Errorf("api-key not set, see how to set it with \"mgc object-storage api-key -h\"")
		return
	}

	if err = signHeaders(req, accesskeyId, accessSecretKey, unsignedPayload, ignoredHeaders); err != nil {
		return
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
