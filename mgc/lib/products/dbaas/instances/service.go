/*
import "magalu.cloud/lib/products/dbaas/instances"
*/
package instances

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
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	Resize(parameters ResizeParameters, configs ResizeConfigs) (result ResizeResult, err error)
	Restores(parameters RestoresParameters, configs RestoresConfigs) (result RestoresResult, err error)
	Start(parameters StartParameters, configs StartConfigs) (result StartResult, err error)
	Stop(parameters StopParameters, configs StopConfigs) (result StopResult, err error)
	Update(parameters UpdateParameters, configs UpdateConfigs) (result UpdateResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
