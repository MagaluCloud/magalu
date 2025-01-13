package versioning

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "versioning",
			Description: "Manage bucket versioning",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(),     // object-storage buckets versioning get
				getEnable(),  // object-storage buckets versioning enable
				getSuspend(), // object-storage buckets versioning suspend
			}
		},
	)
})
