/*
import "magalu.cloud/lib/products/workspace"
*/
package workspace

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	CreateContext(ctx context.Context, parameters CreateParameters) (result CreateResult, err error)
	Create(parameters CreateParameters) (result CreateResult, err error)
	DeleteContext(ctx context.Context, parameters DeleteParameters) (result DeleteResult, err error)
	Delete(parameters DeleteParameters) (result DeleteResult, err error)
	GetContext(ctx context.Context) (result GetResult, err error)
	Get() (result GetResult, err error)
	ListContext(ctx context.Context) (result ListResult, err error)
	List() (result ListResult, err error)
	SetContext(ctx context.Context, parameters SetParameters) (result SetResult, err error)
	Set(parameters SetParameters) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
