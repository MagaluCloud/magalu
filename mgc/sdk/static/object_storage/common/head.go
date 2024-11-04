package common

import (
	"context"
	"net/http"
	"strconv"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

type HeadObjectResponse struct {
	AcceptRanges  string
	LastModified  string
	ContentLength int64
	ETag          string
	ContentType   string
	StorageClass  string
}

func newHeadRequest(ctx context.Context, cfg Config, dst mgcSchemaPkg.URI, version string) (*http.Request, error) {
	host, err := BuildBucketHostWithPath(cfg, NewBucketNameFromURI(dst), dst.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, string(host), nil)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	if version != "" {
		query := req.URL.Query()
		query.Set("versionId", version)
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

func HeadFile(ctx context.Context, cfg Config, dst mgcSchemaPkg.URI, version string) (metadata HeadObjectResponse, err error) {
	req, err := newHeadRequest(ctx, cfg, dst, version)
	if err != nil {
		return
	}

	resp, err := SendRequest(ctx, req)
	if err != nil {
		return
	}

	err = ExtractErr(resp, req)
	if err != nil {
		return
	}

	metadata, err = getMetadataFromResponse(resp)
	if err != nil {
		return
	}

	return
}

func parseContentLength(contentLenght string) (int64, error) {
	value, err := strconv.ParseInt(contentLenght, HeadContentLengthBase, HeadContentLengthBitSize)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func getMetadataFromResponse(resp *http.Response) (HeadObjectResponse, error) {
	contentLength, err := parseContentLength(resp.Header.Get("Content-Length"))
	if err != nil {
		return HeadObjectResponse{}, err
	}

	metadata := HeadObjectResponse{
		AcceptRanges:  resp.Header.Get("Accept-Ranges"),
		LastModified:  resp.Header.Get("Last-Modified"),
		ContentLength: contentLength,
		ETag:          resp.Header.Get("ETag"),
		ContentType:   resp.Header.Get("Content-Type"),
		StorageClass:  resp.Header.Get("x-amz-storage-class"),
	}

	return metadata, nil
}
