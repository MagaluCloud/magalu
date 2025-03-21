package buckets

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type getParams struct {
	BucketName common.BucketName `json:"bucket" jsonschema:"description=Name of the bucket to retrieve" mgc:"positional"`
}

type bucketInfo struct {
	BucketName    string   `json:"bucket_name"`
	Versioning    string   `json:"versioning,omitempty"`
	Policy        string   `json:"policy,omitempty"`
	ACL           string   `json:"acl,omitempty"`
	Visibility    string   `json:"visibility,omitempty"`
	CreationDate  string   `json:"creation_date,omitempty"`
	ObjectLocking string   `json:"object_locking,omitempty"`
	OwnerID       string   `json:"owner_id,omitempty"`
	OwnerName     string   `json:"owner_name,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
}

type AccessControlPolicy struct {
	XMLName           xml.Name `xml:"AccessControlPolicy"`
	Owner             Owner    `xml:"Owner"`
	AccessControlList ACL      `xml:"AccessControlList"`
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type ACL struct {
	Grants []Grant `xml:"Grant"`
}

type Grant struct {
	Grantee    Owner  `xml:"Grantee"`
	Permission string `xml:"Permission"`
}

type VersioningConfiguration struct {
	XMLName xml.Name `xml:"VersioningConfiguration"`
	Status  string   `xml:"Status"`
}

type ObjectLockConfiguration struct {
	XMLName          xml.Name `xml:"ObjectLockConfiguration"`
	ObjectLockStatus string   `xml:"ObjectLockEnabled"`
	Rule             struct {
		DefaultRetention struct {
			Mode string `xml:"Mode"`
			Days int    `xml:"Days"`
		} `xml:"DefaultRetention"`
	} `xml:"Rule"`
}

var getBucket = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewReflectedSimpleExecutor[getParams, common.Config, *bucketInfo](
		core.ExecutorSpec{
			DescriptorSpec: core.DescriptorSpec{
				Name:        "get",
				Description: "Retrieve detailed information about a bucket",
				IsInternal:  utils.BoolPtr(true),
			},
		},
		getValidBucket,
	)
})

func getValidBucket(ctx context.Context, params getParams, cfg common.Config) (*bucketInfo, error) {
	info := &bucketInfo{
		BucketName: params.BucketName.String(),
	}

	aclReq, err := newGetRequest(ctx, params.BucketName, cfg, "acl")
	if err != nil {
		return nil, err
	}

	aclRes, err := common.SendRequest(ctx, aclReq)
	if err == nil {
		body, _ := io.ReadAll(aclRes.Body)
		defer aclRes.Body.Close()

		if strings.Contains(string(body), "<AccessControlPolicy") {
			var acp AccessControlPolicy
			if err := xml.Unmarshal(body, &acp); err == nil {
				info.OwnerID = acp.Owner.ID
				info.OwnerName = acp.Owner.DisplayName

				for _, grant := range acp.AccessControlList.Grants {
					info.Permissions = append(info.Permissions, grant.Permission)
				}
			}
		}
	}

	versioningReq, err := newGetRequest(ctx, params.BucketName, cfg, "versioning")
	if err != nil {
		return nil, err
	}

	versioningRes, err := common.SendRequest(ctx, versioningReq)
	if err == nil {
		body, _ := io.ReadAll(versioningRes.Body)
		defer versioningRes.Body.Close()

		if strings.Contains(string(body), "<VersioningConfiguration") {
			var vc VersioningConfiguration
			if err := xml.Unmarshal(body, &vc); err == nil {
				info.Versioning = vc.Status
			}
		}
	}

	objectLockReq, err := newGetRequest(ctx, params.BucketName, cfg, "object-lock")
	if err != nil {
		return nil, err
	}

	objectLockRes, err := common.SendRequest(ctx, objectLockReq)
	if err == nil {
		body, _ := io.ReadAll(objectLockRes.Body)
		defer objectLockRes.Body.Close()

		if strings.Contains(string(body), "<ObjectLockConfiguration") {
			var ol ObjectLockConfiguration
			if err := xml.Unmarshal(body, &ol); err == nil {
				info.ObjectLocking = fmt.Sprintf("%s (%d days)", ol.Rule.DefaultRetention.Mode, ol.Rule.DefaultRetention.Days)
			}
		}
	}

	return info, nil
}

func newGetRequest(ctx context.Context, bucketName common.BucketName, cfg common.Config, params ...string) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	for _, param := range params {
		query.Set(param, "")
	}

	url.RawQuery = query.Encode()
	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
