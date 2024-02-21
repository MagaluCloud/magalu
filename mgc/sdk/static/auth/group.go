package auth

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/auth/scopes"
	"magalu.cloud/sdk/static/auth/tenant"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "auth",
			Summary: "Actions with ID Magalu to log in, refresh tokens, change tenants and others",
			Description: `The authentication credentials set here will be used as a basis for a variety
of HTTP requests using the MgcSDK. Authentication is done via Magalu Cloud account
(Object Storage requires special keys, refer to it for more info)`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getLogin(),
				getAccessToken(),
				scopes.GetGroup(),
				tenant.GetGroup(),
			}
		},
	)
})
