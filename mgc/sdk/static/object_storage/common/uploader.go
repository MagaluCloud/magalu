package common

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

type uploader interface {
	Upload(context.Context) error
}

func NewUploader(cfg Config, src mgcSchemaPkg.FilePath, dst mgcSchemaPkg.URI) (uploader, error) {
	fileInfo, err := os.Stat(src.String())
	if err != nil {
		return nil, fmt.Errorf("error reading object: %w", err)
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("cannot upload a directory, use 'upload-dir' instead")
	}

	size := fileInfo.Size()
	mimeType := mime.TypeByExtension(filepath.Ext(fileInfo.Name()))

	chunkN := int(math.Ceil(float64(size) / float64(cfg.chunkSizeInBytes())))

	if chunkN > 1 {
		return &bigFileUploader{
			cfg:      cfg,
			dst:      dst,
			mimeType: mimeType,
			fileInfo: fileInfo,
			filePath: src,
			workerN:  cfg.Workers,
		}, nil
	} else {
		return &smallFileUploader{
			cfg:      cfg,
			dst:      dst,
			mimeType: mimeType,
			fileInfo: fileInfo,
			filePath: src,
		}, nil
	}
}

func newUploadRequest(ctx context.Context, cfg Config, dst mgcSchemaPkg.URI, newReader func() (io.ReadCloser, error)) (*http.Request, error) {
	host, err := BuildBucketHostWithPath(cfg, NewBucketNameFromURI(dst), dst.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	var body io.ReadCloser

	if newReader != nil {
		body, err = newReader()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, string(host), body)
	if err != nil {
		return nil, err
	}
	req.GetBody = newReader
	return req, nil
}

func readContent(p mgcSchemaPkg.FilePath, file fs.FileInfo) (*os.File, error) {
	path := p.String()

	switch mode := file.Mode(); {
	case mode&os.ModeSymlink != 0:
		resolvedPath, err := os.Readlink(path)
		if err != nil {
			return nil, err
		}
		reader, err := os.Open(resolvedPath)
		return reader, err
	case mode.IsRegular():
		reader, err := os.Open(path)
		return reader, err
	default:
		return nil, fmt.Errorf("file type %s not supported", mode.Type())
	}
}
