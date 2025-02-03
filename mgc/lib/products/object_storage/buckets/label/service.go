/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/buckets/label"
*/
package label

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	SetContext(ctx context.Context, parameters SetParameters, configs SetConfigs) (result SetResult, err error)
	Set(parameters SetParameters, configs SetConfigs) (result SetResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
