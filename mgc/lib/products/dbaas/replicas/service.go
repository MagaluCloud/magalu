/*
import "magalu.cloud/lib/products/dbaas/replicas"
*/
package replicas

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
	//Create(	parameters CreateParameters, configs CreateConfigs,) ( result CreateResult, err error,)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (err error)
	//Delete(	parameters DeleteParameters, configs DeleteConfigs,) ( err error,)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	//Get(	parameters GetParameters, configs GetConfigs,) ( result GetResult, err error,)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	//List(	parameters ListParameters, configs ListConfigs,) ( result ListResult, err error,)
	ResizeContext(ctx context.Context, parameters ResizeParameters, configs ResizeConfigs) (result ResizeResult, err error)
	//Resize(	parameters ResizeParameters, configs ResizeConfigs,) ( result ResizeResult, err error,)
	StartContext(ctx context.Context, parameters StartParameters, configs StartConfigs) (result StartResult, err error)
	//Start(	parameters StartParameters, configs StartConfigs,) ( result StartResult, err error,)
	StopContext(ctx context.Context, parameters StopParameters, configs StopConfigs) (result StopResult, err error)
	//Stop(	parameters StopParameters, configs StopConfigs,) ( result StopResult, err error,)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
