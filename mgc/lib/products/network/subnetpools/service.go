/*
import "github.com/MagaluCloud/magalu/mgc/lib/products/network/subnetpools"
*/
package subnetpools

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
	CreateBookCidrContext(ctx context.Context, parameters CreateBookCidrParameters, configs CreateBookCidrConfigs) (result CreateBookCidrResult, err error)
	CreateBookCidr(parameters CreateBookCidrParameters, configs CreateBookCidrConfigs) (result CreateBookCidrResult, err error)
	CreateUnbookCidrContext(ctx context.Context, parameters CreateUnbookCidrParameters, configs CreateUnbookCidrConfigs) (err error)
	CreateUnbookCidr(parameters CreateUnbookCidrParameters, configs CreateUnbookCidrConfigs) (err error)
	DeleteContext(ctx context.Context, parameters DeleteParameters, configs DeleteConfigs) (err error)
	Delete(parameters DeleteParameters, configs DeleteConfigs) (err error)
	GetContext(ctx context.Context, parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	Get(parameters GetParameters, configs GetConfigs) (result GetResult, err error)
	ListContext(ctx context.Context, parameters ListParameters, configs ListConfigs) (result ListResult, err error)
	List(parameters ListParameters, configs ListConfigs) (result ListResult, err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
