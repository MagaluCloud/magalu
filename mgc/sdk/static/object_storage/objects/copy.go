package objects

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"magalu.cloud/core"
	"magalu.cloud/core/progress_report"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type copyObjectParams struct {
	Source      mgcSchemaPkg.URI `json:"src" jsonschema:"description=Path of the object to be copied,example=s3://bucket1/file.txt" mgc:"positional"`
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Full destination path in the bucket with desired filename,example=s3://bucket2/dir/file.txt" mgc:"positional"`
}

var getCopy = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "copy",
			Description: "Copy an object from a bucket to another",
		},
		copy,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Copied from {{.src}} to {{.dst}}\n"
	})
})

func newCopyRequest(ctx context.Context, cfg common.Config, dst mgcSchemaPkg.URI) (*http.Request, error) {
	host, err := common.BuildBucketHostWithPath(cfg, common.NewBucketNameFromURI(dst), dst.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}
	return http.NewRequestWithContext(ctx, http.MethodPut, string(host), nil)
}

func copySingleFile(ctx context.Context, cfg common.Config, src mgcSchemaPkg.URI, dst mgcSchemaPkg.URI) error {
	req, err := newCopyRequest(ctx, cfg, dst)
	if err != nil {
		return err
	}

	copySource, err := url.JoinPath(src.Hostname(), src.Path())
	if err != nil {
		return core.UsageError{Err: fmt.Errorf("badly specified source URI: %w", err)}
	}

	req.Header.Set("x-amz-copy-source", copySource)

	_, err = common.SendRequest(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func copy(ctx context.Context, p copyObjectParams, cfg common.Config) (result core.Value, err error) {
	reportProgress := progress_report.FromContext(ctx)
	reportMsg := "Copying object from " + p.Source.String() + " to " + p.Destination.String()
	progress := uint64(0)
	total := uint64(1)

	reportProgress(reportMsg, progress, total, progress_report.UnitsNone, nil)

	err = copySingleFile(ctx, cfg, p.Source, p.Destination)
	if err != nil {
		reportProgress(reportMsg, progress, total, progress_report.UnitsNone, err)
		return nil, err
	}

	reportProgress(reportMsg, total, total, progress_report.UnitsNone, progress_report.ErrorProgressDone)

	return copyObjectParams{Source: p.Source, Destination: p.Destination}, err
}
