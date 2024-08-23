package static

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/auth"
	"magalu.cloud/sdk/static/config"
	"magalu.cloud/sdk/static/http"
	"magalu.cloud/sdk/static/object_storage"
	"magalu.cloud/sdk/static/profile"
	"magalu.cloud/sdk/static/workspace"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{Name: "Static Groups Root"},
		func() []core.Descriptor {
			return []core.Descriptor{
				auth.GetGroup(),
				config.GetGroup(),
				object_storage.GetGroup(),
				workspace.GetGroup(),
				http.GetGroup(),
				profile.GetGroup(),
			}
		},
	)
})
