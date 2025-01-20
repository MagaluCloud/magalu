package objects

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/objects/acl"
	object_lock "github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/objects/object-lock"
)

var GetGroup = utils.NewLazyLoader[core.Grouper](func() core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "objects",
			Description: "Object operations for Object Storage API",
		},
		func() []core.Descriptor {
			return []core.Descriptor{
				acl.GetGroup(),         // object-storage objects acl
				getCopy(),              // object-storage objects copy
				getCopyAll(),           // object-storage objects copy-all
				getDelete(),            // object-storage objects delete
				getDeleteAll(),         // object-storage objects delete-all
				getDownload(),          // object-storage objects download
				getDownloadAll(),       // object-storage objects download-all
				getHead(),              // object-storage objects head
				getList(),              // object-storage objects list
				getMoveDir(),           // object-storage objects move-dir
				getMove(),              // object-storage objects move
				object_lock.GetGroup(), // object-storage objects object-lock
				getSync(),              // object-storage objects sync
				getUpload(),            // object-storage objects upload
				getUploadDir(),         // object-storage objects upload-dir
				getPresign(),           // object-storage objects presigned
				getPublicUrl(),         // object-storage objects public-url
				getVersions(),          // object-storage objects versions
			}
		},
	)
})
