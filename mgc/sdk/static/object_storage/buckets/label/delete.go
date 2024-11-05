package label

import (
	"context"
	"fmt"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type deleteBucketLabelParams struct {
	Bucket common.BucketName `json:"bucket" jsonschema:"description=Name of the bucket to delete labels from,example=my-bucket" mgc:"positional"`
	Labels string            `json:"label" jsonschema:"description=Label values, comma separated, without whitespaces" mgc:"positional"`
}

var getDelete = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete",
			Description: "Delete labels for the specified bucket",
		},
		deleteLabels,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully deleted labels for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func deleteLabels(ctx context.Context, params deleteBucketLabelParams, cfg common.Config) (_ core.Value, err error) {
	res, err := getTags(ctx, GetBucketLabelParams{Bucket: params.Bucket}, cfg)
	if err != nil {
		return
	}

	labels, hasMGCLabels := findMGCLabels(res.Tags)
	if !hasMGCLabels || labels == "" {
		return
	}

	savedLabels := strings.Split(labels, ",")
	inputLabels := strings.Split(params.Labels, ",")
	tagLabel := removeMatchingStrings(savedLabels, inputLabels)
	updateMGCLabels(&res.Tags, tagLabel)

	taggingParams := setBucketTaggingParams{Bucket: params.Bucket, TagSet: res}
	req, err := newSetBucketTaggingRequest(ctx, taggingParams, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	err = common.ExtractErr(resp, req)
	if err != nil {
		return
	}

	return
}

func removeMatchingStrings(tags1 []string, tags2 []string) string {
	tags2Map := make(map[string]Unit)
	for _, tag := range tags2 {
		tags2Map[tag] = struct{}{}
	}

	updatedTags := make([]string, 0, len(tags1))
	for _, tag := range tags1 {
		if _, found := tags2Map[tag]; !found {
			updatedTags = append(updatedTags, tag)
		}
	}

	return strings.Join(updatedTags, ",")
}
