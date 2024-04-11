package acl

import (
	"context"

	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/schema"
	"magalu.cloud/core/utils"

	"magalu.cloud/sdk/static/object_storage/common"
)

type getObjectACLParams struct {
	Destination schema.URI `json:"dst" jsonschema:"description=The full object URL to get the ACL information from,example:my-bucket/file.txt" mgc:"positional"`
}

var getGet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get ACL information for the specified object",
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

	resp, err := common.SendRequest(ctx, req)
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
