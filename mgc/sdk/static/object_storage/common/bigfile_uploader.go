package common

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/http"
	"sort"

	"go.uber.org/zap"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

var deleteBucketsLogger *zap.SugaredLogger

func bigfileUploaderLogger() *zap.SugaredLogger {
	if deleteBucketsLogger == nil {
		deleteBucketsLogger = logger().Named("bigfileUploader")
	}
	return deleteBucketsLogger
}

type preparationResponse struct {
	XMLName  xml.Name `xml:"InitiateMultipartUploadResult"`
	Bucket   string
	Key      string
	UploadId string
}

type completionPart struct {
	Etag       string `xml:",innerxml"`
	PartNumber int
}

func NewCompletionPart(partNumber int, etag string) completionPart {
	return completionPart{
		// Manual string is needed here because content has double quotes.
		// Using xml.Marshal normally, the quotes are escaped, which cannot happen
		// Using `xml:",innerxml"` in struct solves this but removes tags
		Etag:       fmt.Sprintf("<ETag>%s</ETag>", etag),
		PartNumber: partNumber,
	}
}

type completionRequest struct {
	XMLName   xml.Name         `xml:"CompleteMultipartUpload"`
	Namespace string           `xml:"xmlns,attr"`
	Parts     []completionPart `xml:"Part"`
}

type bigFileUploader struct {
	cfg          Config
	dst          mgcSchemaPkg.URI
	mimeType     string
	fileInfo     fs.FileInfo
	filePath     mgcSchemaPkg.FilePath
	workerN      int
	uploadId     string
	storageClass string
}

var _ uploader = (*bigFileUploader)(nil)

func (u *bigFileUploader) newPreparationRequest(ctx context.Context) (*http.Request, error) {
	req, err := newUploadRequest(ctx, u.cfg, u.dst, nil)
	if err != nil {
		return nil, err
	}
	req.Method = http.MethodPost
	req.Header.Set("Content-Type", "application/octet-stream")

	if u.storageClass != "" {
		req.Header.Set("X-Amz-Storage-Class", u.storageClass)
	}

	q := req.URL.Query()
	q.Set("uploads", "")
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (u *bigFileUploader) getUploadId(ctx context.Context) (string, error) {
	if u.uploadId == "" {
		req, err := u.newPreparationRequest(ctx)
		if err != nil {
			return "", err
		}

		resp, err := SendRequest(ctx, req)
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

func (u *bigFileUploader) createMultipartRequest(ctx context.Context, partNumber int, body func() (io.ReadCloser, error)) (*http.Request, error) {
	uploadId, err := u.getUploadId(ctx)
	if err != nil {
		return nil, err
	}
	req, err := newUploadRequest(ctx, u.cfg, u.dst, body)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("uploadId", uploadId)
	q.Set("partNumber", fmt.Sprint(partNumber))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", u.mimeType)

	return req, nil
}

func (u *bigFileUploader) sendCompletionRequest(ctx context.Context, parts []completionPart, uploadId string) (err error) {
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].PartNumber < parts[j].PartNumber
	})
	body := completionRequest{
		Parts:     parts,
		Namespace: "http://s3.amazonaws.com/doc/2006-03-01/",
	}
	parsed, err := xml.Marshal(body)
	if err != nil {
		return err
	}

	bigfileUploaderLogger().Infow("All file parts uploaded, sending completion", "etags", parts)

	newReader := func() (io.ReadCloser, error) {
		reader := bytes.NewReader(parsed)
		return io.NopCloser(reader), nil
	}
	req, err := newUploadRequest(ctx, u.cfg, u.dst, newReader)
	if err != nil {
		return err
	}
	req.Method = http.MethodPost
	q := req.URL.Query()
	q.Set("uploadId", uploadId)
	req.URL.RawQuery = q.Encode()

	resp, err := SendRequestWithIgnoredHeaders(ctx, req, bigFileCopierExcludedHeaders)
	if err != nil {
		return err
	}

	err = ExtractErr(resp, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *bigFileUploader) createPartSenderProcessor(cancel context.CancelCauseFunc, totalParts int, uploadId string) pipeline.Processor[pipeline.ReadableChunk, completionPart] {
	return func(ctx context.Context, chunk pipeline.ReadableChunk) (part completionPart, status pipeline.ProcessStatus) {
		var err error

		newReader := func() (io.ReadCloser, error) {
			return io.NopCloser(io.NewSectionReader(chunk.Reader, 0, int64(u.cfg.chunkSizeInBytes()))), nil
		}

		partNumber := int(chunk.StartOffset/int64(u.cfg.chunkSizeInBytes())) + 1
		req, err := u.createMultipartRequest(ctx, partNumber, newReader)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		bigfileUploaderLogger().Debugw("Sending part", "part", partNumber, "total", totalParts)
		res, err := SendRequest(ctx, req)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		err = ExtractErr(res, req)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		return NewCompletionPart(partNumber, res.Header.Get("etag")), pipeline.ProcessOutput
	}
}

func (u *bigFileUploader) Upload(ctx context.Context) error {
	bigfileUploaderLogger().Debug("start")

	progressReportMsg := fmt.Sprintf("Uploading %q", u.fileInfo.Name())
	progressReporter := progress_report.NewBytesReporter(ctx, progressReportMsg, uint64(u.fileInfo.Size()))
	progressReporter.Start()
	defer progressReporter.End()
	ctx = progress_report.NewBytesReporterContext(ctx, progressReporter)

	ctx, cancel := context.WithCancelCause(ctx)

	var err error
	defer func() {
		if err == nil {
			err = ctx.Err()
		}

		progressReporter.Report(0, err)
		cancel(err)
	}()

	uploadId, err := u.getUploadId(ctx)
	if err != nil {
		return err
	}

	reader, err := readContent(u.filePath, u.fileInfo)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	totalParts := int(math.Ceil(float64(u.fileInfo.Size()) / float64(u.cfg.chunkSizeInBytes())))
	chunkChan := pipeline.ReadChunks(ctx, reader, u.fileInfo.Size(), int64(u.cfg.chunkSizeInBytes()))

	partChan := pipeline.ParallelProcess(ctx, u.workerN, chunkChan, u.createPartSenderProcessor(cancel, totalParts, uploadId), nil)

	parts, err := pipeline.SliceItemConsumer[[]completionPart](ctx, partChan)
	if err != nil {
		return err
	}

	return u.sendCompletionRequest(ctx, parts, uploadId)
}
