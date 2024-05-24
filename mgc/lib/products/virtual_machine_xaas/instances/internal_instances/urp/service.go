/*
import "magalu.cloud/lib/products/virtual_machine_xaas/instances/internal_instances/urp"
*/
package urp

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Update(parameters UpdateParameters, configs UpdateConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
