/*
import "magalu.cloud/lib/products/kubernetes/info"
*/
package info

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Flavors(configs FlavorsConfigs) (result FlavorsResult, err error)
	Versions(configs VersionsConfigs) (result VersionsResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
