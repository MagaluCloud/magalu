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
	Copy(parameters CopyParameters, configs CopyConfigs) (err error)
}

func NewService(ctx context.Context, client *mgcClient.Client) Service {
	return &service{ctx, client}
}
