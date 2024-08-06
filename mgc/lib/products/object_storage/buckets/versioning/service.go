/*
import "magalu.cloud/lib/products/object_storage/buckets/versioning"
*/
package versioning

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Enable(parameters EnableParameters, configs EnableConfigs) (result EnableResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Suspend(parameters SuspendParameters, configs SuspendConfigs) (result SuspendResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
