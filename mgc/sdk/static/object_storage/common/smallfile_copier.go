package common

import (
	"context"

	mgcSchemaPkg "magalu.cloud/core/schema"
)

type smallFileCopier struct {
	cfg          Config
	src          mgcSchemaPkg.URI
	dst          mgcSchemaPkg.URI
	version      string
	storageClass string
}

var _ copier = (*smallFileCopier)(nil)

func (u *smallFileCopier) Copy(ctx context.Context) error {
	req, err := newCopyRequest(ctx, u.cfg, u.src, u.dst, u.version)
	if err != nil {
		return err
	}

	if u.storageClass != "" {
		req.Header.Set("X-Amz-Storage-Class", u.storageClass)
	}

	resp, err := SendRequest(ctx, req)
	if err != nil {
		return err
	}

	return ExtractErr(resp, req)
}
