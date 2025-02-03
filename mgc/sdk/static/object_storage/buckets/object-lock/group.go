package object_lock

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "object-lock",
			Description: "Object locking commands",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(),   // object-storage buckets object-lock get
				getSet(),   // object-storage buckets object-lock set
				getUnset(), // object-storage buckets object-lock unset
			}
		},
	)
})
