package buckets

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

var deleteBucketsLogger *zap.SugaredLogger

func deleteLogger() *zap.SugaredLogger {
	if deleteBucketsLogger == nil {
		deleteBucketsLogger = logger().Named("delete")
	}
	return deleteBucketsLogger
}

type deleteParams struct {
	BucketName common.BucketName `json:"bucket" jsonschema:"description=Name of the bucket to be deleted" mgc:"positional"`
}

var getDelete = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var executor core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete",
			Description: "Delete an existing Bucket",
		},
		deleteBucket,
	)

	executor = core.NewConfirmableExecutor(
		executor,
		core.ConfirmPromptWithTemplate(
			"This command will delete bucket {{.parameters.bucket}}, and its result is NOT reversible.",
		),
	)

	return core.NewExecuteFormat(executor, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Deleted bucket %q", result.Source().Parameters["bucket"])
	})
})

func deleteBucket(ctx context.Context, params deleteParams, cfg common.Config) (bool, error) {
	logger := deleteLogger().Named("delete").With(
		"params", params,
		"cfg", cfg,
	)

	err := common.DeleteAllObjectsInBucket(ctx, common.DeleteAllObjectsInBucketParams{BucketName: params.BucketName, BatchSize: common.MaxBatchSize}, cfg)
	if err != nil {
		return false, err
	}

	dst := params.BucketName.AsURI()
	err = common.DeleteBucket(ctx, common.DeleteBucketParams{Destination: dst}, cfg)
	if err != nil {
		return false, err
	}

	logger.Info("Deleted bucket")
	return true, err
}
