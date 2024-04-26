/*
import "magalu.cloud/lib/products/object_storage/objects/acl"
*/
package acl

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Set(parameters SetParameters, configs SetConfigs) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
