/*
import "magalu.cloud/lib/products/container_registry/credentials"
*/
package credentials

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	List(configs ListConfigs) (result ListResult, err error)
	Password(configs PasswordConfigs) (result PasswordResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
