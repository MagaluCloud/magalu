package alpha

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
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
			Summary:     "Early preview commands",
			Description: "Early preview commands. They may change or stop working, and the API refuses them unless your tenant was invited",
			IsInternal:  utils.BoolPtr(true),
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				project.GetGroup(),
			}
		},
	)
})
