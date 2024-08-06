/*
import "magalu.cloud/lib/products/object_storage/buckets"
*/
package buckets

import (
	"context"

	mgcClient "magalu.cloud/lib"
)

type service struct {
	ctx    context.Context
	client *mgcClient.Client
}

type Service interface {
	Create(parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (result DeleteResult, err error)
	List(configs ListConfigs) (result ListResult, err error)
	PublicUrl(parameters PublicUrlParameters, configs PublicUrlConfigs) (result PublicUrlResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
