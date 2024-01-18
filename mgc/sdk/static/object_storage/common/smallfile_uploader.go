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

func (u *smallFileUploader) createProgressReporter(ctx context.Context) progress_report.ReportRead {
	reportProgress := progress_report.FromContext(ctx)
	fileName := u.fileInfo.Name()
	total := uint64(u.fileInfo.Size())
	sentBytes := uint64(0)
	return func(n int, err error) {
		sentBytes += uint64(n)
		reportProgress(fileName, max(0, sentBytes-1), total, progress_report.UnitsBytes, err)
	}
}

func (u *smallFileUploader) Upload(ctx context.Context) error {
	progressReporter := u.createProgressReporter(ctx)
	wrappedReader := progress_report.NewReporterReader(u.reader, progressReporter)
	req, err := newUploadRequest(ctx, u.cfg, u.dst, wrappedReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", u.mimeType)

	_, err = SendRequest(ctx, req)
	if err != nil {
		progressReporter(0, err)
		return err
	}
	progressReporter(1, nil)
	return nil
}
