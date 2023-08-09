package auth

import (
	"magalu.cloud/core"
	"magalu.cloud/sdk/static/auth/tenant"
)

func NewGroup() *core.StaticGroup {
	return core.NewStaticGroup(
		"auth",
		"",
		"",
		[]core.Descriptor{
			newLogin(),       // cmd: auth login
			newAccessToken(), // cmd: auth access_token
			tenant.NewTenant(),
		},
	)
}
