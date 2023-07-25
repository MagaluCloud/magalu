package static

import (
	"magalu.cloud/core"
	"magalu.cloud/sdk/static/auth"
	"magalu.cloud/sdk/static/config"
	"magalu.cloud/sdk/static/vpc"
)

func NewGroup() *core.StaticGroup {
	return core.NewStaticGroup(
		"Static Groups Root",
		"",
		"",
		[]core.Descriptor{
			newStatic(),       // cmd: "static"
			auth.NewGroup(),   // cmd: "auth"
			vpc.NewGroup(),    // cmd: "vpc"
			config.NewGroup(), // cmd: "config"
		},
	)
}
