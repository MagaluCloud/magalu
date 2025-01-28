package buckets

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/acl"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/label"
	object_lock "github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/object-lock"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/policy"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets/versioning"
)

var GetGroup = utils.NewLazyLoader[core.Grouper](func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "buckets",
			Description: "Bucket operations for Object Storage API",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getCreate(),            // object-storage buckets create
				getDelete(),            // object-storage buckets delete
				getList(),              // object-storage buckets list
				getBucket(),            // object-storage buckets get
				getPublicUrl(),         // object-storage objects public-url
				acl.GetGroup(),         // object-storage buckets acl
				versioning.GetGroup(),  // object-storage buckets versioning
				policy.GetGroup(),      // object-storage buckets policy
				label.GetGroup(),       // object-storage buckets label
				object_lock.GetGroup(), // object-storage buckets object-lock
			}
		},
	)
})
