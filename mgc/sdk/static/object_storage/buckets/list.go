package buckets

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/buckets/label"
	"magalu.cloud/sdk/static/object_storage/common"
)

type BucketResponse struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
	BucketSize   string `json:"BucketSize,omitempty" xml:"BucketSize"`
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
		},
		list,
	)
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

	var wg sync.WaitGroup
	sizeChannel := make(chan *BucketResponse)

	for _, bucket := range result.Buckets {
		wg.Add(1)
		go func(bucket *BucketResponse) {
			defer wg.Done()
			size, err := getBucketSizeFromTag(ctx, bucket.Name, cfg)
			if err != nil {
				bucket.BucketSize = "0B"
			} else {
				bucket.BucketSize = FormatSize(size)
			}
			sizeChannel <- bucket
		}(bucket)
	}

	go func() {
		wg.Wait()
		close(sizeChannel)
	}()

	for range sizeChannel {
	}

	return
}

func getBucketSizeFromTag(ctx context.Context, bucketName string, cfg common.Config) (int64, error) {
	tagSet, err := label.GetTags(ctx, label.GetBucketLabelParams{Bucket: common.BucketName(bucketName)}, cfg)
	if err != nil {
		return 0, err
	}

	for _, tag := range tagSet.Tags {
		if tag.Key == "MGC_SIZE" {
			size, parseErr := strconv.ParseInt(tag.Value, 10, 64)
			if parseErr != nil {
				return 0, parseErr
			}
			return size, nil
		}
	}

	return 0, nil
}

func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%ciB",
		float64(size)/float64(div), "KMGTPE"[exp])
}
