/*
import "magalu.cloud/lib/products/network/port/ports"
*/
package ports

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Attach(parameters AttachParameters, configs AttachConfigs) (result AttachResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	Detach(parameters DetachParameters, configs DetachConfigs) (result DetachResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}