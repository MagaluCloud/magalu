/*
import "magalu.cloud/lib/products/network/backoffice_conciliator"
*/
package backofficeConciliator

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	CreateContext(ctx context.Context, configs CreateConfigs) (result CreateResult, err error)
	Create(configs CreateConfigs) (result CreateResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
