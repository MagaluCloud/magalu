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
	Attach(parameters AttachParameters, configs AttachConfigs) (err error)
	Detach(parameters DetachParameters, configs DetachConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
