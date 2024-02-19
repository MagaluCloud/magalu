package common

import (
	"context"
	"io"
	"io/fs"

	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

type smallFileUploader struct {
	cfg      Config
	dst      mgcSchemaPkg.URI
	mimeType string
	reader   io.Reader
	fileInfo fs.FileInfo
}

var _ uploader = (*smallFileUploader)(nil)

func (u *smallFileUploader) Upload(ctx context.Context) error {
	progressReporter := progress_report.NewBytesReporter(ctx, "Uploading "+u.fileInfo.Name(), uint64(u.fileInfo.Size()))
	progressReporter.Start()
	defer progressReporter.End()

	wrappedReader := progress_report.NewReporterReader(u.reader, progressReporter.Report)
	req, err := newUploadRequest(ctx, u.cfg, u.dst, wrappedReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", u.mimeType)

	resp, err := SendRequest(ctx, req)
	if err != nil {
		return err
	}

	return ExtractErr(resp, req)
}
