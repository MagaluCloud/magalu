package acl

import (
	"context"
	"fmt"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type setBucketACLParams struct {
	Bucket                common.BucketName `json:"dst" jsonschema:"description=Name of the bucket to set permissions for,example=my-bucket" mgc:"positional"`
	common.ACLPermissions `json:",squash"`  // nolint
}

var getSet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "set permission information for the specified bucket",
		},
		setACL,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		fmt.Println(result.Source().Context.Err().Error())
		return fmt.Sprintf("Successfully set ACL for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func setACL(ctx context.Context, params setBucketACLParams, cfg common.Config) (result core.Value, err error) {
	err = params.ACLPermissions.Validate()
	if err != nil {
		return
	}

	req, err := newSetBucketACLRequest(ctx, params, cfg)
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

func newSetBucketACLRequest(ctx context.Context, p setBucketACLParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("acl", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url.String(), nil)
	if err != nil {
		return nil, err
	}

	if p.ACLPermissions.IsEmpty() {
		return nil, core.UsageError{Err: fmt.Errorf("needs to pass either grant permissions or canned info")}
	}

	err = p.ACLPermissions.SetHeaders(req, cfg)
	if err != nil {
		return nil, err
	}

	return req, nil
}
