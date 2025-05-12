package common

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	"github.com/MagaluCloud/magalu/mgc/core/progress_report"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

type bigFileCopier struct {
	cfg              Config
	src              mgcSchemaPkg.URI
	dst              mgcSchemaPkg.URI
	fileSize         int64
	totalParts       int
	uploadId         string
	progressReporter *progress_report.BytesReporter
	version          string
	storageClass     string
}

var _ copier = (*bigFileCopier)(nil)

type uploadPartCopyRequestResponse struct {
	ETag string `xml:"ETag"`
}

func (u *bigFileCopier) newPreparationRequest(ctx context.Context) (*http.Request, error) {
	req, err := newUploadRequest(ctx, u.cfg, u.dst, nil)
	if err != nil {
		return nil, err
	}
	req.Method = http.MethodPost
	req.Header.Set("Content-Type", "application/octet-stream")
	q := req.URL.Query()
	q.Set("uploads", "")

	if u.version != "" {
		q.Set("versionId", u.version)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (u *bigFileCopier) getUploadId(ctx context.Context) (string, error) {
	if u.uploadId == "" {
		req, err := u.newPreparationRequest(ctx)
		if err != nil {
			return "", err
		}

		resp, err := SendRequest(ctx, req, u.cfg)
		if err != nil {
			return "", err
		}

		pr, err := UnwrapResponse[preparationResponse](resp, req)
		if err != nil {
			return "", err
		}

		u.uploadId = pr.UploadId
	}
	return u.uploadId, nil
}

func (u *bigFileCopier) createMultipartRequest(ctx context.Context, partNumber int, startOffset int64, endOffset int64) (*http.Request, error) {
	uploadId, err := u.getUploadId(ctx)
	if err != nil {
		return nil, err
	}
	req, err := newCopyRequest(ctx, u.cfg, u.src, u.dst, u.version)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("uploadId", uploadId)
	q.Set("partNumber", fmt.Sprint(partNumber))
	req.URL.RawQuery = q.Encode()

	downloadByteRange := fmt.Sprintf("bytes=%d-%d", startOffset, endOffset)
	req.Header.Set("x-amz-copy-source-range", downloadByteRange)

	return req, nil
}

func (u *bigFileCopier) sendCompletionRequest(ctx context.Context, parts []completionPart, uploadId string) (err error) {
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].PartNumber < parts[j].PartNumber
	})
	body := completionRequest{
		Parts:     parts,
		Namespace: "http://s3.amazonaws.com/doc/2006-03-01/",
	}
	parsed, err := xml.Marshal(body)
	if err != nil {
		return fmt.Errorf("unable to marshal completion request body: %w\nCopy parts requests were successful but copy won't be finalized", err)
	}

	bigfileUploaderLogger().Debugw("All file parts uploaded, sending completion", "etags", parts)

	readerFunc := func() (io.ReadCloser, error) {
		reader := bytes.NewReader(parsed)
		return io.NopCloser(reader), nil
	}

	req, err := newUploadRequest(ctx, u.cfg, u.dst, readerFunc)
	if err != nil {
		return err
	}
	req.Method = http.MethodPost
	q := req.URL.Query()
	q.Set("uploadId", uploadId)
	req.URL.RawQuery = q.Encode()

	resp, err := SendRequestWithIgnoredHeaders(ctx, req, u.cfg, bigFileCopierExcludedHeaders)
	if err != nil {
		return err
	}

	err = ExtractErr(resp, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *bigFileCopier) createPartSenderProcessor(cancel context.CancelCauseFunc, uploadId string) pipeline.Processor[pipeline.WriteableChunk, completionPart] {
	return func(ctx context.Context, chunk pipeline.WriteableChunk) (part completionPart, status pipeline.ProcessStatus) {
		var err error
		defer func() { u.progressReporter.Report(0, err) }()

		partNumber := int(chunk.StartOffset/int64(u.cfg.chunkSizeInBytes())) + 1
		req, err := u.createMultipartRequest(ctx, partNumber, chunk.StartOffset, chunk.EndOffset)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		bigfileUploaderLogger().Debugw("Sending part", "part", partNumber, "total", u.totalParts)
		res, err := SendRequest(ctx, req, u.cfg)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		u.progressReporter.Report(uint64(chunk.EndOffset-chunk.StartOffset), err)

		err = ExtractErr(res, req)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		result, err := UnwrapResponse[uploadPartCopyRequestResponse](res, req)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		return NewCompletionPart(partNumber, result.ETag), pipeline.ProcessOutput
	}
}

func (u *bigFileCopier) Copy(ctx context.Context) error {
	bigfileUploaderLogger().Debug("start")

	name := "Preparing to copy " + u.src.String()
	u.progressReporter = progress_report.NewBytesReporter(ctx, name, uint64(u.fileSize))
	u.progressReporter.Start()
	defer u.progressReporter.End()
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	uploadId, err := u.getUploadId(ctx)
	if err != nil {
		return err
	}

	chunkChan := pipeline.PrepareWriteChunks(ctx, nil, u.fileSize, int64(u.cfg.chunkSizeInBytes()))
	partChan := pipeline.ParallelProcess(ctx, u.cfg.Workers, chunkChan, u.createPartSenderProcessor(cancel, uploadId), nil)

	parts, err := pipeline.SliceItemConsumer[[]completionPart](ctx, partChan)
	if err != nil {
		return err
	}

	return u.sendCompletionRequest(ctx, parts, uploadId)
}
