/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/kubernetes/info"
*/
package info

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	FlavorsContext(ctx context.Context, configs FlavorsConfigs) (result FlavorsResult, err error)
	Flavors(configs FlavorsConfigs) (result FlavorsResult, err error)
	VersionsContext(ctx context.Context, configs VersionsConfigs) (result VersionsResult, err error)
	Versions(configs VersionsConfigs) (result VersionsResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
