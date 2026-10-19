package alpha

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/project"
)

const Name = "alpha"

var GetRoot = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{Name: Name + " root", Description: "Alpha commands that ship with the CLI"},
		func() []core.Descriptor {
			return []core.Descriptor{getGroup()}
		},
	)
})

var getGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        Name,
			Summary:     preview.Alpha.Summary,
			Description: "Early preview commands, for tenants invited. 'mgc preview list' shows the ones your tenant was given, and 'mgc preview install' adds them",
			IsInternal:  utils.BoolPtr(true),
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				project.GetGroup(),
			}
		},
	)
})
