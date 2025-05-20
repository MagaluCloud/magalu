package common

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

type smallFileUploader struct {
	cfg          Config
	dst          mgcSchemaPkg.URI
	mimeType     string
	fileInfo     fs.FileInfo
	filePath     mgcSchemaPkg.FilePath
	storageClass string
}

var _ uploader = (*smallFileUploader)(nil)

func (u *smallFileUploader) Upload(ctx context.Context) error {

	newReader := func() (io.ReadCloser, error) {
		reader, err := readContent(u.filePath, u.fileInfo)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		return reader, nil
	}

	var err error
	// TODO: This will only work sometimes... sometimes the error won't be nil but it won't
	// be updated in the progress bar

	req, err := newUploadRequest(ctx, u.cfg, u.dst, newReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", u.mimeType)

	if u.storageClass != "" {
		req.Header.Set("X-Amz-Storage-Class", u.storageClass)
	}

	resp, err := SendRequest(ctx, req, u.cfg)
	if err != nil {
		return err
	}

	err = ExtractErr(resp, req)
	return err
}
