package block_storage

import (
	"magalu.cloud/core"
	"magalu.cloud/sdk/static/block_storage/attachment"
)

func NewGroup() core.Grouper {
	return core.NewStaticGroup(
		"block-storage",
		"",
		"Operations for Block Storage API",
		[]core.Descriptor{
			attachment.NewGroup(),
		},
	)
}
