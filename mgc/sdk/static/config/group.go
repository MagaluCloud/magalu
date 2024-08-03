package config

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "config",
			Summary: "Manage Configuration values",
			GroupID: "settings",
			Description: `Configuration values are available to be set so that they persist between
different executions of the MgcSDK. They reside in a YAML file when set.
Config values may also be loaded via Environment Variables. Any Config available
(see 'list') may be exported as an env variable in uppercase with the 'MGC_' prefix`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getList(),
				getGet(),
				getGetSchema(),
				getSet(),
				getDelete(),
			}
		},
	)
})
