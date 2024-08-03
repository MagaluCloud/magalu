package profile

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name: "profile",
			Description: `Profiles hold auth and runtime configuration for the MgcSDK, like tokens and log filter settings.
Users can create as many profiles as they choose to. Auth and config operations will affect only the
current profile, so users can alter and switch between profiles without loosing the previous configuration`,
			Summary: "Profile related commands",
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
