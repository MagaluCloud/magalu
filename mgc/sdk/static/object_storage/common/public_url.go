package common

import (
	"context"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

type PublicUrlResult struct {
	URL mgcSchemaPkg.URI `json:"url"`
}

func PublicUrl(ctx context.Context, cfg Config, dst mgcSchemaPkg.URI) (url *PublicUrlResult, err error) {
	resourceUrl, err := BuildBucketHostWithPath(cfg, NewBucketNameFromURI(dst), dst.Path())
	if err != nil {
		return
	}

	return &PublicUrlResult{
		URL: mgcSchemaPkg.URI(resourceUrl),
	}, nil
}
