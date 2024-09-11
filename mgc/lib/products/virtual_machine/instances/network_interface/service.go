/*
import "magalu.cloud/lib/products/virtual_machine/instances/network_interface"
*/
package networkInterface

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	AttachContext(ctx context.Context, parameters AttachParameters, configs AttachConfigs) (err error)
	Attach(parameters AttachParameters, configs AttachConfigs) (err error)
	DetachContext(ctx context.Context, parameters DetachParameters, configs DetachConfigs) (err error)
	Detach(parameters DetachParameters, configs DetachConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
