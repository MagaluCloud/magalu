package object_lock

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type setObjectLockParams struct {
	Object          mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Specifies the object whose lock is being requested" mgc:"positional"`
	RetainUntilDate string           `json:"retain_until_date" jsonschema:"description=Timestamp in ISO 8601 format,example=2025-10-03T00:00:00"`
	Mode            string           `json:"mode,omitempty" jsonschema:"description=Lock mode,enum=COMPLIANCE,enum=GOVERNANCE,default=COMPLIANCE" mgc:"hidden"`
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
		return fmt.Sprintf("Successfully set Object Locking for object %q", result.Source().Parameters["dst"])
	})

	return exec
})

func setObjectLocking(ctx context.Context, params setObjectLockParams, cfg common.Config) (result core.Value, err error) {
	req, err := newSetObjectLockingRequest(ctx, params, cfg)
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

func newSetObjectLockingRequest(ctx context.Context, p setObjectLockParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostWithPath(cfg, common.NewBucketNameFromURI(p.Object), p.Object.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, string(url), nil)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := req.URL.Query()
	query.Add("retention", "")
	req.URL.RawQuery = query.Encode()

	getBody := func() (io.ReadCloser, error) {
		var parsedTime time.Time

		parsedTime, err = time.Parse("2006-01-02T15:04:05", p.RetainUntilDate)
		if err != nil {
			return nil, core.UsageError{Err: err}
		}
		bodyObj := common.DefaultObjectRetentionBody(parsedTime.In(time.Now().Location()))
		if p.Mode == string(common.ObjectLockModeGovernance) {
			bodyObj.Mode = common.ObjectLockModeGovernance
		}
		body, err := xml.MarshalIndent(bodyObj, "", "  ")
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
