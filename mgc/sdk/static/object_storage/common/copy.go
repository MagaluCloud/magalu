package common

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	"github.com/MagaluCloud/magalu/mgc/core/progress_report"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

var copyAllObjectsLogger *zap.SugaredLogger

func copyAllLogger() *zap.SugaredLogger {
	if copyAllObjectsLogger == nil {
		copyAllObjectsLogger = logger().Named("copy")
	}
	return copyAllObjectsLogger
}

type CopyObjectParams struct {
	Source       mgcSchemaPkg.URI `json:"src" jsonschema:"description=Path of the object in a bucket to be copied,example=bucket1/file.txt" mgc:"positional"`
	Destination  mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Full destination path in the bucket with desired filename,example=bucket2/dir/file.txt" mgc:"positional"`
	Version      string           `json:"obj_version,omitempty" jsonschema:"description=Version of the object to be copied"`
	StorageClass string           `json:"storage_class,omitempty" jsonschema:"description=Copy objects to other storage classes,example=cold,enum=,enum=standard,enum=cold,enum=glacier_ir,enum=cold_instant,default="`
}

type CopyAllObjectsParams struct {
	Source       mgcSchemaPkg.URI `json:"src" jsonschema:"description=Path of objects in a bucket to be copied,example=bucket1" mgc:"positional"`
	Destination  mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Full destination path in the bucket,example=bucket2/dir/" mgc:"positional"`
	StorageClass string           `json:"storage_class,omitempty" jsonschema:"description=Copy objects to other storage classes,example=cold,enum=,enum=standard,enum=cold,enum=glacier_ir,enum=cold_instant,default="`
	Filters      `json:",squash"` // nolint
}

type copier interface {
	Copy(context.Context) error
}

func newCopyRequest(ctx context.Context, cfg Config, src mgcSchemaPkg.URI, dst mgcSchemaPkg.URI, version string) (*http.Request, error) {
	host, err := BuildBucketHostWithPath(cfg, NewBucketNameFromURI(dst), dst.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, string(host), nil)
	if err != nil {
		return nil, err
	}

	copySource, err := url.JoinPath(src.Hostname(), src.Path())
	if err != nil {
		return nil, core.UsageError{Err: fmt.Errorf("badly specified source URI: %w", err)}
	}

	req.Header.Set("x-amz-copy-source", copySource)

	if version != "" {
		query := req.URL.Query()
		query.Set("versionId", version)
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

func createObjectCopyProcessor(cfg Config, params CopyAllObjectsParams, progressReporter *progress_report.UnitsReporter) pipeline.Processor[pipeline.WalkDirEntry, error] {
	return func(ctx context.Context, dirEntry pipeline.WalkDirEntry) (error, pipeline.ProcessStatus) {
		bucketName := NewBucketNameFromURI(params.Source)
		rootURI := bucketName.AsURI()
		var err error

		defer func() { progressReporter.Report(1, 0, err) }()
		path := dirEntry.Path()
		objURI := rootURI.JoinPath(path)

		if dirEntry.Err() != nil {
			err = &ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: dirEntry.Err()}
			return err, pipeline.ProcessOutput
		}

		_, ok := dirEntry.DirEntry().(*BucketContent)
		if !ok {
			err = &ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: fmt.Errorf("expected object, got directory")}
			return err, pipeline.ProcessOutput
		}

		copyAllLogger().Infow("Copying object", "uri", objURI)
		err = CopySingleFile(ctx, cfg, objURI, params.Destination.JoinPath(dirEntry.Path()), params.StorageClass)
		if err != nil {
			return err, pipeline.ProcessAbort
		}

		if err != nil {
			err = &ObjectError{Url: mgcSchemaPkg.URI(objURI), Err: err}
			return err, pipeline.ProcessOutput
		}

		return nil, pipeline.ProcessOutput
	}
}

func CopyMultipleFiles(ctx context.Context, cfg Config, params CopyAllObjectsParams) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	listParams := ListObjectsParams{
		Destination: params.Source,
		Recursive:   true,
		PaginationParams: PaginationParams{
			MaxItems: math.MaxInt64,
		},
	}
	progressReportMsg := fmt.Sprintf("Copying objects from %q to %q", params.Source, params.Destination)
	progressReporter := progress_report.NewUnitsReporter(ctx, progressReportMsg, 0)
	progressReporter.Start()
	defer progressReporter.End()

	onNewPage := func(objCount uint64) {
		progressReporter.Report(0, objCount, nil)
	}

	objs := ListGenerator(ctx, listParams, cfg, onNewPage)
	objs = ApplyFilters(ctx, objs, params.FilterParams, cancel)

	copyObjectsErrorChan := pipeline.ParallelProcess(ctx, cfg.Workers, objs, createObjectCopyProcessor(cfg, params, progressReporter), nil)
	copyObjectsErrorChan = pipeline.Filter(ctx, copyObjectsErrorChan, pipeline.FilterNonNil[error]{})

	objErr, err := pipeline.SliceItemConsumer[utils.MultiError](ctx, copyObjectsErrorChan)
	if err != nil {
		progressReporter.Report(0, 0, err)
		return err
	}
	if len(objErr) > 0 {
		progressReporter.Report(0, 0, objErr)
		return objErr
	}

	return nil
}

func CopySingleFile(ctx context.Context, cfg Config, src mgcSchemaPkg.URI, dst mgcSchemaPkg.URI, storageClass string) error {
	if dst.IsRoot() {
		dst = dst.JoinPath(src.Filename())
	}

	req, err := newCopyRequest(ctx, cfg, src, dst, "")
	if err != nil {
		return err
	}

	if storageClass != "" {
		req.Header.Set("X-Amz-Storage-Class", storageClass)
	}

	resp, err := SendRequest(ctx, req, cfg)
	if err != nil {
		return err
	}

	return ExtractErr(resp, req)
}

func NewCopier(ctx context.Context, cfg Config, src mgcSchemaPkg.URI, dst mgcSchemaPkg.URI, version string, storageClass string) (copier, error) {
	metadata, err := HeadFile(ctx, cfg, src, version)
	if err != nil {
		return nil, err
	}

	totalCopyParts := int(math.Ceil(float64(metadata.ContentLength) / float64(cfg.chunkSizeInBytes())))

	if totalCopyParts > 1 {
		return &bigFileCopier{
			cfg:          cfg,
			src:          src,
			dst:          dst,
			fileSize:     metadata.ContentLength,
			totalParts:   totalCopyParts,
			storageClass: storageClass,
		}, nil
	} else {
		return &smallFileCopier{
			cfg:          cfg,
			src:          src,
			dst:          dst,
			storageClass: storageClass,
		}, nil
	}
}
