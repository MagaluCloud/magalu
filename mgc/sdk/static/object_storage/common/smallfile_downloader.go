package common

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/MagaluCloud/magalu/mgc/core/progress_report"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type smallFileDownloader struct {
	cfg     Config
	src     mgcSchemaPkg.URI
	dst     mgcSchemaPkg.FilePath
	version string
}

var _ downloader = (*smallFileDownloader)(nil)

func (u *smallFileDownloader) Download(ctx context.Context) error {
	req, err := NewDownloadRequest(ctx, u.cfg, u.src, u.version)
	if err != nil {
		return err
	}

	resp, err := SendRequest(ctx, req, u.cfg)
	if err != nil {
		return err
	}

	err = ExtractErr(resp, req)
	if err != nil {
		return err
	}

	progressReporter := progress_report.NewBytesReporter(ctx, fmt.Sprintf("Downloading %q", u.src), uint64(resp.ContentLength))
	progressReporter.Start()
	defer progressReporter.End()

	resp.Body = progress_report.NewReporterReader(resp.Body, progressReporter.Report)

	dir := path.Dir(u.dst.String())
	if len(dir) != 0 {
		if err := os.MkdirAll(dir, utils.DIR_PERMISSION); err != nil {
			return err
		}
	}

	if err := WriteToFile(ctx, resp.Body, resp.ContentLength, u.dst); err != nil {
		return err
	}

	return nil
}
