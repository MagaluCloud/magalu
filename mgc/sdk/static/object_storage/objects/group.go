package objects

import (
	"magalu.cloud/core"
)

func NewGroup() core.Grouper {
	return core.NewStaticGroup(
		"objects",
		"",
		"Object operations for Object Storage API",
		[]core.Descriptor{
			newDelete(),   // object-storage objects delete
			newDownload(), // object-storage objects download
			newList(),     // object-storage objects list
			newUpload(),   // object-storage objects upload
		},
	)
}
