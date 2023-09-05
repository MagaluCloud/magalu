package static

import (
	"magalu.cloud/core"
	"magalu.cloud/sdk/static/auth"
	"magalu.cloud/sdk/static/config"
	"magalu.cloud/sdk/static/object_storage"
)

func NewGroup() *core.StaticGroup {
	return core.NewStaticGroup(
		"Static Groups Root",
		"",
		"",
		[]core.Descriptor{
			auth.NewGroup(),           // cmd: "auth"
			config.NewGroup(),         // cmd: "config"
			object_storage.NewGroup(), // cmd: "object-storage"
		},
	)
}
