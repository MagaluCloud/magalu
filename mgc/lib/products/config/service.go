/*
import "magalu.cloud/lib/products/config"
*/
package config

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Delete(parameters DeleteParameters) (result DeleteResult, err error)
	Get(parameters GetParameters) (result GetResult, err error)
	GetSchema(parameters GetSchemaParameters) (result GetSchemaResult, err error)
	List() (result ListResult, err error)
	Set(parameters SetParameters) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
