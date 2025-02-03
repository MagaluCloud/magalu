package common

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"math"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	"github.com/MagaluCloud/magalu/mgc/core/progress_report"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

var deleteObjectsLogger *zap.SugaredLogger

func deleteLogger() *zap.SugaredLogger {
	if deleteObjectsLogger == nil {
		deleteObjectsLogger = logger().Named("delete")
	}
	return deleteObjectsLogger
}

type DeleteObjectParams struct {
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Path of the object to be deleted,example=bucket1/file.txt" mgc:"positional"`
	Version     string           `json:"objVersion,omitempty" jsonschema:"description=Version of the object to be deleted"`
}

type DeleteBucketParams struct {
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Path of the bucket to be deleted,example=bucket1" mgc:"positional"`
}

type DeleteObjectsParams struct {
	Destination mgcSchemaPkg.URI
	ToDelete    <-chan pipeline.WalkDirEntry
	BatchSize   int
}

type DeleteAllObjectsInBucketParams struct {
	BucketName BucketName       `json:"bucket" jsonschema:"description=Name of the bucket to delete objects from" mgc:"positional"`
	BatchSize  int              `json:"batch_size,omitempty" jsonschema:"description=Limit of items per batch to delete,default=1000,minimum=1,maximum=1000" example:"1000"`
	Filters    `json:",squash"` // nolint
}

func newDeleteRequest(ctx context.Context, cfg Config, dst mgcSchemaPkg.URI) (*http.Request, error) {
	host, err := BuildBucketHostWithPath(cfg, NewBucketNameFromURI(dst), dst.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}
	return http.NewRequestWithContext(ctx, http.MethodDelete, string(host), nil)
}

type objectIdentifier struct {
	Key       string `xml:"Key"`
	VersionId string `xml:"VersionId,omitempty"`
}

type deleteBatchRequestBody struct {
	XMLName struct{}           `xml:"Delete"`
	Objects []objectIdentifier `xml:"Object"`
}

func DeleteSingle(ctx context.Context, params DeleteObjectParams, cfg Config) error {
	objectKey := params.Destination.AsFilePath().String()
	versionID := params.Version

	req, err := newDeleteSingleRequest(ctx, cfg, NewBucketNameFromURI(params.Destination), objectKey, versionID)
	if err != nil {
		return err
	}

	resp, err := SendRequest(ctx, req)
	if err != nil {
		return err
	}

	err = ExtractErr(resp, req)
	if err != nil {
		return err
	}

	return nil
}

func newDeleteSingleRequest(ctx context.Context, cfg Config, bucketName BucketName, objectKey string, versionID string) (*http.Request, error) {
	host, err := BuildBucketHost(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	url := fmt.Sprintf("%s/%s", host, objectKey)

	if versionID != "" {
		url = fmt.Sprintf("%s?versionId=%s", url, versionID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func newDeleteBatchRequest(ctx context.Context, cfg Config, bucketName BucketName, objKeys []objectIdentifier) (*http.Request, error) {
	host, err := BuildBucketHost(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}
	body := deleteBatchRequestBody{
		Objects: objKeys,
	}
	marshalledBody, err := xml.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, string(host), bytes.NewBuffer(marshalledBody))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("delete", "")
	req.URL.RawQuery = query.Encode()

	return req, nil
}

// Deleting an object does not yield result except there is an error. So this processor will *Skip*
// success results and *Output* errors
func createObjectDeletionProcessor(cfg Config, bucketName BucketName, progressReporter *progress_report.UnitsReporter) pipeline.Processor[[]pipeline.WalkDirEntry, error] {
	return func(ctx context.Context, dirEntries []pipeline.WalkDirEntry) (error, pipeline.ProcessStatus) {
		progressReporter.Report(0, uint64(len(dirEntries)), nil)

		var objIdentifiers []objectIdentifier
		var err error

		for _, dirEntry := range dirEntries {
			if err = dirEntry.Err(); err != nil {
				progressReporter.Report(0, 0, err)
				return &ObjectError{Err: err}, pipeline.ProcessAbort
			}

			obj, ok := dirEntry.DirEntry().(*BucketContent)
			if !ok {
				err = fmt.Errorf("expected object, got directory")
				progressReporter.Report(0, 0, err)
				return &ObjectError{Err: err}, pipeline.ProcessAbort
			}

			objIdentifiers = append(objIdentifiers, objectIdentifier{Key: obj.Key})
		}

		defer func() { progressReporter.Report(uint64(len(dirEntries)), 0, err) }()

		req, err := newDeleteBatchRequest(ctx, cfg, bucketName, objIdentifiers)
		if err != nil {
			return &ObjectError{Err: err}, pipeline.ProcessAbort
		}

		resp, err := SendRequest(ctx, req)
		if err != nil {
			return &ObjectError{Url: mgcSchemaPkg.URI(bucketName), Err: err}, pipeline.ProcessOutput
		}

		err = ExtractErr(resp, req)
		if err != nil {
			return &ObjectError{Err: err}, pipeline.ProcessAbort
		}

		deleteLogger().Infow("Deleted objects", "uri", URIPrefix+bucketName)
		return nil, pipeline.ProcessOutput
	}
}

func DeleteAllObjectsInBucket(ctx context.Context, params DeleteAllObjectsInBucketParams, cfg Config) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	dst := params.BucketName.AsURI()
	listParams := ListObjectsParams{
		Destination: dst,
		Recursive:   true,
		PaginationParams: PaginationParams{
			MaxItems: math.MaxInt64,
		},
	}

	progressReportMsg := fmt.Sprintf("Deleting objects from %q", params.BucketName)
	progressReporter := progress_report.NewUnitsReporter(ctx, progressReportMsg, 0)
	progressReporter.Start()
	defer progressReporter.End()

	onNewPage := func(objCount uint64) {
		progressReporter.Report(0, objCount, nil)
	}

	objs := ListGenerator(ctx, listParams, cfg, onNewPage)
	objs = ApplyFilters(ctx, objs, params.FilterParams, cancel)

	if params.BatchSize < MinBatchSize || params.BatchSize > MaxBatchSize {
		return core.UsageError{Err: fmt.Errorf("invalid item limit per request BatchSize, must not be lower than %d and must not be higher than %d: %d", MinBatchSize, MaxBatchSize, params.BatchSize)}
	}

	objsBatch := pipeline.Batch(ctx, objs, params.BatchSize)
	deleteObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, objsBatch, createObjectDeletionProcessor(cfg, params.BucketName, progressReporter), nil)
	deleteObjectsErrorChan = pipeline.Filter(ctx, deleteObjectsErrorChan, pipeline.FilterNonNil[error]{})

	objErr, err := pipeline.SliceItemConsumer[utils.MultiError](ctx, deleteObjectsErrorChan)
	if err != nil {
		return err
	}
	if len(objErr) > 0 {
		return objErr
	}

	return nil
}

func DeleteBucket(ctx context.Context, params DeleteBucketParams, cfg Config) error {
	req, err := newDeleteRequest(ctx, cfg, params.Destination)
	if err != nil {
		return err
	}

	resp, err := SendRequest(ctx, req)
	if err != nil {
		return err
	}

	// GA - TEMP
	if resp.StatusCode == 409 {
		return fmt.Errorf("the bucket may not be empty or may be locked.\nPlease clear up before attempting deletion.\n")
	}

	return ExtractErr(resp, req)
}

func Delete(ctx context.Context, params DeleteObjectParams, cfg Config) error {
	objKeys := []objectIdentifier{{Key: params.Destination.AsFilePath().String(), VersionId: params.Version}}

	if len(objKeys) > 1 {
		req, err := newDeleteBatchRequest(ctx, cfg, NewBucketNameFromURI(params.Destination), objKeys)
		if err != nil {
			return err
		}

		resp, err := SendRequest(ctx, req)
		if err != nil {
			return err
		}

		err = ExtractErr(resp, req)
		if err != nil {
			return err
		}
	} else {
		err := DeleteSingle(ctx, params, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteObjects(ctx context.Context, params DeleteObjectsParams, cfg Config) error {
	progressReportMsg := fmt.Sprintf("Deleting objects from %q", params.Destination.String())
	progressReporter := progress_report.NewUnitsReporter(ctx, progressReportMsg, 0)
	progressReporter.Start()
	defer progressReporter.End()

	bucketName := NewBucketNameFromURI(params.Destination)

	if params.BatchSize < MinBatchSize || params.BatchSize > MaxBatchSize {
		return core.UsageError{Err: fmt.Errorf("invalid item limit per request BatchSize, must not be lower than %d and must not be higher than %d: %d", MinBatchSize, MaxBatchSize, params.BatchSize)}
	}

	objsBatch := pipeline.Batch(ctx, params.ToDelete, params.BatchSize)
	deleteObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, objsBatch, createObjectDeletionProcessor(cfg, bucketName, progressReporter), nil)
	deleteObjectsErrorChan = pipeline.Filter(ctx, deleteObjectsErrorChan, pipeline.FilterNonNil[error]{})

	objErr, _ := pipeline.SliceItemConsumer[utils.MultiError](ctx, deleteObjectsErrorChan)
	if len(objErr) > 0 {
		return objErr
	}
	return nil
}
