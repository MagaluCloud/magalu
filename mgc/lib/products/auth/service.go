/*
import "magalu.cloud/lib/products/auth"
*/
package auth

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	AccessTokenContext(ctx context.Context, parameters AccessTokenParameters) (result AccessTokenResult, err error)
	//AccessToken(	parameters AccessTokenParameters,) ( result AccessTokenResult, err error,)
	LoginContext(ctx context.Context, parameters LoginParameters) (result LoginResult, err error)
	//Login(	parameters LoginParameters,) ( result LoginResult, err error,)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
