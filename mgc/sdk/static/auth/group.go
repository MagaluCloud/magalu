package auth

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/auth/api_key"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/auth/clients"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/auth/tenant"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "auth",
			Summary: "Actions with ID Magalu to log in, API Keys, refresh tokens, change tenants and others",
			Description: `The authentication credentials set here will be used as a basis for a variety
of HTTP requests using the MgcSDK. Authentication is done via Magalu Cloud account
(Object Storage requires special keys, refer to it for more info)`,
			GroupID: "settings",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getLogin(),
				getAccessToken(),
				getLogout(),
				tenant.GetGroup(),
				clients.GetGroup(),
				api_key.GetGroup(),
			}
		},
	)
})
