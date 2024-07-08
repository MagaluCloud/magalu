package api_key

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "key-pair",
			Summary: "Manage credentials to use Object Storage",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getCreate(),
				getGet(),
				getGetCurrent(),
				getList(),
				getRevoke(),
				getSetCurrent(),
				getAdd(),
			}
		},
	)
})
