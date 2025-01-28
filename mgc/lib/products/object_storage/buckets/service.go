/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/buckets"
*/
package buckets

import (
	"context"

	mgcClient "github.com/MagaluCloud/magalu/mgc/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	CreateContext(ctx context.Context, parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	Create(parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	ListContext(ctx context.Context, configs ListConfigs) (result ListResult, err error)
	List(configs ListConfigs) (result ListResult, err error)
	PublicUrlContext(ctx context.Context, parameters PublicUrlParameters, configs PublicUrlConfigs) (result PublicUrlResult, err error)
	PublicUrl(parameters PublicUrlParameters, configs PublicUrlConfigs) (result PublicUrlResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
