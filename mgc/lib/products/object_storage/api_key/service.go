/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/api_key"
*/
package apiKey

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
	CurrentContext(ctx context.Context) (result CurrentResult, err error)
	Current() (result CurrentResult, err error)
	GetContext(ctx context.Context, parameters GetParameters) (result GetResult, err error)
	Get(parameters GetParameters) (result GetResult, err error)
	ListContext(ctx context.Context) (result ListResult, err error)
	List() (result ListResult, err error)
	RevokeContext(ctx context.Context, parameters RevokeParameters) (result RevokeResult, err error)
	Revoke(parameters RevokeParameters) (result RevokeResult, err error)
	SetContext(ctx context.Context, parameters SetParameters) (result SetResult, err error)
	Set(parameters SetParameters) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
