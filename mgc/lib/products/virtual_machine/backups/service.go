/*
import "magalu.cloud/lib/products/virtual_machine/backups"
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
	CreateContext(ctx context.Context, parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	Create(parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	RenameContext(ctx context.Context, parameters RenameParameters, configs RenameConfigs) (err error)
	Rename(parameters RenameParameters, configs RenameConfigs) (err error)
	RestoreContext(ctx context.Context, parameters RestoreParameters, configs RestoreConfigs) (result RestoreResult, err error)
	Restore(parameters RestoreParameters, configs RestoreConfigs) (result RestoreResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
