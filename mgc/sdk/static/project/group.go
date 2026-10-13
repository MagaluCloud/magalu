package project

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "project",
			Summary: "Manage the projects of your tenant",
			Description: `Projects group the resources of a tenant. The selected project scopes the
Virtual Machine, Kubernetes and Block Storage requests. Listing and creating
projects talk to Magalu Cloud IAM and require a user login (see 'mgc auth login')`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getList(),
				getCreate(),
				getSet(),
				getUnset(),
				getCurrent(),
			}
		},
	)
})
