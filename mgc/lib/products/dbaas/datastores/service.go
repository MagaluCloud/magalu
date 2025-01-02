/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/dbaas/datastores"
*/
package datastores

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
