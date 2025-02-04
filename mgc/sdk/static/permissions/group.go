package permissions

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "permissions",
			Summary:     "Manage account (RBAC) - CHANGE-ME",
			Description: `CHANGE-ME`,
			GroupID:     "settings",
		},
		func() []core.Descriptor {
			return []core.Descriptor{}
		},
	)
})
