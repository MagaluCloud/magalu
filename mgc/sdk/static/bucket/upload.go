package bucket

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/s3"
)

type uploadParams struct {
	Destination string `json:"dst" jsonschema:"description=Full destination path in the bucket with s3 prefix, i.e: s3://<bucket-name>"`
	Source      string `json:"src" jsonschema:"description=Source file path to be uploaded"`
}

type uploadTemplateResult struct {
	File string `json:"file"`
	URI  string `json:"uri"`
}

func newUpload() core.Executor {
	executor := core.NewStaticExecute(
		"upload",
		"",
		"Upload a file to a bucket",
		upload,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Value) string {
		return "template=Uploaded file {{.file}} to {{.uri}}\n"
	})
}

func newUploadRequest(ctx context.Context, region, dst string, body []byte) (*http.Request, error) {
	host := s3.BuildHost(region)
	url, err := url.JoinPath(host, dst)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
}

func readContent(path string) ([]byte, error) {
	file, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	switch mode := file.Mode(); {
	case mode&os.ModeSymlink != 0:
		resolvedPath, err := os.Readlink(path)
		if err != nil {
			return nil, err
		}
		return os.ReadFile(resolvedPath)
	case mode.IsRegular():
		return os.ReadFile(path)
	default:
		// TODO: treat directory recursively
		return nil, fmt.Errorf("file type %s not supported", mode.Type())
	}
}

func formatURI(bucketURI, filename string) string {
	if bucketURI[len(bucketURI)-1] == '/' {
		return bucketURI + filename
	}

	return bucketURI + "/" + filename
}

func upload(ctx context.Context, params uploadParams, cfg s3.Config) (*uploadTemplateResult, error) {
	bucketURI, _ := strings.CutPrefix(params.Destination, s3.URIPrefix)
	_, fileName := path.Split(params.Source)

	fileContent, err := readContent(params.Source)
	if err != nil {
		return nil, fmt.Errorf("error reading object: %w", err)
	}
	req, err := newUploadRequest(ctx, cfg.Region, path.Join(bucketURI, fileName), fileContent)
	if err != nil {
		return nil, err
	}

	_, err = s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, nil)
	if err != nil {
		return nil, err
	}

	return &uploadTemplateResult{
		// path.Join will remove URI double slash prefix on s3://
		URI:  formatURI(params.Destination, fileName),
		File: fileName,
	}, nil
}
