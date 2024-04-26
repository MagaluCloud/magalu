/*
import "magalu.cloud/lib/products/profile"
*/
package profile

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
	Current() (result CurrentResult, err error)
	Delete(parameters DeleteParameters) (result DeleteResult, err error)
	List() (result ListResult, err error)
	SetCurrent(parameters SetCurrentParameters) (result SetCurrentResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
