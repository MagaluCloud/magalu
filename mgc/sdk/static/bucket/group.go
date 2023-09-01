package bucket

import (
	"magalu.cloud/core"
)

func NewGroup() core.Grouper {
	return core.NewStaticGroup(
		"bucket",
		"",
		"Bucket operations for Object Storage API",
		[]core.Descriptor{
			newCreate(),
			newDelete(),
			newDeleteObject(),
			newList(),
			newListObjects(),
			newUpload(),
		},
	)
}
