package object_storage

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/api_key"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/buckets"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/objects"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "object-storage",
			Summary:     "Operations for Object Storage",
			Description: `Create and manage Buckets and Objects via the Object Storage API`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				buckets.GetGroup(), // object-storage buckets
				objects.GetGroup(), // object-storage objects
				api_key.GetGroup(), // object-storage api-keys
			}
		},
	)
})
