/*
import "magalu.cloud/lib/products/container_registry/images"
*/
package images

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (err error)
	//Delete(	parameters DeleteParameters, configs DeleteConfigs,) ( err error,)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	//Get(	parameters GetParameters, configs GetConfigs,) ( result GetResult, err error,)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	//List(	parameters ListParameters, configs ListConfigs,) ( result ListResult, err error,)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
