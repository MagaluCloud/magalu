/*
import "magalu.cloud/lib/products/virtual_machine/instances"
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
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	Password(parameters PasswordParameters, configs PasswordConfigs) (result PasswordResult, err error)
	Reboot(parameters RebootParameters, configs RebootConfigs) (err error)
	Rename(parameters RenameParameters, configs RenameConfigs) (err error)
	Retype(parameters RetypeParameters, configs RetypeConfigs) (err error)
	Start(parameters StartParameters, configs StartConfigs) (err error)
	Stop(parameters StopParameters, configs StopConfigs) (err error)
	Suspend(parameters SuspendParameters, configs SuspendConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
