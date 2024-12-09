package object_lock

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "object-lock",
			Description: "Object locking commands",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(), // object-storage object object-lock get
				getSet(), // object-storage object object-lock set
			}
		},
	)
})
