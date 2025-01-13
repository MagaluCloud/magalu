/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/auth/clients"
*/
package clients

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	CreateContext(ctx context.Context, parameters CreateParameters) (result CreateResult, err error)
	Create(parameters CreateParameters) (result CreateResult, err error)
	ListContext(ctx context.Context) (result ListResult, err error)
	List() (result ListResult, err error)
	UpdateContext(ctx context.Context, parameters UpdateParameters) (result UpdateResult, err error)
	Update(parameters UpdateParameters) (result UpdateResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
