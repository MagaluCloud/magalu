package buckets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type BucketResponse struct {
	CreationDate string       `xml:"CreationDate"`
	Name         string       `xml:"Name"`
	Size         SizeResponse `xml:"Size"`
}

type SizeResponse struct {
	Standard    string `xml:"Standard"`
	ColdInstant string `xml:"ColdInstant"`
	Versions    string `xml:"Versions"`
	UpdatedAt   string `xml:"UpdatedAt"`
}

type ListResponse struct {
	Buckets []*BucketResponse `xml:"Buckets>Bucket"`
	Owner   *common.Owner     `xml:"Owner"`
}

type Bucket struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
	Size         Size   `xml:"Size"`
}

type Size struct {
	Standard    float64 `xml:"Standard"`
	ColdInstant float64 `xml:"ColdInstant"`
	Versions    float64 `xml:"Versions"`
	UpdatedAt   string  `xml:"UpdatedAt"`
}

type APIReturn struct {
	Buckets []*Bucket     `xml:"Buckets>Bucket"`
	Owner   *common.Owner `xml:"Owner"`
}

func newListRequest(ctx context.Context, cfg common.Config) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(common.BuildHost(cfg)), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("size", "")
	req.URL.RawQuery = q.Encode()

	return req, nil
}

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all existing Buckets",
		},
		list,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "table"
	})
	return exec
})

func list(ctx context.Context, _ struct{}, cfg common.Config) (*ListResponse, error) {
	req, err := newListRequest(ctx, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return nil, err
	}

	response, err := common.UnwrapResponse[APIReturn](resp, req)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	result.Owner = response.Owner
	for _, bucket := range response.Buckets {
		result.Buckets = append(result.Buckets, &BucketResponse{
			CreationDate: bucket.CreationDate,
			Name:         bucket.Name,
			Size: SizeResponse{
				Standard:    FormatSize(bucket.Size.Standard),
				ColdInstant: FormatSize(bucket.Size.ColdInstant),
				Versions:    FormatSize(bucket.Size.Versions),
				UpdatedAt:   bucket.Size.UpdatedAt,
			},
		})
	}

	return &result, nil
}

func FormatSize(size float64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%.1fB", size)
	}

	suffixes := []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	i := -1
	sizeFloat := float64(size)

	for sizeFloat >= unit && i < len(suffixes)-1 {
		i++
		sizeFloat /= unit
	}

	return fmt.Sprintf("%.1f %s", sizeFloat, suffixes[i])
}
