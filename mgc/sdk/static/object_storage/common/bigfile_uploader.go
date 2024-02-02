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
	"magalu.cloud/core/utils"
)

type progressReport struct {
	bytes uint64
	err   error
}

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
	cfg        Config
	dst        mgcSchemaPkg.URI
	mimeType   string
	reader     io.ReaderAt
	fileInfo   fs.FileInfo
	workerN    int
	uploadId   string
	reportChan chan progressReport
}

var _ uploader = (*bigFileUploader)(nil)

func (u *bigFileUploader) newPreparationRequest(ctx context.Context) (*http.Request, error) {
	req, err := newUploadRequest(ctx, u.cfg, u.dst, nil)
	if err != nil {
		return nil, err
	}
	req.Method = http.MethodPost
	req.Header.Set("Content-Type", "application/octet-stream")
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

		pr, err := UnwrapResponse[preparationResponse](resp)
		if err != nil {
			return "", err
		}

		u.uploadId = pr.UploadId
	}
	return u.uploadId, nil
}

func (u *bigFileUploader) createMultipartRequest(ctx context.Context, partNumber int, body io.Reader) (*http.Request, error) {
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
	defer func() { u.reportProgress(0, err) }()

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

	bigfileUploaderLogger().Debugw("All file parts uploaded, sending completion", "etags", parts)

	req, err := newUploadRequest(ctx, u.cfg, u.dst, bytes.NewReader(parsed))
	if err != nil {
		return err
	}
	req.Method = http.MethodPost
	q := req.URL.Query()
	q.Set("uploadId", uploadId)
	req.URL.RawQuery = q.Encode()

	// excludedHeaders is a global variable that needs to be altered specifically
	// for this request, so set the correct headers and resets after
	excludedHeaders["Content-Type"] = nil
	excludedHeaders["Content-MD5"] = nil
	defer func() {
		delete(excludedHeaders, "Content-Type")
		delete(excludedHeaders, "Content-MD5")
	}()

	resp, err := SendRequest(ctx, req)
	if err != nil {
		return err
	}

	err = ExtractErr(resp)
	if err != nil {
		return err
	}

	return nil
}

func (u *bigFileUploader) createPartSenderProcessor(cancel context.CancelCauseFunc, totalParts int, uploadId string) pipeline.Processor[pipeline.ReadableChunk, completionPart] {
	return func(ctx context.Context, chunk pipeline.ReadableChunk) (part completionPart, status pipeline.ProcessStatus) {
		var err error
		defer func() { u.reportProgress(1, err) }()

		partNumber := int(chunk.StartOffset/CHUNK_SIZE) + 1
		reader := progress_report.NewReporterReader(chunk.Reader, u.reportProgress)
		req, err := u.createMultipartRequest(ctx, partNumber, reader)
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

		err = ExtractErr(res)
		if err != nil {
			cancel(err)
			return part, pipeline.ProcessAbort
		}

		return NewCompletionPart(partNumber, res.Header.Get("etag")), pipeline.ProcessOutput
	}
}

func (u *bigFileUploader) Upload(ctx context.Context) error {
	bigfileUploaderLogger().Debug("start")

	reportProgress := progress_report.FromContext(ctx)
	u.reportChan = make(chan progressReport)
	defer close(u.reportChan)

	totalParts := int(math.Ceil(float64(u.fileInfo.Size()) / float64(CHUNK_SIZE)))
	chunkChan := pipeline.ReadChunks(ctx, u.reader, u.fileInfo.Size(), CHUNK_SIZE)

	go progressReportSubroutine(reportProgress, totalParts, u.reportChan, u.fileInfo)

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	uploadId, err := u.getUploadId(ctx)
	if err != nil {
		return err
	}

	partChan := pipeline.ParallelProcess(ctx, u.workerN, chunkChan, u.createPartSenderProcessor(cancel, totalParts, uploadId), nil)

	parts, err := pipeline.SliceItemConsumer[[]completionPart](ctx, partChan)
	if err != nil {
		return err
	}

	return u.sendCompletionRequest(ctx, parts, uploadId)
}

func (u *bigFileUploader) reportProgress(n int, err error) {
	if u.reportChan == nil {
		return
	}

	u.reportChan <- progressReport{bytes: uint64(n), err: err}
}

func progressReportSubroutine(
	reportProgress progress_report.ReportProgress,
	totalParts int,
	reportChan <-chan progressReport,
	fileInfo fs.FileInfo,
) {
	reportMsg := "Upload " + fileInfo.Name()
	total := uint64(totalParts)
	progress := uint64(0)

	// total here must be reported as one, otherwise the progress-bar shows
	// an animation we do not wish the user to see
	reportProgress(reportMsg, progress, total, progress_report.UnitsNone, nil)

	var errors utils.MultiError
	for report := range reportChan {
		progress += report.bytes

		if report.err != nil {
			errors = append(errors, report.err)
		}

		reportProgress(reportMsg, progress, total, progress_report.UnitsNone, nil)
	}

	if len(errors) > 0 {
		reportProgress(reportMsg, progress, total, progress_report.UnitsNone, errors)
		return
	}

	reportProgress(reportMsg, total, total, progress_report.UnitsNone, progress_report.ErrorProgressDone)
}
