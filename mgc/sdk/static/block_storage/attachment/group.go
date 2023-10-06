package attachment

import (
	"magalu.cloud/core"
)

func NewGroup() core.Grouper {
	return core.NewStaticGroup(
		"attachment",
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
