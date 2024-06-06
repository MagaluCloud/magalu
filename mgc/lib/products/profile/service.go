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
	Delete(parameters DeleteParameters) (result DeleteResult, err error)
	Get() (result GetResult, err error)
	List() (result ListResult, err error)
	Set(parameters SetParameters) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
