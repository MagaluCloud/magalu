package objectstorage

import (
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

var GetGroup = utils.NewLazyLoader(func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:    "object-storage",
			Summary: "Credentials used for Object Storage",
			Description: `Object Storage uses a different set of credentials when compared
to normal HTTP requests. Two keys are needed, the 'SecretKey' and
'AccessKeyId'. Instructions on how to create these can be found
here: https://id.magalu.com/api-keys`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				getSet(),
				getGet(),
			}
		},
	)
})
