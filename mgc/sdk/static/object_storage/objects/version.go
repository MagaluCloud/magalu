package objects

import (
	"context"
	"encoding/xml"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type versioningObjectParams struct {
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Path of the object to retrieve versions from,example=bucket1/file.txt" mgc:"positional"`
}

type ListObjectVersionsResponse struct {
	XMLName  xml.Name        `xml:"ListVersionsResult"`
	Text     string          `xml:",chardata"`
	Versions []ObjectVersion `xml:"Version"`
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ObjectVersion.html
type ObjectVersion struct {
	Text         string       `xml:",chardata"`
	VersionID    string       `xml:"VersionId"`
	IsLatest     bool         `xml:"IsLatest"`
	Key          string       `xml:"Key"`
	LastModified string       `xml:"LastModified"`
	ETag         string       `xml:"ETag"`
	Size         int64        `xml:"Size"`
	Owner        common.Owner `xml:"Owner"`
	StorageClass string       `xml:"StorageClass"`
}

var getVersions = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "versions",
			Description: "Retrieve all versions of an object",
		},
		getObjectVersioning,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getObjectVersioning(ctx context.Context, params versioningObjectParams, cfg common.Config) (result []ObjectVersion, err error) {
	req, err := newGetObjectVersioningRequest(ctx, cfg, params)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	resposta, err := common.UnwrapResponse[ListObjectVersionsResponse](resp, req)
	if err != nil {
		return nil, err
	}

	return resposta.Versions, nil
}

func newGetObjectVersioningRequest(ctx context.Context, cfg common.Config, params versioningObjectParams) (*http.Request, error) {
	url, err := common.BuildBucketHostWithPathURL(cfg, common.NewBucketNameFromURI(params.Destination), "")
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectVersions.html#:~:text=in%20the%20specified-,bucket,-.
	query := url.Query()
	query.Set("prefix", params.Destination.Path())
	query.Set("encoding-type", "url")
	query.Set("versions", "")

	url.RawQuery = query.Encode()
	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
