/*
import "magalu.cloud/lib/products/virtual_machine_xaas/xaas_keypairs"
*/
package xaasKeypairs

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Create(parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	DeleteKeypairName(parameters DeleteKeypairNameParameters, configs DeleteKeypairNameConfigs) (err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
