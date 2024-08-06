package policy

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "policy",
			Description: "Policy-related commands",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(),    // object-storage buckets acl get
				getSet(),    // object-storage buckets acl set
				getDelete(), // object-storage buckets acl delete
			}
		},
	)
})
