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
	Attach(parameters AttachParameters, configs AttachConfigs) (err error)
	Create(parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	Detach(parameters DetachParameters, configs DetachConfigs) (err error)
	Extend(parameters ExtendParameters, configs ExtendConfigs) (err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	Rename(parameters RenameParameters, configs RenameConfigs) (err error)
	Retype(parameters RetypeParameters, configs RetypeConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
