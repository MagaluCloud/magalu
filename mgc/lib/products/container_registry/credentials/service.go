/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/container_registry/credentials"
*/
package credentials

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	ListContext(ctx context.Context, configs ListConfigs) (result ListResult, err error)
	List(configs ListConfigs) (result ListResult, err error)
	PasswordContext(ctx context.Context, configs PasswordConfigs) (result PasswordResult, err error)
	Password(configs PasswordConfigs) (result PasswordResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
