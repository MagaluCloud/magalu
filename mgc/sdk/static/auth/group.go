package auth

import (
	"magalu.cloud/core"
	"magalu.cloud/sdk/static/auth/tenant"
)

func NewGroup() *core.StaticGroup {
	return core.NewStaticGroup(
		"auth",
		"",
		"Actions with ID Magalu to login, refresh tokens, change tenants and others",
		[]core.Descriptor{
			newLogin(),
			newAccessToken(),
			tenant.NewGroup(),
		},
	)
}
