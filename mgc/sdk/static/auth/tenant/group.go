package tenant

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "tenant",
			Summary: "Manage Tenants",
			Description: `Tenants work like sub-accounts. You may have more than one Tenant under your
Magalu Cloud account and they each store their data separately, but are billed
under the same account`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getList(),
				getSet(),
				getCurrent(),
			}
		},
	)
})
