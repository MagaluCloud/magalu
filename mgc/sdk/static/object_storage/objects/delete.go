package objects

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/object_storage/s3"
)

type DeleteObjectParams struct {
	Destination string `json:"dst" jsonschema:"description=Path of the object to be deleted" example:"s3://bucket1/file1"`
}

func newDelete() core.Executor {
	exec := core.NewStaticExecute(
		"delete",
		"",
		"Delete an object from a bucket",
		Delete,
	)

	msg := "This command will delete the object at {{.parameters.dst}}, and it's result is NOT reversible."

	return core.NewConfirmableExecutor(
		exec,
		core.ConfirmPromptWithTemplate(msg),
	)
}

func newDeleteRequest(ctx context.Context, cfg s3.Config, pathURIs ...string) (*http.Request, error) {
	host := s3.BuildHost(cfg)
	url, err := url.JoinPath(host, pathURIs...)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
}

func Delete(ctx context.Context, params DeleteObjectParams, cfg s3.Config) (result core.Value, err error) {
	bucketURI, _ := strings.CutPrefix(params.Destination, s3.URIPrefix)
	req, err := newDeleteRequest(ctx, cfg, bucketURI)
	if err != nil {
		return nil, err
	}

	result, _, err = s3.SendRequest[core.Value](ctx, req)
	return
}
