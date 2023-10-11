package objects

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/s3"
)

var downloadObjectsLogger *zap.SugaredLogger

func downloadLogger() *zap.SugaredLogger {
	if downloadObjectsLogger == nil {
		downloadObjectsLogger = logger().Named("download")
	}
	return downloadObjectsLogger
}

type downloadObjectsError struct {
	errorMap map[string]error
}

func (o downloadObjectsError) Error() string {
	var errorMsg string
	for file, err := range o.errorMap {
		errorMsg += fmt.Sprintf("%s - %s, ", file, err)
	}
	// Remove trailing `, `
	if len(errorMsg) != 0 {
		errorMsg = errorMsg[:len(errorMsg)-2]
	}
	return fmt.Sprintf("failed to download some objects from bucket: %s", errorMsg)
}

func (o downloadObjectsError) Add(uri string, err error) {
	o.errorMap[uri] = err
}

func (o downloadObjectsError) HasError() bool {
	return len(o.errorMap) != 0
}

func NewDownloadObjectsError() downloadObjectsError {
	return downloadObjectsError{
		errorMap: make(map[string]error),
	}
}

type downloadObjectParams struct {
	Source      string `json:"src" jsonschema:"description=Path of the object to be downloaded" example:"s3://bucket1/file1"`
	Destination string `json:"dst,omitempty" jsonschema:"description=Name of the file to be saved" example:"file1.txt"`
}

func newDownload() core.Executor {
	executor := core.NewStaticExecute(
		"download",
		"",
		"download an object from a bucket",
		download,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Downloaded from {{.src}} to {{.dst}}\n"
	})
}

func newDownloadRequest(ctx context.Context, cfg s3.Config, pathURIs ...string) (*http.Request, error) {
	host := s3.BuildHost(cfg)
	url, err := url.JoinPath(host, pathURIs...)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func writeToFile(reader io.ReadCloser, outFile string) (err error) {
	defer reader.Close()

	writer, err := os.OpenFile(outFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	n, err := writer.ReadFrom(reader)
	defer writer.Close()
	if err != nil {
		return fmt.Errorf("error writing to file (wrote %d bytes): %w", n, err)
	}
	return nil
}

func downloadSingleFile(ctx context.Context, cfg s3.Config, src, dst string) error {
	bucketURI, _ := strings.CutPrefix(src, s3.URIPrefix)
	downloadLogger().Infof("Downloading %s", bucketURI)
	req, err := newDownloadRequest(ctx, cfg, bucketURI)
	if err != nil {
		return err
	}

	closer, _, err := s3.SendRequest[io.ReadCloser](ctx, req, cfg.AccessKeyID, cfg.SecretKey)
	if err != nil {
		return err
	}

	if !isFilePath(dst) {
		_, fname := path.Split(bucketURI)
		dst = path.Join(dst, fname)
	}

	dir, _ := path.Split(dst)
	if len(dir) != 0 {
		if err := os.MkdirAll(dir, utils.DIR_PERMISSION); err != nil {
			return err
		}
	}

	if err := writeToFile(closer, dst); err != nil {
		return err
	}

	return nil
}

func downloadMultipleFiles(ctx context.Context, cfg s3.Config, src, dst string) error {
	bucketURI, _ := strings.CutPrefix(src, s3.URIPrefix)
	bucketName := strings.Split(bucketURI, "/")[0]
	objs, err := List(ctx, ListObjectsParams{Destination: bucketURI}, cfg)
	if err != nil {
		return err
	}

	objError := NewDownloadObjectsError()
	for _, obj := range objs.Contents {
		objURI := path.Join(bucketName, obj.Key)
		downloadLogger().Infof("Downloading %s", objURI)
		req, err := newDownloadRequest(ctx, cfg, objURI)
		if err != nil {
			objError.Add(objURI, err)
			continue
		}

		closer, _, err := s3.SendRequest[io.ReadCloser](ctx, req, cfg.AccessKeyID, cfg.SecretKey)
		if err != nil || closer == nil {
			objError.Add(objURI, err)
			continue
		}

		dir, _ := path.Split(obj.Key)
		if err := os.MkdirAll(path.Join(dst, dir), utils.DIR_PERMISSION); err != nil {
			objError.Add(objURI, err)
			continue
		}

		if err := writeToFile(closer, path.Join(dst, obj.Key)); err != nil {
			objError.Add(objURI, err)
			continue
		}
	}

	if objError.HasError() {
		return objError
	}

	return nil
}

func isFilePath(fpath string) bool {
	// TODO: find a better way to infer if it's a file - requesting metadata to s3?
	return path.Ext(fpath) != ""
}

func download(ctx context.Context, p downloadObjectParams, cfg s3.Config) (result core.Value, err error) {
	dst := p.Destination
	if dst == "" {
		dst, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("no destination specified and could not use local dir: %w", err)
		}
		_, fname := path.Split(p.Source)
		dst = path.Join(dst, fname)
	}
	if isFilePath(p.Source) {
		// User specified a directory, append the file name to it
		if !isFilePath(dst) {
			_, fname := path.Split(p.Source)
			dst = path.Join(dst, fname)
		}
		err = downloadSingleFile(ctx, cfg, p.Source, dst)
	} else {
		err = downloadMultipleFiles(ctx, cfg, p.Source, dst)
	}

	if err != nil {
		return nil, err
	}

	return downloadObjectParams{Source: p.Source, Destination: dst}, nil
}
