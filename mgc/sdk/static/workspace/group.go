package workspace

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name: "workspace",
			Description: `Workspace hold auth and runtime configuration, like tokens and log filter settings.
Users can create as many workspaces as they choose to. Auth and config operations will affect only the
current workspace, so users can alter and switch between workspaces without loosing the previous configuration`,
			Summary: "Manage workspaces for isolated auth and config settings",
			GroupID: "settings",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getGet(),
				getCreate(),
				getSet(),
				getList(),
				getDelete(),
			}
		},
	)
})
