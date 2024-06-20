/*
import "magalu.cloud/lib/products/virtual_machine_xaas/keypairs"
*/
package keypairs

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
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
