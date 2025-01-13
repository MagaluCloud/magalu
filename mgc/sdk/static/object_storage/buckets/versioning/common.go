package versioning

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type versioningConfiguration struct {
	Status    string `xml:"Status"`
	MfaDelete string `xml:"MfaDelete,omitempty"`

	Namespace string   `xml:"xmlns,omitempty,attr" json:"-"`
	XMLName   struct{} `xml:"VersioningConfiguration" json:"-"`
}

func newSetBucketVersioningRequest(
	ctx context.Context,
	bucketName common.BucketName,
	cfg common.Config,
	versioning versioningConfiguration,
) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Set("versioning", "")

	url.RawQuery = query.Encode()

	body, err := xml.Marshal(versioning)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(ctx, http.MethodPut, url.String(), bytes.NewBuffer(body))
}
