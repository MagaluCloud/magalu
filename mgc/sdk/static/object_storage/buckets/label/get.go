package label

import (
	"context"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type Label struct {
	Value string `json:"Labels"`
}

type Tag struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

type TagSet struct {
	Tags []Tag `xml:"TagSet>Tag" json:"Labels"`
}

type GetBucketLabelParams struct {
	Bucket common.BucketName `json:"bucket" jsonschema:"description=Specifies the bucket whose labels is being requested" mgc:"positional"`
}

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get labels for the specified bucket",
		},
		getLabels,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "xml"
	})
	return exec
})

func getLabels(ctx context.Context, params GetBucketLabelParams, cfg common.Config) (_ Label, err error) {
	res, err := getTags(ctx, GetBucketLabelParams{Bucket: params.Bucket}, cfg)
	if err != nil {
		return
	}
	labels, _ := findMGCLabels(res.Tags)
	return Label{Value: labels}, err
}

func getTags(ctx context.Context, params GetBucketLabelParams, cfg common.Config) (_ TagSet, err error) {
	req, err := newGetTaggingRequest(ctx, cfg, params.Bucket)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}

	return common.UnwrapResponse[TagSet](res, req)
}

func newGetTaggingRequest(ctx context.Context, cfg common.Config, bucketName common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("tagging", "")
	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
