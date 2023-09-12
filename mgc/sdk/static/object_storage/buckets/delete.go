package buckets

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/sdk/static/object_storage/objects"
	"magalu.cloud/sdk/static/object_storage/s3"
)

var deleteBucketsLogger *zap.SugaredLogger

func deleteLogger() *zap.SugaredLogger {
	if deleteBucketsLogger == nil {
		deleteBucketsLogger = logger().Named("delete")
	}
	return deleteBucketsLogger
}

type deleteParams struct {
	Name string `json:"name" jsonschema:"description=Name of the bucket to be deleted"`
}

type deleteObjectsError struct {
	errorMap map[string]error
}

func (o deleteObjectsError) Error() string {
	var errorMsg string
	for file, err := range o.errorMap {
		errorMsg += fmt.Sprintf("%s - %s, ", file, err)
	}
	// Remove trailing `, `
	if len(errorMsg) != 0 {
		errorMsg = errorMsg[:len(errorMsg)-2]
	}
	return fmt.Sprintf("failed to delete objects from bucket: %s", errorMsg)
}

func (o deleteObjectsError) Add(uri string, err error) {
	o.errorMap[uri] = err
}

func (o deleteObjectsError) HasError() bool {
	return len(o.errorMap) != 0
}

func NewDeleteObjectsError() deleteObjectsError {
	return deleteObjectsError{
		errorMap: make(map[string]error),
	}
}

func newDelete() core.Executor {
	executor := core.NewStaticExecute(
		"delete",
		"",
		"Delete a bucket",
		delete,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Value) string {
		return "template=Deleted bucket {{.name}}\n"
	})
}

func newDeleteRequest(ctx context.Context, region string, pathURIs ...string) (*http.Request, error) {
	host := s3.BuildHost(region)
	url, err := url.JoinPath(host, pathURIs...)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
}

func delete(ctx context.Context, params deleteParams, cfg s3.Config) (core.Value, error) {
	objs, err := objects.List(ctx, objects.ListObjectsParams{Destination: params.Name}, cfg)
	if err != nil {
		return nil, err
	}

	objErr := NewDeleteObjectsError()
	makeRequest := func(uri string, done chan<- core.Value) {
		_, err := objects.Delete(
			ctx,
			objects.DeleteObjectParams{Destination: uri},
			cfg,
		)

		if err != nil {
			objErr.Add(uri, err)
		} else {
			deleteLogger().Infof("Deleted %s%s", s3.URIPrefix, uri)
		}

		done <- true
	}

	done := make(chan core.Value)
	defer close(done)

	for _, obj := range objs.Contents {
		objURI := path.Join(params.Name, obj.Key)
		go makeRequest(objURI, done)
	}

	for range objs.Contents {
		<-done
	}

	if objErr.HasError() {
		return nil, objErr
	}

	req, err := newDeleteRequest(ctx, cfg.Region, params.Name)
	if err != nil {
		return nil, err
	}

	_, err = s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, (*core.Value)(nil))
	if err != nil {
		return nil, err
	}

	return params, nil
}
