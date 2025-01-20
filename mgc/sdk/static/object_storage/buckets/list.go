package buckets

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type BucketResponse struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
	BucketSize   string `xml:"BucketSize"`
}

type ListResponse struct {
	Buckets []*BucketResponse `xml:"Buckets>Bucket"`
	Owner   *common.Owner     `xml:"Owner"`
}

func newListRequest(ctx context.Context, cfg common.Config) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, string(common.BuildHost(cfg)), nil)
}

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all existing Buckets",
			// Scopes:      core.Scopes{"object-storage.read"},
		},
		list,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "table"
	})
	return exec
})

func list(ctx context.Context, _ struct{}, cfg common.Config) (result ListResponse, err error) {
	req, err := newListRequest(ctx, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	result, err = common.UnwrapResponse[ListResponse](resp, req)
	if err != nil {
		return
	}

	for _, bucket := range result.Buckets {
		size, err := strconv.ParseInt(bucket.BucketSize, 10, 64)
		if err != nil {
			bucket.BucketSize = "-"
		} else {
			bucket.BucketSize = FormatSize(size)
		}
	}

	return
}

func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}

	suffixes := []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	i := 0
	sizeFloat := float64(size)

	for sizeFloat >= unit && i < len(suffixes)-1 {
		sizeFloat /= unit
		i++
	}

	return fmt.Sprintf("%.1f %s", sizeFloat, suffixes[i])
}
