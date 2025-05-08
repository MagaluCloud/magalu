package common

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	"github.com/MagaluCloud/magalu/mgc/core/progress_report"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type bigFileDownloader struct {
	cfg              Config
	src              mgcSchemaPkg.URI
	dst              mgcSchemaPkg.FilePath
	version          string
	fileSize         int64
	progressReporter *progress_report.BytesReporter
}

func (u *bigFileDownloader) createPartDownloaderProcessor(cancel context.CancelCauseFunc, cfg Config) pipeline.Processor[pipeline.WriteableChunk, error] {
	return func(ctx context.Context, chunk pipeline.WriteableChunk) (error, pipeline.ProcessStatus) {
		req, err := NewDownloadRequest(ctx, cfg, u.src, u.version)
		if err != nil {
			cancel(err)
			return err, pipeline.ProcessAbort
		}

		downloadByteRange := fmt.Sprintf("bytes=%d-%d", chunk.StartOffset, chunk.EndOffset)
		req.Header.Set("Range", downloadByteRange)

		resp, err := SendRequest(ctx, req, cfg)
		if err != nil {
			cancel(err)
			return err, pipeline.ProcessAbort
		}

		err = ExtractErr(resp, req)
		if err != nil {
			cancel(err)
			return err, pipeline.ProcessAbort
		}

		reporterWriter := progress_report.NewReporterWriter(chunk.Writer, u.progressReporter.Report)

		_, err = io.Copy(reporterWriter, resp.Body)
		if err != nil {
			return err, pipeline.ProcessAbort
		}

		return nil, pipeline.ProcessOutput
	}
}

func (u *bigFileDownloader) Download(ctx context.Context) error {
	u.progressReporter = progress_report.NewBytesReporter(ctx, fmt.Sprintf("Downloading %q", u.src), uint64(u.fileSize))
	u.progressReporter.Start()
	defer u.progressReporter.End()

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	dir := path.Dir(u.dst.String())
	if len(dir) != 0 {
		if err := os.MkdirAll(dir, utils.DIR_PERMISSION); err != nil {
			return err
		}
	}
	writer, err := os.OpenFile(u.dst.String(), os.O_WRONLY|os.O_CREATE, utils.FILE_PERMISSION)
	if err != nil {
		return err
	}

	chunkChan := pipeline.PrepareWriteChunks(ctx, writer, u.fileSize, int64(u.cfg.chunkSizeInBytes()))

	bigFileDownloadErrorChan := pipeline.ParallelProcess(ctx, u.cfg.Workers, chunkChan, u.createPartDownloaderProcessor(cancel, u.cfg), nil)
	bigFileDownloadErrorChan = pipeline.Filter(ctx, bigFileDownloadErrorChan, pipeline.FilterNonNil[error]{})

	objErr, _ := pipeline.SliceItemConsumer[utils.MultiError](ctx, bigFileDownloadErrorChan)
	if len(objErr) > 0 {
		return objErr
	}

	return nil
}
