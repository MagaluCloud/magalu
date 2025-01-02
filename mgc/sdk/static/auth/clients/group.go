package clients

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
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
