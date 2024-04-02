package clients

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "clients",
			Summary: "Manage Clients (Oauth Applications) to use ID Magalu",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getCreate(),
				getList(),
				getUpdate(),
			}
		},
	)
})
