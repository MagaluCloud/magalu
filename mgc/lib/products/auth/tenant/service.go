/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/auth/tenant"
*/
package tenant

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	CurrentContext(ctx context.Context) (result CurrentResult, err error)
	Current() (result CurrentResult, err error)
	ListContext(ctx context.Context) (result ListResult, err error)
	List() (result ListResult, err error)
	SetContext(ctx context.Context, parameters SetParameters) (result SetResult, err error)
	Set(parameters SetParameters) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
