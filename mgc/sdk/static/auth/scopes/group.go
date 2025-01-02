package scopes

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name: "scopes",
			Description: `Some operations require scopes to be executed. These scopes
can be managed here, with operations that change the current
access token used in all other operations.`,
			Summary: "Manage scope operations for current access token",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getAdd(),
				getAddAll(),
				getListAll(),
				getListCurrent(),
				getRemove(),
				getSet(),
			}
		},
	)
})
