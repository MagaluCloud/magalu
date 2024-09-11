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
	CreateContext(ctx context.Context, parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	Create(parameters CreateParameters, configs CreateConfigs) (result CreateResult, err error)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	PasswordContext(ctx context.Context, parameters PasswordParameters, configs PasswordConfigs) (result PasswordResult, err error)
	Password(parameters PasswordParameters, configs PasswordConfigs) (result PasswordResult, err error)
	RebootContext(ctx context.Context, parameters RebootParameters, configs RebootConfigs) (err error)
	Reboot(parameters RebootParameters, configs RebootConfigs) (err error)
	RenameContext(ctx context.Context, parameters RenameParameters, configs RenameConfigs) (err error)
	Rename(parameters RenameParameters, configs RenameConfigs) (err error)
	RetypeContext(ctx context.Context, parameters RetypeParameters, configs RetypeConfigs) (err error)
	Retype(parameters RetypeParameters, configs RetypeConfigs) (err error)
	StartContext(ctx context.Context, parameters StartParameters, configs StartConfigs) (err error)
	Start(parameters StartParameters, configs StartConfigs) (err error)
	StopContext(ctx context.Context, parameters StopParameters, configs StopConfigs) (err error)
	Stop(parameters StopParameters, configs StopConfigs) (err error)
	SuspendContext(ctx context.Context, parameters SuspendParameters, configs SuspendConfigs) (err error)
	Suspend(parameters SuspendParameters, configs SuspendConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
