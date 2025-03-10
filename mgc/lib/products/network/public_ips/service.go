/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/network/public_ips"
*/
package publicIps

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	AttachContext(ctx context.Context, parameters AttachParameters, configs AttachConfigs) (err error)
	Attach(parameters AttachParameters, configs AttachConfigs) (err error)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	DetachContext(ctx context.Context, parameters DetachParameters, configs DetachConfigs) (err error)
	Detach(parameters DetachParameters, configs DetachConfigs) (err error)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	ListContext(ctx context.Context, configs ListConfigs) (result ListResult, err error)
	List(configs ListConfigs) (result ListResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
