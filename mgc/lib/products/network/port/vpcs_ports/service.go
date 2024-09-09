/*
import "magalu.cloud/lib/products/network/port/vpcs_ports"
*/
package vpcsPorts

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
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	//List(	parameters ListParameters, configs ListConfigs,) ( result ListResult, err error,)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
