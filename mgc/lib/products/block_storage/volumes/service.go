/*
import "magalu.cloud/lib/products/block_storage/volumes"
*/
package volumes

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	AttachContext(ctx context.Context, parameters AttachParameters, configs AttachConfigs) (err error)
	//Attach(	parameters AttachParameters, configs AttachConfigs,) ( err error,)
	CreateContext(ctx context.Context, parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	//Create(	parameters CreateParameters, configs CreateConfigs,) ( result CreateResult, err error,)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (err error)
	//Delete(	parameters DeleteParameters, configs DeleteConfigs,) ( err error,)
	DetachContext(ctx context.Context, parameters DetachParameters, configs DetachConfigs) (err error)
	//Detach(	parameters DetachParameters, configs DetachConfigs,) ( err error,)
	ExtendContext(ctx context.Context, parameters ExtendParameters, configs ExtendConfigs) (err error)
	//Extend(	parameters ExtendParameters, configs ExtendConfigs,) ( err error,)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	//Get(	parameters GetParameters, configs GetConfigs,) ( result GetResult, err error,)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	//List(	parameters ListParameters, configs ListConfigs,) ( result ListResult, err error,)
	RenameContext(ctx context.Context, parameters RenameParameters, configs RenameConfigs) (err error)
	//Rename(	parameters RenameParameters, configs RenameConfigs,) ( err error,)
	RetypeContext(ctx context.Context, parameters RetypeParameters, configs RetypeConfigs) (err error)
	//Retype(	parameters RetypeParameters, configs RetypeConfigs,) ( err error,)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
