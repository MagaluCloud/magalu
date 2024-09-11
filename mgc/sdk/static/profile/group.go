package profile

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "profile",
			Summary: "Manage account settings, including SSH keys and related configurations",
			Description: `The profile group provides commands to view and modify user account settings. 
It allows users to manage their SSH keys, update personal information, and configure other 
account-related preferences. This group is essential for maintaining secure access and 
personalizing the user experience within the system.`,
			GroupID: "settings",
		},
		func() []core.Descriptor {
			return []core.Descriptor{}
		},
	)
})
