/*
import "magalu.cloud/lib/products/block_storage/backups"
*/
package backups

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
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	Rename(parameters RenameParameters, configs RenameConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
