package object_storage

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/api_key"
	"magalu.cloud/sdk/static/object_storage/buckets"
	"magalu.cloud/sdk/static/object_storage/objects"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "object-storage",
			Summary:     "Operations for Object Storage API",
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
