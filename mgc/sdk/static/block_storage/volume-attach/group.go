package attachment

import (
	"magalu.cloud/core"
)

func NewGroup() core.Grouper {
	return core.NewStaticGroup(
		"volume-attach",
		"",
		"Block Storage Volume Attachment",
		[]core.Descriptor{
			newCreate(),
			newGet(),
			newUpdate(),
			newDelete(),
		},
	)
}
