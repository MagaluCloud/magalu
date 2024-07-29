package common

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
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
	progressReporter := progress_report.NewBytesReporter(ctx, "Uploading "+u.fileInfo.Name(), uint64(u.fileInfo.Size()))
	progressReporter.Start()
	defer progressReporter.End()

	ctx = progress_report.NewBytesReporterContext(ctx, progressReporter)
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
	defer func() { progressReporter.Report(0, err) }()

	req, err := newUploadRequest(ctx, u.cfg, u.dst, newReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", u.mimeType)

	if u.storageClass != "" {
		req.Header.Set("X-Amz-Storage-Class", u.storageClass)
	}

	resp, err := SendRequest(ctx, req)
	if err != nil {
		return err
	}

	err = ExtractErr(resp, req)
	return err
}
