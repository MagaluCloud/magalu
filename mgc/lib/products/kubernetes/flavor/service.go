/*
import "magalu.cloud/lib/products/kubernetes/flavor"
*/
package flavor

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	ListContext(ctx context.Context, configs ListConfigs) (result ListResult, err error)
	List(configs ListConfigs) (result ListResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
