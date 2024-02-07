package keys

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "keys",
			Summary: "Manage credentials to use Object Storage",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getKeys(),
			}
		},
	)
})
