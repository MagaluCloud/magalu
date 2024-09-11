package label

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "label",
			Description: "Label-related commands",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(),    // object-storage buckets label get
				getSet(),    // object-storage buckets label set
				getDelete(), // object-storage buckets label delete
			}
		},
	)
})
