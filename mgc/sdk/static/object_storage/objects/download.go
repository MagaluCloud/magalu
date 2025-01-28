package objects

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

var getDownload = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "download",
			Summary:     "Download an object from a bucket",
			Description: "Download an object from a bucket. If no destination is specified, the default is the current working directory",
		},
		download,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Downloaded from {{.src}} to {{.dst}}\n"
	})
})

func download(ctx context.Context, p common.DownloadObjectParams, cfg common.Config) (result core.Value, err error) {
	if p.Source.Path() == "" {
		return nil, core.UsageError{Err: fmt.Errorf("invalid source specified. Please include the object key in addition to the bucket name")}
	}

	dst, err := common.GetDownloadFileDst(p.Destination, p.Source)
	if err != nil {
		return nil, fmt.Errorf("no destination specified and could not use local dir: %w", err)
	}

	downloader, err := common.NewDownloader(ctx, cfg, p.Source, dst, p.Version)
	if err != nil {
		return nil, err
	}

	if err = downloader.Download(ctx); err != nil {
		return nil, err
	}

	return common.DownloadObjectParams{Source: p.Source, Destination: mgcSchemaPkg.FilePath(dst)}, nil
}
