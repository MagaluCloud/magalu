package static

import (
	"magalu.cloud/core"
	"magalu.cloud/sdk/static/auth"
	"magalu.cloud/sdk/static/bucket"
	"magalu.cloud/sdk/static/config"
)

func NewGroup() *core.StaticGroup {
	return core.NewStaticGroup(
		"Static Groups Root",
		"",
		"",
		[]core.Descriptor{
			auth.NewGroup(),   // cmd: "auth"
			config.NewGroup(), // cmd: "config"
			bucket.NewGroup(), // cmd: "bucket"
		},
	)
}
