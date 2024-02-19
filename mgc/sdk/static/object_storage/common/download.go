package common

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type DownloadObjectParams struct {
	Source      mgcSchemaPkg.URI      `json:"src" jsonschema:"description=Path of the object to be downloaded,example=bucket1/file.txt" mgc:"positional"`
	Destination mgcSchemaPkg.FilePath `json:"dst,omitempty" jsonschema:"description=Path and file name to be saved (relative or absolute).If not specified it defaults to the current working directory,example=file.txt" mgc:"positional"`
}

type downloader interface {
	Download(context.Context) error
}

func NewDownloadRequest(ctx context.Context, cfg Config, src mgcSchemaPkg.URI) (*http.Request, error) {
	host, err := BuildBucketHostWithPath(cfg, NewBucketNameFromURI(src), src.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, string(host), nil)
}

func WriteToFile(ctx context.Context, reader io.ReadCloser, fileSize int64, outFile mgcSchemaPkg.FilePath) (err error) {
	defer reader.Close()

	writer, err := os.OpenFile(outFile.String(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, utils.FILE_PERMISSION)
	if err != nil {
		return err
	}

	n, err := io.Copy(writer, reader)
	defer writer.Close()
	if err != nil {
		return fmt.Errorf("error writing to file (wrote %d bytes): %w", n, err)
	}
	return nil
}

func NewDownloader(ctx context.Context, cfg Config, src mgcSchemaPkg.URI, dst mgcSchemaPkg.FilePath) (downloader, error) {
	metadata, err := HeadFile(ctx, cfg, src)
	if err != nil {
		return nil, err
	}

	totalDownloadParts := int(math.Ceil(float64(metadata.ContentLength) / float64(CHUNK_SIZE)))

	if totalDownloadParts > 1 {
		return &bigFileDownloader{
			cfg:      cfg,
			src:      src,
			dst:      dst,
			fileSize: metadata.ContentLength,
		}, nil
	} else {
		return &smallFileDownloader{
			cfg: cfg,
			src: src,
			dst: dst,
		}, nil
	}
}

func GetDestination(dst mgcSchemaPkg.FilePath, src mgcSchemaPkg.URI) (mgcSchemaPkg.FilePath, error) {
	if dst == "" {
		d, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("no destination specified and could not use local dir: %w", err)
		}
		_, fname := path.Split(src.Path())
		return mgcSchemaPkg.FilePath(path.Join(d, fname)), nil
	}
	return dst, nil
}
