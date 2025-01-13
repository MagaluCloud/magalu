package api_key

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "api-key",
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
