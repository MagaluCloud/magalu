package acl

import (
	"context"

	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"

	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type getObjectACLParams struct {
	Destination schema.URI `json:"dst" jsonschema:"description=The full object URL to get the ACL information from,example:my-bucket/file.txt" mgc:"positional"`
}

var getGet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:         "get",
			Description:  "Get ACL information for the specified object",
			Observations: "object = \"id:\" (require tenant ID) - Example:id=\"a4900b57-7dbb-4906-b7e8-efed938e325c\"",
		},
		getACL,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getACL(ctx context.Context, p getObjectACLParams, cfg common.Config) (result common.AccessControlPolicy, err error) {
	req, err := newGetObjectAclRequest(ctx, p, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}

	return common.UnwrapResponse[common.AccessControlPolicy](resp, req)
}

func newGetObjectAclRequest(ctx context.Context, p getObjectACLParams, cfg common.Config) (*http.Request, error) {
	bucketName := common.NewBucketNameFromURI(p.Destination)
	url, err := common.BuildBucketHostWithPathURL(cfg, bucketName, p.Destination.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	q := url.Query()
	q.Add("acl", "")

	url.RawQuery = q.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
