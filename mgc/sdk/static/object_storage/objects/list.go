package objects

import (
	"context"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/pipeline"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type ListObjectsParams struct {
	Destination string `json:"dst" jsonschema:"description=Path of the bucket to list objects from" example:"s3://bucket1/"`
}

type prefix struct {
	Path string `xml:"Prefix"`
}

type ListObjectsResponse struct {
	Name           string           `xml:"Name"`
	Contents       []*BucketContent `xml:"Contents"`
	CommonPrefixes []*prefix        `xml:"CommonPrefixes" json:"SubDirectories"`
}

type BucketContent struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ContentSize  int    `xml:"Size"`
}

func (b *BucketContent) ModTime() time.Time {
	modTime, err := time.Parse(time.RFC3339, b.LastModified)
	if err != nil {
		listObjectsLogger().Errorw("Failed to parse time", "err", err)
		modTime = time.Time{}
	}
	return modTime
}

func (b *BucketContent) Mode() fs.FileMode {
	return utils.FILE_PERMISSION
}

func (b *BucketContent) Size() int64 {
	return int64(b.ContentSize)
}

func (b *BucketContent) Sys() any {
	return nil
}

func (b *BucketContent) Info() (fs.FileInfo, error) {
	return b, nil
}

func (b *BucketContent) IsDir() bool {
	return false
}

func (b *BucketContent) Name() string {
	return b.Key
}

func (b *BucketContent) Type() fs.FileMode {
	return utils.FILE_PERMISSION
}

var _ fs.DirEntry = (*BucketContent)(nil)
var _ fs.FileInfo = (*BucketContent)(nil)

var listObjectsLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("list")
})

func newListRequest(ctx context.Context, cfg common.Config, bucket string) (*http.Request, error) {
	parsedUrl, err := parseURL(cfg, bucket)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
}

var getList = utils.NewLazyLoader[core.Executor](newList)

func newList() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all objects from a bucket",
		},
		List,
	)
}

func parseURL(cfg common.Config, bucketURI string) (*url.URL, error) {
	// Bucket URI cannot end in '/' as this makes it search for a
	// non existing directory
	bucketURI = strings.TrimSuffix(bucketURI, "/")
	dirs := strings.Split(bucketURI, "/")
	path, err := url.JoinPath(common.BuildHost(cfg), dirs[0])
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	if len(dirs) <= 1 {
		return u, nil
	}
	q := u.Query()
	delimiter := "/"
	prefixQ := strings.Join(dirs[1:], delimiter)
	lastChar := string(prefixQ[len(prefixQ)-1])
	if lastChar != delimiter {
		prefixQ += delimiter
	}
	q.Set("prefix", prefixQ)
	q.Set("delimiter", delimiter)
	q.Set("encoding-type", "url")
	u.RawQuery = q.Encode()
	return u, nil
}

func List(ctx context.Context, params ListObjectsParams, cfg common.Config) (result ListObjectsResponse, err error) {
	objChan := ListGenerator(ctx, params, cfg)

	entries, err := pipeline.SliceItemConsumer[[]pipeline.WalkDirEntry](ctx, objChan)
	if err != nil {
		return result, err
	}

	var contents []*BucketContent
	for _, entry := range entries {
		content, ok := entry.DirEntry.(*BucketContent)
		if !ok {
			continue
		}
		contents = append(contents, content)

		if entry.Err != nil {
			return result, entry.Err
		}
	}

	result = ListObjectsResponse{
		Contents: contents,
	}
	return result, nil
}

func ListGenerator(ctx context.Context, params ListObjectsParams, cfg common.Config) (outputChan <-chan pipeline.WalkDirEntry) {
	ch := make(chan pipeline.WalkDirEntry)
	outputChan = ch

	generator := func() {
		defer func() {
			listObjectsLogger().Info("closing output channel")
			close(ch)
		}()

		bucket, _ := strings.CutPrefix(params.Destination, common.URIPrefix)
		req, err := newListRequest(ctx, cfg, bucket)
		if err != nil {
			listObjectsLogger().Errorw("newListRequest() failed", "err", err)
			ch <- pipeline.WalkDirEntry{Err: err}
			return
		}

		result, _, err := common.SendRequest[ListObjectsResponse](ctx, req)
		if err != nil {
			listObjectsLogger().Errorw("s3.SendRequest() failed", "err", err)
			ch <- pipeline.WalkDirEntry{Err: err}
			return
		}

		for _, content := range result.Contents {
			dirEntry := pipeline.WalkDirEntry{
				Path:     path.Join(params.Destination, content.Key),
				DirEntry: content,
				Err:      err,
			}

			select {
			case <-ctx.Done():
				listObjectsLogger().Debugw("context.Done()", "err", ctx.Err())
				return
			case ch <- dirEntry:
			}
		}
		listObjectsLogger().Info("finished reading contents")
	}

	listObjectsLogger().Info("list generation start")
	go generator()
	return
}
