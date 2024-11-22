package sharing

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "share",
			Summary:     "Shared Access (Beta)",
			Description: `You may share access to your account or organization with other people using Cloud Access Share`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getList(),
				getCreate(),
				getRemove(),
			}
		},
	)
})
