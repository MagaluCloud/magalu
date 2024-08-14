/*
import "magalu.cloud/lib/products/auth/api_key"
*/
package apiKey

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
	Get(parameters GetParameters) (result GetResult, err error)
	List(parameters ListParameters) (result ListResult, err error)
	Revoke(parameters RevokeParameters) (result RevokeResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
