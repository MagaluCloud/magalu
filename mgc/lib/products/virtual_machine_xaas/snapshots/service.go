/*
import "magalu.cloud/lib/products/virtual_machine_xaas/snapshots"
*/
package snapshots

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Update(parameters UpdateParameters, configs UpdateConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
