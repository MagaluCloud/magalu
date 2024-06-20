package objects

import (
	"context"

	"magalu.cloud/core"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type listParams struct {
	common.ListObjectsParams `json:",squash" mgc:"positional"` // nolint
	common.Filters           `json:",squash"`                  // nolint
}

type listResponse struct {
	Contents       []*common.BucketContent `xml:"Contents"`
	CommonPrefixes []*common.Prefix        `xml:"CommonPrefixes"`
}

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all objects from a bucket",
		},
		List,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "table=transpose,DIRS:$.CommonPrefixes,FILES:$.Contents::`KEY:$[*].Key,SIZE:fileSize($[*].ContentSize),LASTMODIFIED:humanTime($[*].LastModified)`"
	})
})

func List(ctx context.Context, params listParams, cfg common.Config) (result listResponse, err error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	objects := common.ListGenerator(ctx, params.ListObjectsParams, cfg, nil)
	objects = common.ApplyFilters(ctx, objects, params.FilterParams, cancel)
	entries, err := pipeline.SliceItemLimitedConsumer[[]pipeline.WalkDirEntry](ctx, params.MaxItems, objects)
	if err != nil {
		return result, err
	}

	contents := make([]*common.BucketContent, 0, len(entries))
	commonPrefixes := make([]*common.Prefix, 0)
	for _, entry := range entries {
		if entry.Err() != nil {
			return result, entry.Err()
		}
		if entry.DirEntry().IsDir() {
			commonPrefixes = append(commonPrefixes, entry.DirEntry().(*common.Prefix))
		} else {
			contents = append(contents, entry.DirEntry().(*common.BucketContent))
		}
	}

	result = listResponse{
		Contents:       contents,
		CommonPrefixes: commonPrefixes,
	}
	return result, nil
}
