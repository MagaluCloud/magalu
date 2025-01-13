package acl

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "acl",
			Description: "ACL related operations",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(), // object-storage objects acl get
				getSet(), // object-storage objects acl set
			}
		},
	)
})
