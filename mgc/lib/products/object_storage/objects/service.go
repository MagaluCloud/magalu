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
	CopyContext(ctx context.Context, parameters CopyParameters, configs CopyConfigs) (result CopyResult, err error)
	Copy(parameters CopyParameters, configs CopyConfigs) (result CopyResult, err error)
	CopyAllContext(ctx context.Context, parameters CopyAllParameters, configs CopyAllConfigs) (result CopyAllResult, err error)
	CopyAll(parameters CopyAllParameters, configs CopyAllConfigs) (result CopyAllResult, err error)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	DeleteAllContext(ctx context.Context, parameters DeleteAllParameters, configs DeleteAllConfigs) (result DeleteAllResult, err error)
	DeleteAll(parameters DeleteAllParameters, configs DeleteAllConfigs) (result DeleteAllResult, err error)
	DownloadContext(ctx context.Context, parameters DownloadParameters, configs DownloadConfigs) (result DownloadResult, err error)
	Download(parameters DownloadParameters, configs DownloadConfigs) (result DownloadResult, err error)
	DownloadAllContext(ctx context.Context, parameters DownloadAllParameters, configs DownloadAllConfigs) (result DownloadAllResult, err error)
	DownloadAll(parameters DownloadAllParameters, configs DownloadAllConfigs) (result DownloadAllResult, err error)
	HeadContext(ctx context.Context, parameters HeadParameters, configs HeadConfigs) (result HeadResult, err error)
	Head(parameters HeadParameters, configs HeadConfigs) (result HeadResult, err error)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	MoveContext(ctx context.Context, parameters MoveParameters, configs MoveConfigs) (result MoveResult, err error)
	Move(parameters MoveParameters, configs MoveConfigs) (result MoveResult, err error)
	MoveDirContext(ctx context.Context, parameters MoveDirParameters, configs MoveDirConfigs) (result MoveDirResult, err error)
	MoveDir(parameters MoveDirParameters, configs MoveDirConfigs) (result MoveDirResult, err error)
	PresignContext(ctx context.Context, parameters PresignParameters, configs PresignConfigs) (result PresignResult, err error)
	Presign(parameters PresignParameters, configs PresignConfigs) (result PresignResult, err error)
	PublicUrlContext(ctx context.Context, parameters PublicUrlParameters, configs PublicUrlConfigs) (result PublicUrlResult, err error)
	PublicUrl(parameters PublicUrlParameters, configs PublicUrlConfigs) (result PublicUrlResult, err error)
	SyncContext(ctx context.Context, parameters SyncParameters, configs SyncConfigs) (result SyncResult, err error)
	Sync(parameters SyncParameters, configs SyncConfigs) (result SyncResult, err error)
	UploadContext(ctx context.Context, parameters UploadParameters, configs UploadConfigs) (result UploadResult, err error)
	Upload(parameters UploadParameters, configs UploadConfigs) (result UploadResult, err error)
	UploadDirContext(ctx context.Context, parameters UploadDirParameters, configs UploadDirConfigs) (result UploadDirResult, err error)
	UploadDir(parameters UploadDirParameters, configs UploadDirConfigs) (result UploadDirResult, err error)
	VersionsContext(ctx context.Context, parameters VersionsParameters, configs VersionsConfigs) (result VersionsResult, err error)
	Versions(parameters VersionsParameters, configs VersionsConfigs) (result VersionsResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
