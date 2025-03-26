package buckets

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/acl"
	object_lock "github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/object-lock"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/versioning"
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
	info := &bucketInfo{BucketName: params.BucketName.String()}

	aclRes, err := acl.GetACL(ctx, acl.GetBucketACLParams{Bucket: params.BucketName}, cfg)
	if err != nil {
		return nil, err
	}

	info.OwnerID = aclRes.Owner.ID
	info.OwnerName = aclRes.Owner.DisplayName
	for _, grant := range aclRes.AccessControlList.Grant {
		info.Permissions = append(info.Permissions, grant.Permission)
	}

	versioningRes, err := versioning.GetBucketVersioning(ctx, versioning.GetBucketVersioningParams{Bucket: params.BucketName}, cfg)
	if err != nil {
		return nil, err
	}

	info.Versioning = versioningRes.Status

	objectLockRes, err := getObjectLocking(ctx, params.BucketName, cfg)
	if err != nil {
		return nil, err
	}
	if objectLockRes != nil {
		info.ObjectLocking = fmt.Sprintf("%s (%d days)", objectLockRes.Rule.DefaultRetention.Mode, objectLockRes.Rule.DefaultRetention.Days)
	} else {
		info.ObjectLocking = "Not Configured"
	}

	return info, nil
}

func getObjectLocking(ctx context.Context, bucket common.BucketName, cfg common.Config) (*object_lock.GetBucketObjectLockResponse, error) {
	objectLockRes, err := object_lock.GetObjectLocking(ctx, object_lock.GetBucketObjectLockParams{Bucket: bucket}, cfg)
	if err != nil {
		if errors.Is(err, object_lock.ErrBucketMissingObjectLockConfiguration) {
			return nil, nil
		}
		return nil, err
	}
	return &objectLockRes, nil
}
