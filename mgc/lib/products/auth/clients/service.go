/*
import "magalu.cloud/lib/products/auth/clients"
*/
package clients

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Create(parameters CreateParameters) (result CreateResult, err error)
	List() (result ListResult, err error)
	Update(parameters UpdateParameters) (result UpdateResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
