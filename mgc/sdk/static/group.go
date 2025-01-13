package static

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/auth"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/config"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/http"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/profile"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/workspace"
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
