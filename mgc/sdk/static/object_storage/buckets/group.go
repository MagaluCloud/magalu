package buckets

import (
	"magalu.cloud/core"
)

func NewGroup() core.Grouper {
	return core.NewStaticGroup(
		"buckets",
		"",
		"Bucket operations for Object Storage API",
		[]core.Descriptor{
			newCreate(), // object-storage buckets create
			newDelete(), // object-storage buckets delete
			newList(),   // object-storage buckets list
		},
	)
}
