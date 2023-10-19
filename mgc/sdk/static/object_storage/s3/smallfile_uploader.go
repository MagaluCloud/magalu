package s3

import (
	"context"
	"io"
)

type smallFileUploader struct {
	cfg      Config
	dst      string
	mimeType string
	reader   io.Reader
}

var _ uploader = (*smallFileUploader)(nil)

func (u *smallFileUploader) Upload(ctx context.Context) error {
	req, err := newUploadRequest(ctx, u.cfg, u.dst, u.reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", u.mimeType)

	_, _, err = SendRequest[any](ctx, req)
	if err != nil {
		return err
	}
	return nil
}
