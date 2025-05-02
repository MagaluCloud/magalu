package common

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

var listObjectsLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("list")
})

type ListObjectsParams struct {
	Destination      mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Path of the bucket to list objects from,example=bucket1" mgc:"positional"`
	PaginationParams `json:",squash"` // nolint
	Recursive        bool             `json:"recursive,omitempty" jsonschema:"description=List folders and subfolders,default=false"`
}

type PaginationParams struct {
	MaxItems          int    `json:"max-items,omitempty" jsonschema:"description=Limit of items to be listed,default=1000,minimum=1,example=1000,required"`
	ContinuationToken string `json:"continuation-token,omitempty" jsonschema:"description=Token of result page to continue from"`
}

type Prefix struct {
	Path string `xml:"Prefix"`
}

func (p *Prefix) ModTime() time.Time {
	var modTime time.Time
	return modTime
}

func (p *Prefix) Mode() fs.FileMode {
	return utils.DIR_PERMISSION | fs.ModeDir
}

func (p *Prefix) Size() int64 {
	return 0
}

func (p *Prefix) Sys() any {
	return nil
}

func (p *Prefix) Info() (fs.FileInfo, error) {
	return p, nil
}

func (p *Prefix) IsDir() bool {
	return true
}

func (p *Prefix) Name() string {
	return path.Base(path.Dir(p.Path))
}

func (p *Prefix) Type() fs.FileMode {
	return utils.DIR_PERMISSION | fs.ModeDir
}

var _ fs.DirEntry = (*Prefix)(nil)
var _ fs.FileInfo = (*Prefix)(nil)

type listObjectsRequestResponse struct {
	Name                   string           `xml:"Name"`
	Contents               []*BucketContent `xml:"Contents"`
	CommonPrefixes         []*Prefix        `xml:"CommonPrefixes" json:"SubDirectories"`
	paginationResponseInfo `json:",squash"` // nolint
}

type paginationResponseInfo struct {
	NextContinuationToken string `xml:"NextContinuationToken"`
	IsTruncated           bool   `xml:"IsTruncated"`
}

type BucketContent struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ContentSize  int64  `xml:"Size"`
	StorageClass string `xml:"StorageClass"`
}

type BucketContentDirEntry = *pipeline.SimpleWalkDirEntry[*BucketContent]

func (b *BucketContent) ModTime() time.Time {
	modTime, err := time.Parse(time.RFC3339, b.LastModified)
	if err != nil {
		listObjectsLogger().Named("BucketContent.ModTime()").Errorw("failed to parse time", "err", err, "key", b.Key, "lastModified", b.LastModified)
		modTime = time.Time{}
	}
	return modTime
}

func (b *BucketContent) Mode() fs.FileMode {
	return utils.FILE_PERMISSION
}

func (b *BucketContent) Size() int64 {
	return b.ContentSize
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
	return path.Base(b.Key)
}

func (b *BucketContent) Type() fs.FileMode {
	return utils.FILE_PERMISSION
}

var _ fs.DirEntry = (*BucketContent)(nil)
var _ fs.FileInfo = (*BucketContent)(nil)

// FIXME: The way we handle sigv4 in here (and probably other places) is really brittle.
//
//  1. We have to manually create the querystring, escaping spaces with '%20' instead of '+' (as the url.QueryEscape() would do)
//
//  2. sigv4 is dependent on the order of the query string params, this forces us to sort the query string params
//     alphabetically (also manaully)
//
//     The module that handles signature should actually take care of all of this under the hood so that this logic is not
//     repeated every time we need to make a signed API call.
//
// More info on https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html
func newListRequest(ctx context.Context, cfg Config, bucketURI mgcSchemaPkg.URI, page PaginationParams, recursive bool) (*http.Request, error) {
	finalUrl, err := buildListRequestURL(cfg, bucketURI)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	var queryStringParts []string

	if prefix := bucketURI.Path(); prefix != "" {
		lastChar := string(prefix[len(prefix)-1])
		if lastChar != delimiter {
			prefix += delimiter
		}

		queryEscapedValue := url.QueryEscape(prefix)

		// How for the "fun" part: the aws uri encoding scheme is not the same as go's.
		//
		// From the docs:
		// URI encode every byte. UriEncode() must enforce the following rules:
		//
		//   - URI encode every byte except the unreserved characters: 'A'-'Z', 'a'-'z', '0'-'9', '-', '.', '_', and '~'.
		//   - The space character is a reserved character and must be encoded as "%20" (and not as "+").
		//   - Each URI encoded byte is formed by a '%' and the two-digit hexadecimal value of the byte.
		//   - Letters in the hexadecimal value must be uppercase, for example "%1A".
		//   - Encode the forward slash character, '/', everywhere except in the object key name. For example, if the object key name is photos/Jan/sample.jpg, the forward slash in the key name is not encoded.
		//
		// Source: https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html#example-signature-calculations
		awsEscapedValue := strings.ReplaceAll(queryEscapedValue, "+", "%20")
		awsEscapedValue = strings.ReplaceAll(awsEscapedValue, "*", "%2A")
		awsEscapedValue = strings.ReplaceAll(awsEscapedValue, "%7E", "~")
		queryStringParts = append(queryStringParts, "prefix="+awsEscapedValue)
	}

	queryStringParts = append(queryStringParts, "list-type=2")
	if page.ContinuationToken != "" {
		queryStringParts = append(queryStringParts, "continuation-token="+url.QueryEscape(page.ContinuationToken))
	}

	if page.MaxItems <= 0 {
		return nil, core.UsageError{Err: fmt.Errorf("invalid item limit MaxItems, must be higher than zero: %d", page.MaxItems)}
	} else if page.MaxItems > ApiLimitMaxItems {
		page.MaxItems = ApiLimitMaxItems
	}

	queryStringParts = append(queryStringParts, "max-keys="+fmt.Sprint(page.MaxItems))
	if !recursive {
		queryStringParts = append(queryStringParts, "delimiter="+url.QueryEscape(delimiter))
	}

	sort.Strings(queryStringParts)
	finalUrl.RawQuery = strings.Join(queryStringParts, "&")

	return http.NewRequestWithContext(ctx, http.MethodGet, finalUrl.String(), nil)
}

func buildListRequestURL(cfg Config, bucketURI mgcSchemaPkg.URI) (*url.URL, error) {
	u, err := BuildBucketHostURL(cfg, NewBucketNameFromURI(bucketURI))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func ListGenerator(ctx context.Context, params ListObjectsParams, cfg Config, onNewPage func(objCount uint64)) (outputChan <-chan pipeline.WalkDirEntry) {
	ch := make(chan pipeline.WalkDirEntry)
	outputChan = ch

	logger := listObjectsLogger().Named("ListGenerator").With(
		"params", params,
	)

	generator := func() {
		defer func() {
			close(ch)
			logger.Info("closed output channel")
		}()

		dst := params.Destination
		page := params.PaginationParams
		var requestedItems int
	MainLoop:
		for {
			requestedItems = 0

			req, err := newListRequest(ctx, cfg, dst, page, params.Recursive)
			if err != nil {
				logger.Warnw("failed to create request", "err", err)
				return
			}

			resp, err := SendRequest(ctx, req)
			if err != nil {
				logger.Warnw("failed to send request", "err", err)
				return
			}

			var result listObjectsRequestResponse
			result, err = UnwrapResponse[listObjectsRequestResponse](resp, req)
			if err != nil {
				logger.Warnw("list request failed", "err", err, "req", (*mgcHttpPkg.LogRequest)(req))
				select {
				case <-ctx.Done():
					logger.Debugw("context.Done()", "err", err)
				case ch <- pipeline.NewSimpleWalkDirEntry[*BucketContent](dst.Path(), nil, err):
				}
				return
			}

			if onNewPage != nil {
				onNewPage(uint64(len(result.Contents)))
			}
			for _, prefix := range result.CommonPrefixes {
				dirEntry := pipeline.NewSimpleWalkDirEntry(
					path.Join(dst.Path(), prefix.Path),
					prefix,
					nil,
				)
				select {
				case <-ctx.Done():
					logger.Debugw("context.Done()", "err", ctx.Err())
					return
				case ch <- dirEntry:
					requestedItems++
					if requestedItems >= page.MaxItems {
						logger.Infow("item limit reached", "limit", params.PaginationParams.MaxItems)
						break MainLoop
					}
				}
			}

			for _, content := range result.Contents {
				dirEntry := pipeline.NewSimpleWalkDirEntry(
					content.Key,
					content,
					nil,
				)

				select {
				case <-ctx.Done():
					logger.Debugw("context.Done()", "err", ctx.Err())
					return
				case ch <- dirEntry:
					requestedItems++
					if requestedItems >= page.MaxItems {
						logger.Infow("item limit reached", "limit", params.PaginationParams.MaxItems)
						break MainLoop
					}
				}
			}

			page.ContinuationToken = result.NextContinuationToken
			page.MaxItems = page.MaxItems - requestedItems
			if !result.IsTruncated || page.MaxItems <= 0 {
				logger.Info("finished reading contents")
				break
			}

		}
	}

	logger.Info("list generation start")
	go generator()
	return
}
