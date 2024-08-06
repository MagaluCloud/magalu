/*
import "magalu.cloud/lib/products/auth/tenant"
*/
package tenant

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Current() (result CurrentResult, err error)
	List() (result ListResult, err error)
	Set(parameters SetParameters) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
