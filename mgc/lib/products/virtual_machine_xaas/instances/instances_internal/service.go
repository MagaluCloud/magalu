/*
import "magalu.cloud/lib/products/virtual_machine_xaas/instances/instances_internal"
*/
package instancesInternal

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
	Urp(parameters UrpParameters, configs UrpConfigs) (result UrpResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
