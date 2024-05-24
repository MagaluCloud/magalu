/*
import "magalu.cloud/lib/products/virtual_machine_xaas/images"
*/
package images

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
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	List(configs ListConfigs) (result ListResult, err error)
	Rename(parameters RenameParameters, configs RenameConfigs) (err error)
	Update(parameters UpdateParameters, configs UpdateConfigs) (err error)
	Urp(parameters UrpParameters, configs UrpConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
