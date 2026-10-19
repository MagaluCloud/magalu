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
			Summary:     preview.Alpha.Summary,
			Description: "Early preview commands, for tenants invited. 'mgc preview list' shows the ones your tenant was given, and 'mgc preview install' adds them",
			IsInternal:  utils.BoolPtr(true),
		},
		func() []core.Descriptor { return children },
	)

	return core.NewStaticGroup(
		core.DescriptorSpec{Name: Name + " root", Description: "Alpha commands that ship with the CLI"},
		func() []core.Descriptor { return []core.Descriptor{group} },
	)
}
