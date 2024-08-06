/*
import "magalu.cloud/lib/products/object_storage/objects"
*/
package objects

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Copy(parameters CopyParameters, configs CopyConfigs) (result CopyResult, err error)
	CopyAll(parameters CopyAllParameters, configs CopyAllConfigs) (result CopyAllResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	DeleteAll(parameters DeleteAllParameters, configs DeleteAllConfigs) (result DeleteAllResult, err error)
	Download(parameters DownloadParameters, configs DownloadConfigs) (result DownloadResult, err error)
	DownloadAll(parameters DownloadAllParameters, configs DownloadAllConfigs) (result DownloadAllResult, err error)
	Head(parameters HeadParameters, configs HeadConfigs) (result HeadResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	Move(parameters MoveParameters, configs MoveConfigs) (result MoveResult, err error)
	MoveDir(parameters MoveDirParameters, configs MoveDirConfigs) (result MoveDirResult, err error)
	Presign(parameters PresignParameters, configs PresignConfigs) (result PresignResult, err error)
	PublicUrl(parameters PublicUrlParameters, configs PublicUrlConfigs) (result PublicUrlResult, err error)
	Sync(parameters SyncParameters, configs SyncConfigs) (result SyncResult, err error)
	Upload(parameters UploadParameters, configs UploadConfigs) (result UploadResult, err error)
	UploadDir(parameters UploadDirParameters, configs UploadDirConfigs) (result UploadDirResult, err error)
	Versions(parameters VersionsParameters, configs VersionsConfigs) (result VersionsResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
