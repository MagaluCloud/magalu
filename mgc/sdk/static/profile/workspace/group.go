package workspace

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetWorkspace = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name: "workspaces",
			Description: `Workspaces hold auth and runtime configuration for the MGC CLI, like tokens and log filter settings.
Users can create as many workspace as they choose to. Auth and config operations will affect only the
current workspace, so users can alter and switch between workspace without loosing the previous configuration`,
			Summary: "Workspace related commands",
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
