/*
import "magalu.cloud/lib/products/network/subnets/subnets"
*/
package subnets

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
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	UpdateContext(ctx context.Context, parameters UpdateParameters, configs UpdateConfigs) (result UpdateResult, err error)
	Update(parameters UpdateParameters, configs UpdateConfigs) (result UpdateResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
