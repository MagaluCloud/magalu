package objects

import (
	"context"
	"fmt"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

var getCopy = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "copy",
			Description: "Copy an object from a bucket to another bucket",
		},
		copy,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Copied from {{.src}} to {{.dst}}\n"
	})
})

func copy(ctx context.Context, p common.CopyObjectParams, cfg common.Config) (result core.Value, err error) {
	_, err = common.HeadFile(ctx, cfg, p.Source, p.Version)
	if err != nil {
		return nil, fmt.Errorf("error validating source: %w", err)
	}

	fileName := p.Source.Filename()
	if fileName == "" {
		return nil, core.UsageError{Err: fmt.Errorf("source must be a URI to an object")}
	}

	fullDstPath := p.Destination
	if fullDstPath == "" {
		return nil, core.UsageError{Err: fmt.Errorf("destination cannot be empty")}
	}

	if strings.HasSuffix(fullDstPath.String(), "/") || p.Destination.IsRoot() {
		// If it isn't a file path, don't rename, just append source with bucket URI
		fullDstPath = fullDstPath.JoinPath(fileName)
	}

	copier, err := common.NewCopier(ctx, cfg, p.Source, fullDstPath, p.Version, p.StorageClass)
	if err != nil {
		return nil, err
	}

	if err = copier.Copy(ctx); err != nil {
		return nil, err
	}

	return common.CopyObjectParams{Source: p.Source, Destination: fullDstPath}, err
}
