package common

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core/auth"
	mgcHttpPkg "magalu.cloud/core/http"
)

var excludedHeaders = HeaderMap{
	"Authorization":         nil,
	"Accept-Encoding":       nil,
	"Amz-Sdk-Invocation-Id": nil,
	"Amz-Sdk-Request":       nil,
	"User-Agent":            nil,
	"X-Amzn-Trace-Id":       nil,
	"Expect":                nil,
	"Content-Length":        nil,
}

type HostString string
type BucketHostString string

func BuildHost(cfg Config) HostString {
	if cfg.ServerUrl != "" {
		return HostString(cfg.ServerUrl)
	}
	return HostString(strings.ReplaceAll(templateUrl, "{{region}}", cfg.Region))
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

func SendRequest(ctx context.Context, req *http.Request) (res *http.Response, err error) {
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

	if err = sign(req, accesskeyId, accessSecretKey, unsignedPayload, excludedHeaders); err != nil {
		return
	}

	res, err = httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error to send HTTP request: %w", err)
		return
	}

	return
}
