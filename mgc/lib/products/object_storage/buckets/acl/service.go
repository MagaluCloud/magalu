/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/buckets/acl"
*/
package acl

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	SetContext(ctx context.Context, parameters SetParameters, configs SetConfigs) (result SetResult, err error)
	Set(parameters SetParameters, configs SetConfigs) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
