/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/auth"
*/
package auth

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	AccessTokenContext(ctx context.Context, parameters AccessTokenParameters) (result AccessTokenResult, err error)
	AccessToken(parameters AccessTokenParameters) (result AccessTokenResult, err error)
	LoginContext(ctx context.Context, parameters LoginParameters) (result LoginResult, err error)
	Login(parameters LoginParameters) (result LoginResult, err error)
	LogoutContext(ctx context.Context, parameters LogoutParameters) (result LogoutResult, err error)
	Logout(parameters LogoutParameters) (result LogoutResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
