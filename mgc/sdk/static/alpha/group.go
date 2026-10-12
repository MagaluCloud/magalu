package alpha

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
)

const Name = "alpha"

func NewRoot(children ...core.Descriptor) core.Grouper {
	group := core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        Name,
			GroupID:     preview.GroupID,
			Summary:     preview.Alpha.Summary,
			Description: "Early preview commands, for customers invited by Magalu. Run 'mgc alpha sync <invite>' to download the ones your invite gives access to",
			IsInternal:  utils.BoolPtr(true),
		},
		func() []core.Descriptor { return children },
	)

	return core.NewStaticGroup(
		core.DescriptorSpec{Name: Name + " root", Description: "Alpha commands that ship with the CLI"},
		func() []core.Descriptor { return []core.Descriptor{group} },
	)
}
