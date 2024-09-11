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
	DeleteContext(ctx context.Context, parameters DeleteParameters) (result DeleteResult, err error)
	Delete(parameters DeleteParameters) (result DeleteResult, err error)
	GetContext(ctx context.Context, parameters GetParameters) (result GetResult, err error)
	Get(parameters GetParameters) (result GetResult, err error)
	GetSchemaContext(ctx context.Context, parameters GetSchemaParameters) (result GetSchemaResult, err error)
	GetSchema(parameters GetSchemaParameters) (result GetSchemaResult, err error)
	ListContext(ctx context.Context) (result ListResult, err error)
	List() (result ListResult, err error)
	SetContext(ctx context.Context, parameters SetParameters) (result SetResult, err error)
	Set(parameters SetParameters) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
