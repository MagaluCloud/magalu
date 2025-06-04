package cors

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "cors",
			Description: "CORS-related commands",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(),    // object-storage buckets cors get
				getSet(),    // object-storage buckets cors set
				getDelete(), // object-storage buckets cors delete
			}
		},
	)
})
