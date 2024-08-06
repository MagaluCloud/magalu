/*
import "magalu.cloud/lib/products/network/rule/rules"
*/
package rules

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
