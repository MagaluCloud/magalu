/*
import "magalu.cloud/lib/products/network/internal"
*/
package internal

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	ConciliatorContext(ctx context.Context, configs ConciliatorConfigs) (result ConciliatorResult, err error)
	Conciliator(configs ConciliatorConfigs) (result ConciliatorResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
