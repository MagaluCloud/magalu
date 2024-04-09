package acl

import (
	"context"
	"fmt"

	"net/http"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type setObjectACLParams struct {
	Destination            mgcSchemaPkg.URI `json:"dst" jsonschema:"description=The full object URL to set the ACL information,example:my-bucket/file.txt" mgc:"positional"`
	Version                string           `json:"obj_version,omitempty" jsonschema:"description=Version of the object to set the ACL" mgc:"hidden"`
	BucketOwnerRead        bool             `json:"bucket_owner_read,omitempty" mgc:"hidden"`
	BucketOwnerFullControl bool             `json:"bucket_owner_full_control,omitempty" mgc:"hidden"`
	common.ACLPermissions  `json:",squash"` // nolint
}

var getSet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Set ACL information for the specified object",
		},
		set,
	)
	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully set ACL for object %q", result.Source().Parameters["dst"])
	})
	return exec
})

func set(ctx context.Context, p setObjectACLParams, cfg common.Config) (result core.Value, err error) {
	err = p.ACLPermissions.Validate()
	if err != nil {
		return
	}

	req, err := newSetObjectAclRequest(ctx, p, cfg)
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

func newSetObjectAclRequest(ctx context.Context, p setObjectACLParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostWithPathURL(cfg, common.NewBucketNameFromURI(p.Destination), p.Destination.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("acl", "")
	if p.Version != "" {
		query.Add("versionId", p.Version)
	}

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
