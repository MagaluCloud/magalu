package profile

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	workspaces "magalu.cloud/sdk/static/profile/workspace"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "profile",
			Summary: "Actions with profile",
			GroupID: "settings",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				workspaces.GetWorkspace(),
			}
		},
	)
})
