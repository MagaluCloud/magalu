package object_lock

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type setBucketObjectLockParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Name of the bucket to set object locking for its objects,example=my-bucket" mgc:"positional"`
	Days   uint              `json:"days,omitempty" jsonschema:"description=Number of days to lock new objects for. Cannot be used alongside 'years',example=30"`
	Years  uint              `json:"years,omitempty" jsonschema:"description=Number of years to lock new objects for. Cannot be used alongside 'days',example=5"`
	Mode   string            `json:"mode,omitempty" jsonschema:"description=Lock mode,enum=COMPLIANCE,enum=GOVERNANCE,default=COMPLIANCE,required" mgc:"hidden"`
}

var getSet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "set number of either days or years to lock new objects for",
		},
		setObjectLocking,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully set Object Locking for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func setObjectLocking(ctx context.Context, params setBucketObjectLockParams, cfg common.Config) (result core.Value, err error) {
	if params.Days != 0 && params.Years != 0 {
		return nil, fmt.Errorf("Must include either days or years, but not both")
	}
	if params.Days == 0 && params.Years == 0 {
		return nil, fmt.Errorf("Missing parameter `days` or `years`")
	}

	req, err := newSetBucketObjectLockingRequest(ctx, params, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	err = common.ExtractErr(resp, req)
	if err != nil {
		return
	}

	return
}

func newSetBucketObjectLockingRequest(ctx context.Context, p setBucketObjectLockParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("object-lock", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url.String(), nil)
	if err != nil {
		return nil, err
	}

	getBody := func() (io.ReadCloser, error) {
		bodyObj := common.DefaultObjectLockingBody
		if p.Days != 0 {
			bodyObj.Rule.DefaultRetention.Days = int(p.Days)
		} else {
			bodyObj.Rule.DefaultRetention.Years = int(p.Years)
		}
		if p.Mode == string(common.ObjectLockModeGovernance) {
			bodyObj.Rule.DefaultRetention.Mode = common.ObjectLockModeGovernance
		}
		body, err := xml.Marshal(bodyObj)
		// fmt.Println(string(body))
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(body)
		return io.NopCloser(reader), nil
	}

	req.Body, err = getBody()
	if err != nil {
		return nil, err
	}
	req.GetBody = getBody

	return req, nil
}
