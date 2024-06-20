/*
import "magalu.cloud/lib/products/virtual_machine_xaas/healthchecks"
*/
package healthchecks

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Health(configs HealthConfigs) (err error)
	Healthcheck(configs HealthcheckConfigs) (result HealthcheckResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
